// Package deployer will deploy your application.
package deployer

import (
	"fmt"
	"io"
	"net/http"

	"crypto/tls"
	"log"
	"os"

	"encoding/base64"
	"github.com/compozed/deployadactyl/config"
	I "github.com/compozed/deployadactyl/interfaces"
	S "github.com/compozed/deployadactyl/structs"
)

const (
	successfulDeploy = `Your deploy was successful! (^_^)b
If you experience any problems after this point, check that you can manually push your application to Cloud Foundry on a lower environment.
It is likely that it is an error with your application and not with Deployadactyl.
Thanks for using Deployadactyl! Please push down pull up on your lap bar and exit to your left.

`

	deploymentOutput = `Deployment Parameters:
Artifact URL: %s,
Username:     %s,
Environment:  %s,
Org:          %s,
Space:        %s,
AppName:      %s`
)

type SilentDeployer struct {
}

func (d SilentDeployer) Deploy(deploymentInfo *S.DeploymentInfo, env S.Environment, actionCreator I.ActionCreator, response io.ReadWriter) *I.DeployResponse {
	url := os.Getenv("SILENT_DEPLOY_URL")
	deployResponse := &I.DeployResponse{}

	request, err := http.NewRequest("POST", fmt.Sprintf(url+"/%s/%s/%s", deploymentInfo.Org, deploymentInfo.Space, deploymentInfo.AppName), deploymentInfo.Body)
	if err != nil {
		log.Println(fmt.Sprintf("Silent deployer request err: %s", err))
		deployResponse.Error = err
	}
	usernamePassword := base64.StdEncoding.EncodeToString([]byte(deploymentInfo.Username + ":" + deploymentInfo.Password))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", usernamePassword)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}

	resp, err := client.Do(request)
	if err != nil {
		log.Println(fmt.Sprintf("Silent deployer response err: %s", err))
		deployResponse.StatusCode = resp.StatusCode
		deployResponse.Error = err
	}

	deployResponse.StatusCode = resp.StatusCode
	deployResponse.Error = err
	return deployResponse
}

type DeployerConstructor func(config config.Config, blueGreener I.BlueGreener, preChecker I.Prechecker, eventManager I.EventManager, randomizer I.Randomizer, errorFinder I.ErrorFinder, logger I.DeploymentLogger) I.Deployer

func NewDeployer(c config.Config, bg I.BlueGreener, p I.Prechecker, em I.EventManager, r I.Randomizer, ef I.ErrorFinder, l I.DeploymentLogger) I.Deployer {
	return &Deployer{
		Config:       c,
		BlueGreener:  bg,
		Prechecker:   p,
		EventManager: em,
		Randomizer:   r,
		ErrorFinder:  ef,
		Log:          l,
	}
}

type Deployer struct {
	Config       config.Config
	BlueGreener  I.BlueGreener
	Prechecker   I.Prechecker
	EventManager I.EventManager
	Randomizer   I.Randomizer
	ErrorFinder  I.ErrorFinder
	Log          I.DeploymentLogger
}

func (d Deployer) Deploy(deploymentInfo *S.DeploymentInfo, env S.Environment, actionCreator I.ActionCreator, response io.ReadWriter) *I.DeployResponse {

	deployResponse := &I.DeployResponse{
		DeploymentInfo: deploymentInfo,
	}

	d.Log.Debug("prechecking the foundations")
	err := d.Prechecker.AssertAllFoundationsUp(env)
	if err != nil {
		d.Log.Error(err)
		deployResponse.StatusCode = http.StatusInternalServerError
		deployResponse.Error = err
		return deployResponse
	}

	defer func() { actionCreator.CleanUp() }()
	err = actionCreator.SetUp()
	if err != nil {
		deployResponse.StatusCode = http.StatusInternalServerError
		deployResponse.Error = err
		return deployResponse
	}

	err = actionCreator.OnStart()
	if err != nil {
		deployResponse.StatusCode = http.StatusInternalServerError
		deployResponse.Error = err
		return deployResponse
	}

	err = d.BlueGreener.Execute(actionCreator, env, response)

	resp := actionCreator.OnFinish(env, response, err)
	resp.DeploymentInfo = deploymentInfo
	return &resp
}
