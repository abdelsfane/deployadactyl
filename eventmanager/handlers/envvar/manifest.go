package envvar

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/cloudfoundry-incubator/candiedyaml"
	I "github.com/compozed/deployadactyl/interfaces"
	"github.com/spf13/afero"
)

type manifestYaml struct {
	Applications []Application `yaml:"applications"`
}

type Application struct {
	Name              string   `yaml:"name"`
	Memory            string   `yaml:"memory,omitempty"`
	Timeout           *uint16  `yaml:"timeout,omitempty"`
	Instances         *uint16  `yaml:"instances,omitempty"`
	Path              string   `yaml:"path,omitempty"`
	Java_opts         string   `yaml:"JAVA_OPTS,omitempty"`
	Command           string   `yaml:"command,omitempty"`
	Buildpack         string   `yaml:"buildpack,omitempty"`
	Disk_quota        string   `yaml:"disk_quota,omitempty"`
	Domain            string   `yaml:"domain,omitempty"`
	Domains           []string `yaml:"domains,omitempty"`
	Stack             string   `yaml:"stack,omitempty"`
	Health_check_type string   `yaml:"health-check-type,omitempty"`
	Host              string   `yaml:"host,omitempty"`
	Hosts             []string `yaml:"hosts,omitempty"`
	No_Hostname       string   `yaml:"no-hostname,omitempty"`
	Routes            []struct {
		Route string `yaml:"route,omitempty"`
	} `yaml:"routes,omitempty"`
	Services []string          `yaml:"services,omitempty"`
	Env      map[string]string `yaml:"env,omitempty"`
}

//Contains state of a m
type Manifest struct {
	Name       string
	Yaml       string
	parsed     bool
	Log        I.DeploymentLogger
	FileSystem *afero.Afero
	Content    manifestYaml
}

type ManifestError struct {
	Err error
}

func (e ManifestError) Error() string {
	return fmt.Sprintf("cannot open or write m file: %s", e.Err)
}

func CreateManifest(appName string, content string, filesystem *afero.Afero, logger I.DeploymentLogger) (manifest *Manifest, err error) {
	manifest = &Manifest{Name: appName, Yaml: content, FileSystem: filesystem, Log: logger}
	_, err = manifest.UnMarshal()

	if err != nil {
		logger.Errorf("Error Occurred during manifest creation/unmarshal! Details: %+v", err)
	}

	return manifest, err
}

// GetInstances reads a Cloud Foundry m as a string and returns the number of Instances
// defined in the m, if there are any.
//
// Returns a point to a uint16. If Instances are not found or less than 1, it returns nil.
func (m *Manifest) GetInstances() *uint16 {
	var (
		result = true
		err    error
	)

	if !m.parsed {
		result, err = m.UnMarshal()
	}

	if !result || err != nil {
		return nil
	}

	if m.Content.Applications == nil || m.Content.Applications[0].Instances == nil || *m.Content.Applications[0].Instances < 1 {
		return nil
	}

	return m.Content.Applications[0].Instances
}

func (m *Manifest) AddEnvVar(name string, value string) (err error) {

	m.Log.Debugf("Attempting to add Map of Environment Variable [%s] to Manifest", name)

	result := true

	if !m.parsed {
		result, err = m.UnMarshal()
	}

	if !result || err != nil {
		return err
	}

	vars := make(map[string]string)

	if m.HasApplications() {

		if m.Content.Applications[0].Env != nil {
			vars = m.Content.Applications[0].Env
		}

		vars[name] = value
		m.Content.Applications[0].Env = vars
	}

	return err
}

func (m *Manifest) AddEnvironmentVariables(env map[string]string) (result bool, err error) {

	m.Log.Debugf("Attempting to add Map of Environment Variables to Manifest")

	result = false

	if m != nil && len(env) > 0 {
		//#Add Environment Variables if any in the request!
		for k, v := range env {
			err = m.AddEnvVar(k, v)
			if err != nil {
				return false, err
			}
		}

		result = true
	}

	return result, err
}

func (m *Manifest) HasApplications() bool {

	var (
		result bool = true
		err    error
	)

	if !m.parsed {
		result, err = m.UnMarshal()
	}

	if !result || err != nil {
		return false
	}

	if m.Content.Applications != nil && len(m.Content.Applications) > 0 {
		return true
	}

	return false
}

func (m *Manifest) UnMarshal() (result bool, err error) {
	result = false

	if m.Yaml != "" {
		m.Log.Debugf("UnMarshaling Yaml => %s", m.Yaml)
		err = candiedyaml.Unmarshal([]byte(m.Yaml), &m.Content)
		if err != nil {
			m.Log.Errorf("Error Unmarshalling Manifest! Details: %v", err)
		} else {
			m.parsed = true
			result = true

			m.Log.Debugf("UnMarshalled Manifest Contents = %+v", m.Content)
		}
	} else {
		m.Log.Infof("Found No Manifest Content for App [%s]. Adding blank m....", m.Name)
		result = true
		m.Content.Applications = append(m.Content.Applications, Application{Name: m.Name})
		m.parsed = true
		return result, nil
	}

	return result, err
}

func (m *Manifest) Marshal() (content string) {
	m.Log.Debugf("Marshaling Manifest Contents = %+v", m.Content)

	resultBytes, err := candiedyaml.Marshal(m.Content)

	if err != nil {
		m.Log.Errorf("Error occurred marshalling Manifest Yaml! Details: %v", err)
		return m.Yaml
	}

	content = string(resultBytes)

	return content

}

func (m *Manifest) WriteManifest(destination string, includePrefix bool) error {

	manifest := m.Marshal()

	if manifest != "" {

		if includePrefix && !strings.HasPrefix(manifest, "---") {
			manifest = fmt.Sprintf("---\n%s", manifest)
		}

		manifestFile, err := m.FileSystem.OpenFile(path.Join(destination, "manifest.yml"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)

		if err != nil {
			return ManifestError{err}
		}
		defer manifestFile.Close()

		_, err = fmt.Fprint(manifestFile, manifest)
		if err != nil {
			return ManifestError{err}
		}
	}

	return nil
}
