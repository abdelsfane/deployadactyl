// Package pusher handles pushing to individual Cloud Foundry instances.
package push

import (
	"fmt"
	"io"

	H "github.com/compozed/deployadactyl/eventmanager/handlers/healthchecker"
	R "github.com/compozed/deployadactyl/eventmanager/handlers/routemapper"
	I "github.com/compozed/deployadactyl/interfaces"
	"github.com/compozed/deployadactyl/state"
	S "github.com/compozed/deployadactyl/structs"
)

// TemporaryNameSuffix is used when deploying the new application in order to
// not overide the existing application name.
const TemporaryNameSuffix = "-new-build-"

// Pusher has a courier used to push applications to Cloud Foundry.
// It represents logging into a single foundation to perform operations.
type Pusher struct {
	Courier        I.Courier
	DeploymentInfo S.DeploymentInfo
	EventManager   I.EventManager
	Response       io.ReadWriter
	Log            I.DeploymentLogger
	FoundationURL  string
	AppPath        string
	Environment    S.Environment
	Fetcher        I.Fetcher
	CFContext      I.CFContext
	Auth           I.Authorization
	HealthChecker  H.HealthChecker
	RouteMapper    R.RouteMapper
}

// Login will login to a Cloud Foundry instance.
func (p Pusher) Initially() error {
	p.Log.Debugf(
		`logging into cloud foundry with parameters:
		foundation URL: %+v
		username: %+v
		org: %+v
		space: %+v`,
		p.FoundationURL, p.DeploymentInfo.Username, p.DeploymentInfo.Org, p.DeploymentInfo.Space,
	)

	output, err := p.Courier.Login(
		p.FoundationURL,
		p.DeploymentInfo.Username,
		p.DeploymentInfo.Password,
		p.DeploymentInfo.Org,
		p.DeploymentInfo.Space,
		p.DeploymentInfo.SkipSSL,
	)
	p.Response.Write(output)
	if err != nil {
		p.Log.Errorf("could not login to %s", p.FoundationURL)
		return state.LoginError{p.FoundationURL, output}
	}

	p.Log.Infof("logged into cloud foundry %s", p.FoundationURL)

	return nil
}

// Push pushes a single application to a Clound Foundry instance using blue green deployment.
// Blue green is done by pushing a new application with the appName+TemporaryNameSuffix+UUID.
// It pushes the new application with the existing appName route.
// It will map a load balanced domain if provided in the config.yml.
//
// Returns Cloud Foundry logs if there is an error.

func (p Pusher) Verify() error {
	return nil
}

