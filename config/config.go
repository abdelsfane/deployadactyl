// Package config holds all specified configuration information aggregated from all possible inputs.
package config

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/cloudfoundry-incubator/candiedyaml"
	"github.com/compozed/deployadactyl/controller/deployer/error_finder"
	"github.com/compozed/deployadactyl/geterrors"
	"github.com/compozed/deployadactyl/interfaces"
	s "github.com/compozed/deployadactyl/structs"
)

const DefaultConfigPath = "./config.yml"

// Config is a representation of a config yaml. It can contain multiple Environments.
type Config struct {
	Username      string
	Password      string
	Environments  map[string]s.Environment
	Port          int
	ErrorMatchers []interfaces.ErrorMatcher
}

type configYaml struct {
	Environments       []s.Environment            `yaml:",flow"`
	MatcherDescriptors []s.ErrorMatcherDescriptor `yaml:"error_matchers,flow"`
}

type foundationYaml struct {
	Foundations []string
}

type ConfigConstructor func() (Config, error)

// Default returns a new Config struct with information from environment variables and the default config file (./config.yml).
func Default(getenv func(string) string) (Config, error) {
	return Custom(getenv, DefaultConfigPath)
}

// Custom returns a new Config struct with information from environment variables and a custom config file.
func Custom(getenv func(string) string, configPath string) (Config, error) {
	foundationConfig, err := parseConfig(configPath)
	if err != nil {
		return Config{}, err
	}

	environments, err := getEnvironmentsFromConfig(foundationConfig)
	if err != nil {
		return Config{}, err
	}

	errormatchers := getErrorMatchersFromConfig(foundationConfig)
	if err != nil {
		return Config{}, err
	}

	return createConfig(getenv, environments, errormatchers)
}

func createConfig(getenv func(string) string, environments map[string]s.Environment, errormatchers []interfaces.ErrorMatcher) (Config, error) {
	getter := geterrors.WrapFunc(getenv)

	username := getter.Get("CF_USERNAME")
	password := getter.Get("CF_PASSWORD")

	if err := getter.Err("missing environment variables"); err != nil {
		return Config{}, err
	}

	port, err := getPortFromEnv(getenv)
	if err != nil {
		return Config{}, err
	}

	config := Config{
		Username:      username,
		Password:      password,
		Port:          port,
		Environments:  environments,
		ErrorMatchers: errormatchers,
	}
	return config, nil
}

func getPortFromEnv(getenv func(string) string) (int, error) {
	envPort := getenv("PORT")
	if envPort == "" {
		envPort = "8080"
	}

	cfgPort, err := strconv.Atoi(envPort)
	if err != nil {
		return 0, fmt.Errorf("cannot parse $PORT: %s: %s", envPort, err)
	}

	return cfgPort, nil
}

func getErrorMatchersFromConfig(foundationConfig configYaml) []interfaces.ErrorMatcher {

	matchers := make([]interfaces.ErrorMatcher, 0, 0)

	if foundationConfig.MatcherDescriptors != nil || len(foundationConfig.MatcherDescriptors) > 0 {
		factory := error_finder.ErrorMatcherFactory{}
		for _, descriptor := range foundationConfig.MatcherDescriptors {
			matcher, err := factory.CreateErrorMatcher(descriptor)
			if err == nil {
				matchers = append(matchers, matcher)
			}
		}
	}
	return matchers
}

func getEnvironmentsFromConfig(foundationConfig configYaml) (map[string]s.Environment, error) {

	if foundationConfig.Environments == nil || len(foundationConfig.Environments) == 0 {
		return nil, EnvironmentsNotSpecifiedError{}
	}

	environments := map[string]s.Environment{}
	for _, environment := range foundationConfig.Environments {
		if environment.Name == "" || environment.Foundations == nil || len(environment.Foundations) == 0 {
			return nil, MissingParameterError{}
		}

		if environment.Instances < 1 {
			environment.Instances = 1
		}

		environments[strings.ToLower(environment.Name)] = environment
	}

	return environments, nil
}

func parseConfig(configPath string) (configYaml, error) {
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		return configYaml{}, err
	}

	foundationConfig, err := parseYamlFromBody(file)
	if err != nil {
		return configYaml{}, err
	}
	return foundationConfig, nil
}

func parseYamlFromBody(data []byte) (configYaml, error) {
	var foundationConfig configYaml

	err := candiedyaml.Unmarshal(data, &foundationConfig)
	if err != nil {
		return configYaml{}, ParseYamlError{err}
	}

	return foundationConfig, nil
}
