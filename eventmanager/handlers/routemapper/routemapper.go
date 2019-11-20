package routemapper

import (
	"strings"

	"github.com/cloudfoundry-incubator/candiedyaml"
	I "github.com/compozed/deployadactyl/interfaces"
	"github.com/spf13/afero"
)

// RouteMapper will map additional routes to an application at
// deploy time if they are specified in the manifest.
type RouteMapper struct {
	Courier    I.Courier
	FileSystem *afero.Afero
}

type manifest struct {
	Applications []application
}

type application struct {
	CustomRoutes []route `yaml:"custom-routes"`
}

type route struct {
	Route string
}

type RouteMapperRequest struct {
	Logger          I.DeploymentLogger
	Courier         I.Courier
	Manifest        string
	AppPath         string
	TempAppWithUUID string
	Application     string
	UUID            string
	FoundationUrl   string
}

func (r RouteMapper) CustomRouteMapper(request RouteMapperRequest) error {
	log := request.Logger.Log
	log.Debugf("%s %s: starting route mapper", request.UUID, request.FoundationUrl)

	r.Courier = request.Courier

	manifestBytes, err := r.readManifest(request.Manifest, request.AppPath, request.Logger, request.UUID, request.FoundationUrl)
	if err != nil || manifestBytes == nil {
		return err
	}
	m := &manifest{}

	log.Debugf("%s %s: looking for routes in the manifest", request.UUID, request.FoundationUrl)
	err = candiedyaml.Unmarshal(manifestBytes, m)
	if err != nil {
		log.Errorf("%s %s: failed to parse manifest: %s", request.UUID, request.FoundationUrl, err.Error())
		return err
	}

	if m.Applications == nil || len(m.Applications[0].CustomRoutes) == 0 {
		log.Infof("%s %s: finished mapping routes: no routes to map", request.UUID, request.FoundationUrl)
		return nil
	}

	log.Infof("%s %s: found %d routes in the manifest", request.UUID, request.FoundationUrl, len(m.Applications[0].CustomRoutes))

	domains, _ := r.Courier.Domains()

	log.Debugf("%s %s: mapping routes to %s", request.UUID, request.FoundationUrl, request.TempAppWithUUID)
	return r.routeMapper(m, request.TempAppWithUUID, domains, request.Application, request.Logger, request.UUID, request.FoundationUrl)
}

func isRouteADomainInTheFoundation(route string, domains []string) bool {
	for _, domain := range domains {

		if route == domain {
			return true
		}
	}
	return false
}

func (r RouteMapper) readManifest(manifest, appPath string, log I.DeploymentLogger, uuid, foundationUrl string) ([]byte, error) {
	var (
		manifestBytes []byte
		err           error
	)
	if manifest != "" {
		manifestBytes = []byte(manifest)
		return manifestBytes, nil
	} else if appPath != "" {
		manifestBytes, err = r.FileSystem.ReadFile(appPath + "/manifest.yml")
		if err != nil {
			log.Errorf("%s %s: failed to read manifest file: %s", uuid, foundationUrl, err.Error())
			return nil, ReadFileError{err}
		}
		return manifestBytes, nil
	} else {
		log.Infof("%s %s: finished mapping routes: no manifest found", uuid, foundationUrl)
		return nil, nil
	}
}

// routeMapper is used to decide how to map an applications routes that are given to it from the manifest.
// if the route does not include appname or path it will map the given domain to the given application by default
// if the route has an app name it will remove the app name so it can map it with the given domain
// if the route has an app name and a path it will remove the app name so it can map it with the given domain and the path as well
func (r RouteMapper) routeMapper(manifest *manifest, tempAppWithUUID string, domains []string, appName string, log I.DeploymentLogger, uuid, foundationUrl string) error {
	for _, route := range manifest.Applications[0].CustomRoutes {
		var domainAndPath []string

		appNameAndDomain := strings.SplitN(route.Route, ".", 2)

		if len(appNameAndDomain) >= 2 {
			domainAndPath = strings.SplitN(appNameAndDomain[1], "/", 2)
		}

		if isRouteADomainInTheFoundation(route.Route, domains) {
			output, err := r.Courier.MapRoute(tempAppWithUUID, route.Route, appName)
			if err != nil {
				log.Errorf("failed to map route: %s: %s", route.Route, string(output))
				return MapRouteError{route.Route, output}
			}
		} else if len(appNameAndDomain) >= 2 && isRouteADomainInTheFoundation(appNameAndDomain[1], domains) {
			output, err := r.Courier.MapRoute(tempAppWithUUID, appNameAndDomain[1], appNameAndDomain[0])
			if err != nil {
				log.Errorf("failed to map route: %s: %s", route.Route, string(output))
				return MapRouteError{route.Route, output}
			}
		} else if domainAndPath != nil && isRouteADomainInTheFoundation(domainAndPath[0], domains) {
			output, err := r.Courier.MapRouteWithPath(tempAppWithUUID, domainAndPath[0], appNameAndDomain[0], domainAndPath[1])
			if err != nil {
				log.Error(MapRouteError{route.Route, output})
				return MapRouteError{route.Route, output}
			}
		} else {
			return InvalidRouteError{route.Route}
		}

		log.Infof("%s %s: mapped route %s to %s", uuid, foundationUrl, route.Route, tempAppWithUUID)
	}

	log.Infof("%s %s: route mapping successful: finished mapping routes", uuid, foundationUrl)
	return nil
}
