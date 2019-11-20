package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/compozed/deployadactyl/creator"
	"github.com/compozed/deployadactyl/interfaces"
	"github.com/compozed/deployadactyl/state/push"
	"github.com/op/go-logging"
)

const (
	defaultConfigFilePath = "./config.yml"
	defaultLogLevel       = "DEBUG"
	logLevelEnvVarName    = "DEPLOYADACTYL_LOGLEVEL"
)

func main() {
	var (
		config               = flag.String("config", defaultConfigFilePath, "location of the config file")
		envVarHandlerEnabled = flag.Bool("env", false, "enable environment variable handling")
	)
	flag.Parse()

	level := os.Getenv(logLevelEnvVarName)
	if level == "" {
		level = defaultLogLevel
	}

	logLevel, err := logging.LogLevel(level)
	if err != nil {
		log.Fatal(err)
	}

	log := interfaces.DefaultLogger(os.Stdout, logLevel, "deployadactyl")
	log.Infof("log level : %s", level)

	c, err := creator.Custom(level, *config, creator.CreatorModuleProvider{})
	if err != nil {
		log.Fatal(err)
	}

	eventBindings := c.GetEventBindings()

	if *envVarHandlerEnabled {
		envVarHandler := c.CreateEnvVarHandler()
		log.Infof("registering environment variable event handler")
		eventBindings.AddBinding(push.NewArtifactRetrievalSuccessEventBinding(envVarHandler.ArtifactRetrievalSuccessEventHandler))
	}

	l := c.CreateListener()
	controller := c.CreateController()

	deploy := c.CreateControllerHandler(controller)

	log.Infof("Listening on Port %d", c.CreateConfig().Port)

	err = http.Serve(l, deploy)
	if err != nil {
		log.Fatal(err)
	}
}
