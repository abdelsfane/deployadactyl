// Package courier interfaces with the Executor to run specific Cloud Foundry CLI commands.
package courier

import (
	"fmt"
	"strings"

	I "github.com/compozed/deployadactyl/interfaces"
	"github.com/go-errors/errors"
)

type CourierConstructor func(executor I.Executor) I.Courier

func NewCourier(executor I.Executor) I.Courier {
	return Courier{
		Executor: executor,
	}
}

// Courier has an Executor to execute Cloud Foundry commands.
type Courier struct {
	Executor I.Executor
}

// Login runs the Cloud Foundry login command.
//
// Returns the combined standard output and standard error.
func (c Courier) Login(foundationURL, username, password, org, space string, skipSSL bool) ([]byte, error) {
	var s string
	if skipSSL {
		s = "--skip-ssl-validation"
	}

	return c.Executor.Execute("login", "-a", foundationURL, "-u", username, "-p", password, "-o", org, "-s", space, s)
}

func (c Courier) CreateService(service, plan, name string) ([]byte, error) {
	return c.Executor.Execute("create-service", service, plan, name)
}

func (c Courier) BindService(appName, dbName string) ([]byte, error) {
	return c.Executor.Execute("bind-service", appName, dbName)
}

func (c Courier) UnbindService(appName, dbName string) ([]byte, error) {
	return c.Executor.Execute("unbind-service", appName, dbName)
}

func (c Courier) DeleteService(serviceName string) ([]byte, error) {
	return c.Executor.Execute("delete-service", serviceName, "-f")
}

func (c Courier) Restage(appName string) ([]byte, error) {
	return c.Executor.Execute("restage", appName)
}

func (c Courier) Start(appName string) ([]byte, error) {
	return c.Executor.Execute("start", appName)
}

func (c Courier) Stop(appName string) ([]byte, error) {
	return c.Executor.Execute("stop", appName)
}

// Delete runs the Cloud Foundry delete command.
// Returns the combined standard output and standard error.
func (c Courier) Delete(appName string) ([]byte, error) {
	return c.Executor.Execute("delete", appName, "-f")
}

// Push runs the Cloud Foundry push command.
//
// Returns the combined standard output and standard error.
func (c Courier) Push(appName, appLocation, hostname string, instances uint16) ([]byte, error) {
	return c.Executor.ExecuteInDirectory(appLocation, "push", appName, "-i", fmt.Sprint(instances), "-n", hostname)
}

// Rename runs the Cloud Foundry rename command.
//
// Returns the combined standard output and standard error.
func (c Courier) Rename(appName, newAppName string) ([]byte, error) {
	return c.Executor.Execute("rename", appName, newAppName)
}

// MapRoute runs the Cloud Foundry map-route command and added path arguement
//
// Returns the combined standard output and standard error.
func (c Courier) MapRouteWithPath(appName, domain, hostname, path string) ([]byte, error) {
	return c.Executor.Execute("map-route", appName, domain, "-n", hostname, "--path", path)
}

// MapRoute runs the Cloud Foundry map-route command.
//
// Returns the combined standard output and standard error.
func (c Courier) MapRoute(appName, domain, hostname string) ([]byte, error) {
	return c.Executor.Execute("map-route", appName, domain, "-n", hostname)
}

// UnmapRoute runs the Cloud Foundry unmap-route command.
//
// Returns the combined standard output and standard error.
func (c Courier) UnmapRouteWithPath(appName, domain, hostname, path string) ([]byte, error) {
	return c.Executor.Execute("unmap-route", appName, domain, "-n", hostname, "--path", path)
}

// UnmapRoute runs the Cloud Foundry unmap-route command.
//
// Returns the combined standard output and standard error.
func (c Courier) UnmapRoute(appName, domain, hostname string) ([]byte, error) {
	return c.Executor.Execute("unmap-route", appName, domain, "-n", hostname)
}

func (c Courier) DeleteRoute(domain, hostname string) ([]byte, error) {
	return c.Executor.Execute("delete-route", domain, "-n", hostname, "-f")
}

// Logs runs the Cloud Foundry logs command.
//
// Returns the combined standard output and standard error.
func (c Courier) Logs(appName string) ([]byte, error) {
	logs, err := c.Executor.Execute("logs", appName, "--recent")
	return logs, err
}

// Cups runs the Cloud Foundry CUPS command to create user provided
// services.
//
// Returns the combined standard output and standard error.
func (c Courier) Cups(appName string, body string) ([]byte, error) {
	return c.Executor.Execute("cups", appName, "-p", body)
}

// Uups runs the Cloud Foundry UUPS command to update a user provided serivce
func (c Courier) Uups(appName string, body string) ([]byte, error) {
	return c.Executor.Execute("uups", appName, "-p", body)
}

// Exists checks to see whether the application name exists already.
//
// Returns true if the application exists.
func (c Courier) Exists(appName string) bool {
	_, err := c.Executor.Execute("app", appName)
	return err == nil
}

// Domains returns a list of domain in a foundation.
//
// Returns the combined standard output and standard error.
func (c Courier) Domains() ([]string, error) {
	output, err := c.Executor.Execute("domains")

	domains := strings.Split(string(output), "\n")[2:]
	for i, domain := range domains {
		domains[i] = strings.Split(domain, " ")[0]
	}

	return domains, err
}

// CleanUp removes the temporary directory created by the Executor.
func (c Courier) CleanUp() error {
	return c.Executor.CleanUp()
}

// Lists all services in current space
//
// Returns a list of services available
func (c Courier) Services() ([]string, error) {
	output, err := c.Executor.Execute("services")
	if err != nil {
		return []string{}, errors.New(fmt.Sprintf("Execution of services call failed: %s", err.Error()))
	}

	services := strings.Split(string(output), "\n")[3:]
	for i, domain := range services {
		services[i] = strings.Split(domain, " ")[0]
	}
	return services, nil
}
