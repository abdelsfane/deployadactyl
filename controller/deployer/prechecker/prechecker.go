// Package prechecker checks that all the Cloud Foundry instances are running before a deploy.
package prechecker

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/compozed/deployadactyl/eventmanager"
	I "github.com/compozed/deployadactyl/interfaces"
	S "github.com/compozed/deployadactyl/structs"
	"github.com/go-errors/errors"
	"reflect"
)

type eventBinding struct {
	etype   reflect.Type
	handler func(event interface{}) error
}

func (s eventBinding) Accepts(event interface{}) bool {
	return reflect.TypeOf(event) == s.etype
}

func (b eventBinding) Emit(event interface{}) error {
	return b.handler(event)
}

type FoundationsUnavailableEvent struct {
	Environment S.Environment
	Description string
}

func (d FoundationsUnavailableEvent) Name() string {
	return "FoundationsUnavailableEvent"
}

func NewFoundationsUnavailableEventBinding(handler func(event FoundationsUnavailableEvent) error) I.Binding {
	return eventBinding{
		etype: reflect.TypeOf(FoundationsUnavailableEvent{}),
		handler: func(gevent interface{}) error {
			event, ok := gevent.(FoundationsUnavailableEvent)
			if ok {
				return handler(event)
			} else {
				return eventmanager.InvalidEventType{errors.New("invalid event type")}
			}
		},
	}
}

type PrecheckerConstructor func(eventManager I.EventManager) I.Prechecker

func NewPrechecker(eventManager I.EventManager) I.Prechecker {
	return Prechecker{
		EventManager: eventManager,
	}
}

// Prechecker has an eventmanager used to manage event if prechecks fail.
type Prechecker struct {
	EventManager I.EventManager
}

// AssertAllFoundationsUp will send a request to each Cloud Foundry instance and check that the response status code is 200 OK.
func (p Prechecker) AssertAllFoundationsUp(environment S.Environment) error {
	precheckerEventData := S.PrecheckerEventData{Environment: environment}
	event := FoundationsUnavailableEvent{
		Environment: environment,
	}

	if len(environment.Foundations) == 0 {
		precheckerEventData.Description = "no foundations configured"

		p.EventManager.Emit(I.Event{Type: "validate.foundationsUnavailable", Data: precheckerEventData})

		event.Description = precheckerEventData.Description

		p.EventManager.EmitEvent(event)
		return NoFoundationsConfiguredError{}
	}

	insecureClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
			ResponseHeaderTimeout: 15 * time.Second,
		},
	}

	for _, foundationURL := range environment.Foundations {
		resp, err := insecureClient.Get(fmt.Sprintf("%s/v2/info", foundationURL))
		if err != nil {
			return InvalidGetRequestError{foundationURL, err}
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			err := FoundationUnavailableError{foundationURL, resp.Status}

			precheckerEventData.Description = err.Error()
			event.Description = err.Error()

			p.EventManager.Emit(I.Event{Type: "validate.foundationsUnavailable", Data: precheckerEventData})
			p.EventManager.EmitEvent(event)

			return err
		}
	}

	return nil
}
