package healthchecker

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	I "github.com/compozed/deployadactyl/interfaces"
)

type HealthCheckerConstructor func(oldURL, newURL, silentDeployURL, silentDeployEnvironment string, client I.Client) HealthChecker

func NewHealthChecker(oldURL, newURL, silentDeployURL, silentDeployEnvironment string, client I.Client) HealthChecker {
	return HealthChecker{
		OldURL:                  oldURL,
		NewURL:                  newURL,
		SilentDeployURL:         silentDeployURL,
		SilentDeployEnvironment: silentDeployEnvironment,
		Client:                  client,
	}
}

type HealthChecker struct {
	// OldURL is the prepend on the foundationURL to replace in order to build the
	// newly pushed application URL.
	// Eg: "api.run.pivotal"
	OldURL string

	// NewUrl is what replaces OldURL in the OnEvent function.
	// Eg: "cfapps"
	NewURL string

	//SilentDeployURL represents any other url that doesn't match cfapps
	SilentDeployURL         string
	SilentDeployEnvironment string

	Client  I.Client
	Courier I.Courier
}

type HealthCheckRequest struct {
	HealthCheckEndpoint string
	Courier             I.Courier
	Logger              I.DeploymentLogger
	Environment         string
	FoundationUrl       string
	TempAppWithUUID     string
	UUID                string
}

func (h HealthChecker) HealthChecker(healthCheckRequest HealthCheckRequest) error {
	var (
		newFoundationURL string
		domain           string
	)

	h.Courier = healthCheckRequest.Courier

	healthCheckRequest.Logger.Log.Debugf("%s %s: starting health check", healthCheckRequest.UUID, healthCheckRequest.FoundationUrl)

	if healthCheckRequest.Environment != h.SilentDeployEnvironment {
		newFoundationURL = strings.Replace(healthCheckRequest.FoundationUrl, h.OldURL, h.NewURL, 1)
		domain = regexp.MustCompile(fmt.Sprintf("%s.*", h.NewURL)).FindString(newFoundationURL)
	} else {
		newFoundationURL = strings.Replace(healthCheckRequest.FoundationUrl, h.OldURL, h.SilentDeployURL, 1)
		domain = regexp.MustCompile(fmt.Sprintf("%s.*", h.SilentDeployURL)).FindString(newFoundationURL)
	}

	err := h.mapTemporaryRoute(healthCheckRequest.TempAppWithUUID, domain, healthCheckRequest.Logger, healthCheckRequest.UUID, healthCheckRequest.FoundationUrl)
	if err != nil {
		return err
	}

	// unmapTemporaryRoute will be called before deleteTemporaryRoute
	defer h.deleteTemporaryRoute(healthCheckRequest.TempAppWithUUID, domain, healthCheckRequest.Logger, healthCheckRequest.UUID, healthCheckRequest.FoundationUrl)
	defer h.unmapTemporaryRoute(healthCheckRequest.TempAppWithUUID, domain, healthCheckRequest.Logger, healthCheckRequest.UUID, healthCheckRequest.FoundationUrl)

	newFoundationURL = strings.Replace(newFoundationURL, h.NewURL, fmt.Sprintf("%s.%s", healthCheckRequest.TempAppWithUUID, h.NewURL), 1)

	return h.Check(newFoundationURL, healthCheckRequest.HealthCheckEndpoint, healthCheckRequest.Logger, healthCheckRequest.UUID, healthCheckRequest.FoundationUrl)
}

// Check takes a url and endpoint. It does an http.Get to get the response
// status and returns an error if it is not http.StatusOK.
func (h HealthChecker) Check(url, endpoint string, log I.DeploymentLogger, uuid, foundationUrl string) error {
	trimmedEndpoint := strings.TrimPrefix(endpoint, "/")

	log.Debugf("%s %s: checking route %s%s", uuid, foundationUrl, url, endpoint)

	resp, err := h.Client.Get(fmt.Sprintf("%s/%s", url, trimmedEndpoint))
	if err != nil {
		log.Error(ClientError{err})
		return ClientError{err}
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Errorf("%s %s: health check failed for %s/%s", uuid, foundationUrl, url, trimmedEndpoint)
		return HealthCheckError{resp.StatusCode, endpoint, body}
	}

	log.Infof("%s %s: health check successful for %s%s", uuid, foundationUrl, url, endpoint)
	return nil
}

func (h HealthChecker) mapTemporaryRoute(tempAppWithUUID, domain string, log I.DeploymentLogger, uuid, foundationUrl string) error {
	log.Debugf("%s %s: mapping temporary route %s.%s", uuid, foundationUrl, tempAppWithUUID, domain)

	out, err := h.Courier.MapRoute(tempAppWithUUID, domain, tempAppWithUUID)
	if err != nil {
		log.Errorf("%s %s: failed to map temporary route: %s", uuid, foundationUrl, out)
		return MapRouteError{tempAppWithUUID, domain}
	}
	log.Infof("%s %s: mapped temporary route %s.%s", uuid, foundationUrl, tempAppWithUUID, domain)

	return nil
}

func (h HealthChecker) deleteTemporaryRoute(tempAppWithUUID, domain string, log I.DeploymentLogger, uuid, foundationUrl string) error {
	log.Debugf("%s %s: deleting temporary route %s.%s", uuid, foundationUrl, tempAppWithUUID, domain)

	out, err := h.Courier.DeleteRoute(domain, tempAppWithUUID)
	if err != nil {
		log.Errorf("%s %s: failed to delete temporary route: %s", uuid, foundationUrl, out)
		return DeleteRouteError{tempAppWithUUID, domain}
	}

	log.Infof("%s %s: deleted temporary route %s.%s", uuid, foundationUrl, tempAppWithUUID, domain)

	return nil
}

func (h HealthChecker) unmapTemporaryRoute(tempAppWithUUID, domain string, log I.DeploymentLogger, uuid, foundationUrl string) {
	log.Debugf("%s %s: unmapping temporary route %s.%s", uuid, foundationUrl, tempAppWithUUID, domain)

	out, err := h.Courier.UnmapRoute(tempAppWithUUID, domain, tempAppWithUUID)
	if err != nil {
		log.Errorf("%s %s: failed to unmap temporary route: %s", uuid, foundationUrl, out)
	} else {
		log.Infof("%s %s: unmapped temporary route %s.%s", uuid, foundationUrl, tempAppWithUUID, domain)
	}

	log.Infof("%s %s: finished health check", uuid, foundationUrl)
}