func (p Pusher) Execute() error {

	var (
		tempAppWithUUID = p.DeploymentInfo.AppName + TemporaryNameSuffix + p.DeploymentInfo.UUID
		err             error
	)

	err = p.pushApplication(tempAppWithUUID, p.AppPath)
	if err != nil {
		return err
	}

	if p.DeploymentInfo.HealthCheckEndpoint != "" {
		healthCheckRequest := H.HealthCheckRequest{
			HealthCheckEndpoint: p.DeploymentInfo.HealthCheckEndpoint,
			Courier:             p.Courier,
			Logger:              p.Log,
			Environment:         p.DeploymentInfo.Environment,
			FoundationUrl:       p.FoundationURL,
			TempAppWithUUID:     tempAppWithUUID,
			UUID:                p.DeploymentInfo.UUID,
		}

		err = p.HealthChecker.HealthChecker(healthCheckRequest)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p Pusher) PostExecute() error {
	tempAppWithUUID := p.DeploymentInfo.AppName + TemporaryNameSuffix + p.DeploymentInfo.UUID

	routeMapperRequest := R.RouteMapperRequest{
		Logger:          p.Log,
		Courier:         p.Courier,
		Manifest:        p.DeploymentInfo.Manifest,
		AppPath:         p.DeploymentInfo.AppPath,
		TempAppWithUUID: tempAppWithUUID,
		Application:     p.DeploymentInfo.AppName,
		UUID:            p.DeploymentInfo.UUID,
		FoundationUrl:   p.FoundationURL,
	}

	err := p.RouteMapper.CustomRouteMapper(routeMapperRequest)
	if err != nil {
		return err
	}

	if p.DeploymentInfo.Domain != "" {
		err := p.mapTempAppToLoadBalancedDomain(tempAppWithUUID)
		if err != nil {
			return err
		}
	}

	return nil
}

// FinishPush will delete the original application if it existed. It will always
// rename the the newly pushed application to the appName.
func (p Pusher) Success() error {
	if p.Courier.Exists(p.DeploymentInfo.AppName) {
		err := p.unMapLoadBalancedRoute()
		if err != nil {
			return err
		}

		err = p.deleteApplication(p.DeploymentInfo.AppName)
		if err != nil {
			return err
		}
	}

	err := p.renameNewBuildToOriginalAppName()
	if err != nil {
		return err
	}

	return nil
}

// UndoPush is only called when a Push fails. If it is not the first deployment, UndoPush will
// delete the temporary application that was pushed.
// If is the first deployment, UndoPush will rename the failed push to have the appName.
func (p Pusher) Undo() error {

	tempAppWithUUID := p.DeploymentInfo.AppName + TemporaryNameSuffix + p.DeploymentInfo.UUID
	if p.Environment.DisableRollback {
		p.Log.Errorf("%s: Failed to deploy, deployment not rolled back due to DisabledRollback=true", p.FoundationURL)

		return p.Success()
	} else {

		if p.Courier.Exists(p.DeploymentInfo.AppName) {
			p.Log.Errorf("%s: rolling back deploy of %s", p.FoundationURL, tempAppWithUUID)

			err := p.deleteApplication(tempAppWithUUID)
			if err != nil {
				return err
			}

		} else {
			p.Log.Errorf("%s: app %s did not previously exist: not rolling back", p.FoundationURL, p.DeploymentInfo.AppName)

			err := p.renameNewBuildToOriginalAppName()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// CleanUp removes the temporary directory created by the Executor.
func (p Pusher) Finally() error {
	return p.Courier.CleanUp()
}

func (p Pusher) pushApplication(appName, appPath string) error {
	p.Log.Debugf("%s: pushing app %s to %s", p.FoundationURL, appName, p.DeploymentInfo.Domain)
	p.Log.Debugf("%s: tempdir for app %s: %s", p.FoundationURL, appName, appPath)

	var (
		pushOutput          []byte
		cloudFoundryLogs    []byte
		err                 error
		cloudFoundryLogsErr error
	)

	defer func() { p.Response.Write(cloudFoundryLogs) }()
	defer func() { p.Response.Write(pushOutput) }()

	pushOutput, err = p.Courier.Push(appName, appPath, p.DeploymentInfo.AppName, p.DeploymentInfo.Instances)
	p.Log.Infof("%s: push output from Cloud Foundry: \n%s", p.FoundationURL, pushOutput)
	if err != nil {
		defer func() { p.Log.Errorf("%s: logs from %s: \n%s", p.FoundationURL, appName, cloudFoundryLogs) }()

		cloudFoundryLogs, cloudFoundryLogsErr = p.Courier.Logs(appName)
		if cloudFoundryLogsErr != nil {
			return state.CloudFoundryGetLogsError{err, cloudFoundryLogsErr}
		}

		return state.PushError{}
	}

	p.Log.Infof("%s: successfully deployed new build %s", p.FoundationURL, appName)

	return nil
}

func (p Pusher) mapTempAppToLoadBalancedDomain(appName string) error {
	p.Log.Debugf("%s: mapping route for %s to %s", p.FoundationURL, p.DeploymentInfo.AppName, p.DeploymentInfo.Domain)

	out, err := p.Courier.MapRoute(appName, p.DeploymentInfo.Domain, p.DeploymentInfo.AppName)
	p.Log.Infof("%s: mapping output from Cloud Foundry: \n%s", p.FoundationURL, out)
	if err != nil {
		p.Log.Errorf("%s: could not map %s to %s", p.FoundationURL, p.DeploymentInfo.AppName, p.DeploymentInfo.Domain)
		return state.MapRouteError{out}
	}

	p.Log.Infof("%s: application route created: %s.%s", p.FoundationURL, p.DeploymentInfo.AppName, p.DeploymentInfo.Domain)

	fmt.Fprintf(p.Response, "application route created: %s.%s", p.DeploymentInfo.AppName, p.DeploymentInfo.Domain)

	return nil
}

func (p Pusher) unMapLoadBalancedRoute() error {
	if p.DeploymentInfo.Domain != "" {
		p.Log.Debugf("%s: unmapping route %s", p.FoundationURL, p.DeploymentInfo.AppName)

		out, err := p.Courier.UnmapRoute(p.DeploymentInfo.AppName, p.DeploymentInfo.Domain, p.DeploymentInfo.AppName)
		p.Log.Infof("%s: unmapping output from Cloud Foundry: \n%s", p.FoundationURL, out)
		if err != nil {
			p.Log.Errorf("%s: could not unmap %s", p.FoundationURL, p.DeploymentInfo.AppName)
			return state.UnmapRouteError{p.DeploymentInfo.AppName, out}
		}

		p.Log.Infof("%s: unmapped route %s", p.FoundationURL, p.DeploymentInfo.AppName)
	}

	return nil
}

func (p Pusher) deleteApplication(appName string) error {
	p.Log.Debugf("%s: deleting %s", p.FoundationURL, appName)

	out, err := p.Courier.Delete(appName)
	p.Log.Infof("%s: deletion output from Cloud Foundry: \n%s", p.FoundationURL, out)
	if err != nil {
		p.Log.Errorf("%s: could not delete %s", p.FoundationURL, appName)
		p.Log.Errorf("%s: deletion error %s", p.FoundationURL, err.Error())
		p.Log.Errorf("%s: deletion output", p.FoundationURL, string(out))
		return state.DeleteApplicationError{appName, out}
	}

	p.Log.Infof("%s: deleted %s", p.FoundationURL, appName)

	return nil
}

func (p Pusher) renameNewBuildToOriginalAppName() error {
	p.Log.Debugf("%s: renaming %s to %s", p.FoundationURL, p.DeploymentInfo.AppName+TemporaryNameSuffix+p.DeploymentInfo.UUID, p.DeploymentInfo.AppName)

	out, err := p.Courier.Rename(p.DeploymentInfo.AppName+TemporaryNameSuffix+p.DeploymentInfo.UUID, p.DeploymentInfo.AppName)
	p.Log.Infof("%s: rename output from Cloud Foundry: \n%s", p.FoundationURL, out)
	if err != nil {
		p.Log.Errorf("%s: could not rename %s to %s", p.FoundationURL, p.DeploymentInfo.AppName+TemporaryNameSuffix+p.DeploymentInfo.UUID, p.DeploymentInfo.AppName)
		return state.RenameError{p.DeploymentInfo.AppName + TemporaryNameSuffix + p.DeploymentInfo.UUID, out}
	}

	p.Log.Infof("%s: renamed %s to %s", p.FoundationURL, p.DeploymentInfo.AppName+TemporaryNameSuffix+p.DeploymentInfo.UUID, p.DeploymentInfo.AppName)

	return nil
}

func (p Pusher) mapTemporaryRoute(tempAppWithUUID, domain string, log I.DeploymentLogger) error {
	log.Debugf("mapping temporary route %s.%s", tempAppWithUUID, domain)

	out, err := p.Courier.MapRoute(tempAppWithUUID, domain, tempAppWithUUID)
	p.Log.Infof("%s: mapping output from Cloud Foundry: \n%s", p.FoundationURL, out)
	if err != nil {
		log.Errorf("failed to map temporary route: %s", out)
		return state.MapRouteError{out}
	}
	log.Infof("mapped temporary route %s.%s", tempAppWithUUID, domain)

	return nil
}

func (p Pusher) deleteTemporaryRoute(tempAppWithUUID, domain string, log I.DeploymentLogger) error {
	log.Debugf("deleting temporary route %s.%s", tempAppWithUUID, domain)

	out, err := p.Courier.DeleteRoute(domain, tempAppWithUUID)
	p.Log.Infof("%s: route deletion output from Cloud Foundry: \n%s", p.FoundationURL, out)
	if err != nil {
		log.Errorf("failed to delete temporary route: %s", out)
		return state.MapRouteError{out}
	}

	log.Infof("deleted temporary route %s.%s", tempAppWithUUID, domain)

	return nil
}

func (p Pusher) unmapTemporaryRoute(tempAppWithUUID, domain string, log I.DeploymentLogger) {
	log.Debugf("unmapping temporary route %s.%s", tempAppWithUUID, domain)

	out, err := p.Courier.UnmapRoute(tempAppWithUUID, domain, tempAppWithUUID)
	p.Log.Infof("%s: unmapping output from Cloud Foundry: \n%s", p.FoundationURL, out)
	if err != nil {
		log.Errorf("failed to unmap temporary route: %s", out)
	} else {
		log.Infof("unmapped temporary route %s.%s", tempAppWithUUID, domain)
	}

	log.Infof("finished health check")
}
