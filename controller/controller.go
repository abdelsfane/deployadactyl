// Package controller is responsible for handling requests from the Server.
package controller

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	"encoding/json"

	I "github.com/compozed/deployadactyl/interfaces"

	"net/http"
	"strings"

	"github.com/compozed/deployadactyl/config"
	"github.com/compozed/deployadactyl/randomizer"
	"github.com/compozed/deployadactyl/request"
	"github.com/gin-gonic/gin"
)

type RequestProcessorFactory func(uuid string, request interface{}, buffer *bytes.Buffer) I.RequestProcessor

// Controller is used to determine the type of request and process it accordingly.
type Controller struct {
	Log                     I.Logger
	RequestProcessorFactory RequestProcessorFactory
	Config                  config.Config
	ErrorFinder             I.ErrorFinder
}

func (c *Controller) PostRequestHandler(g *gin.Context) {
	cfContext := I.CFContext{
		Environment:  strings.ToLower(g.Param("environment")),
		Organization: strings.ToLower(g.Param("org")),
		Space:        strings.ToLower(g.Param("space")),
		Application:  strings.ToLower(g.Param("appName")),
	}

	user, pwd, _ := g.Request.BasicAuth()
	authorization := I.Authorization{
		Username: user,
		Password: pwd,
	}

	deploymentType := g.Request.Header.Get("Content-Type")

	response := &bytes.Buffer{}
	defer io.Copy(g.Writer, response)

	bodyBuffer, _ := ioutil.ReadAll(g.Request.Body)

	g.Request.Body.Close()

	deployment := I.Deployment{
		Authorization: authorization,
		CFContext:     cfContext,
		Type:          deploymentType,
		Body:          &bodyBuffer,
	}

	postRequest := request.PostRequest{}
	if deploymentType == "application/json" {
		err := json.Unmarshal(bodyBuffer, &postRequest)
		if err != nil {
			response.Write([]byte("Invalid request body."))
			g.Writer.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	postDeploymentRequest := request.PostDeploymentRequest{
		Deployment: deployment,
		Request:    postRequest,
	}
	if postRequest.UUID == "" {
		postRequest.UUID = randomizer.StringRunes(10)
	}

	log := I.DeploymentLogger{Log: c.Log, UUID: postRequest.UUID}
	log.Debugf("Request originated from: %+v", g.Request.RemoteAddr)

	deployResponse := c.RequestProcessorFactory(postRequest.UUID, postDeploymentRequest, response).Process()

	if deployResponse.Error != nil {
		g.Writer.WriteHeader(deployResponse.StatusCode)
		fmt.Fprintf(response, "cannot deploy application: %s\n", deployResponse.Error)
		return
	}

	g.Writer.WriteHeader(deployResponse.StatusCode)
}

func (c *Controller) PutRequestHandler(g *gin.Context) {
	cfContext := I.CFContext{
		Environment:  strings.ToLower(g.Param("environment")),
		Organization: strings.ToLower(g.Param("org")),
		Space:        strings.ToLower(g.Param("space")),
		Application:  strings.ToLower(g.Param("appName")),
	}

	response := &bytes.Buffer{}
	defer io.Copy(g.Writer, response)

	user, pwd, _ := g.Request.BasicAuth()
	authorization := I.Authorization{
		Username: user,
		Password: pwd,
	}

	bodyBuffer, _ := ioutil.ReadAll(g.Request.Body)
	g.Request.Body.Close()

	putRequest := request.PutRequest{}
	err := json.Unmarshal(bodyBuffer, &putRequest)
	if err != nil {
		response.Write([]byte("Invalid request body."))
		g.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	deployment := I.Deployment{
		Body:          &bodyBuffer,
		Authorization: authorization,
		CFContext:     cfContext,
		Type:          g.Request.Header.Get("Content-Type"),
	}

	putDeploymentRequest := request.PutDeploymentRequest{
		Deployment: deployment,
		Request:    putRequest,
	}

	if putRequest.UUID == "" {
		putRequest.UUID = randomizer.StringRunes(10)
	}

	log := I.DeploymentLogger{Log: c.Log, UUID: putRequest.UUID}
	log.Debugf("PUT Request originated from: %+v", g.Request.RemoteAddr)

	deployResponse := c.RequestProcessorFactory(putRequest.UUID, putDeploymentRequest, response).Process()
	if deployResponse.Error != nil {
		fmt.Fprintf(response, "cannot deploy application: %s\n", deployResponse.Error)
	}

	g.Writer.WriteHeader(deployResponse.StatusCode)
}

func (c *Controller) DeleteRequestHandler(g *gin.Context) {
	uuid := randomizer.StringRunes(10)
	log := I.DeploymentLogger{Log: c.Log, UUID: uuid}
	log.Debugf("DELETE Request originated from: %+v", g.Request.RemoteAddr)

	cfContext := I.CFContext{
		Environment:  strings.ToLower(g.Param("environment")),
		Organization: strings.ToLower(g.Param("org")),
		Space:        strings.ToLower(g.Param("space")),
		Application:  strings.ToLower(g.Param("appName")),
	}

	response := &bytes.Buffer{}
	defer io.Copy(g.Writer, response)

	user, pwd, _ := g.Request.BasicAuth()
	authorization := I.Authorization{
		Username: user,
		Password: pwd,
	}

	bodyBuffer, _ := ioutil.ReadAll(g.Request.Body)
	g.Request.Body.Close()

	deleteRequest := request.DeleteRequest{}
	err := json.Unmarshal(bodyBuffer, &deleteRequest)
	if err != nil {
		response.Write([]byte("Invalid request body."))
		g.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	deployment := I.Deployment{
		Body:          &bodyBuffer,
		Authorization: authorization,
		CFContext:     cfContext,
		Type:          g.Request.Header.Get("Content-Type"),
	}

	deleteDeploymentRequest := request.DeleteDeploymentRequest{
		Deployment: deployment,
		Request:    deleteRequest,
	}

	deployResponse := c.RequestProcessorFactory(uuid, deleteDeploymentRequest, response).Process()
	if deployResponse.Error != nil {
		fmt.Fprintf(response, "cannot delete application: %s\n", deployResponse.Error)
	}

	g.Writer.WriteHeader(deployResponse.StatusCode)
}
