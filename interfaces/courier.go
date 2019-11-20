package interfaces

type CourierCreator interface {
	CreateCourier() (Courier, error)
}

// Courier interface.
type Courier interface {
	Login(foundationURL, username, password, org, space string, skipSSL bool) ([]byte, error)
	Delete(appName string) ([]byte, error)
	Push(appName, appLocation, hostname string, instances uint16) ([]byte, error)
	Rename(oldName, newName string) ([]byte, error)
	MapRoute(appName, domain, hostname string) ([]byte, error)
	MapRouteWithPath(appName, domain, hostname, path string) ([]byte, error)
	UnmapRoute(appName, domain, hostname string) ([]byte, error)
	UnmapRouteWithPath(appName, domain, hostname, path string) ([]byte, error)
	DeleteRoute(domain, hostname string) ([]byte, error)
	CreateService(service, plan, name string) ([]byte, error)
	BindService(appName, serviceName string) ([]byte, error)
	UnbindService(appName, serviceName string) ([]byte, error)
	DeleteService(serviceName string) ([]byte, error)
	Start(appName string) ([]byte, error)
	Stop(appName string) ([]byte, error)
	Restage(appName string) ([]byte, error)
	Logs(appName string) ([]byte, error)
	Exists(appName string) bool
	Cups(appName string, body string) ([]byte, error)
	Uups(appName string, body string) ([]byte, error)
	Domains() ([]string, error)
	CleanUp() error
	Services() ([]string, error)
}
