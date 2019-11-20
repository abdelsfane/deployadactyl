package prechecker_test

import (
	"errors"
	"net/http"
	"net/http/httptest"

	. "github.com/compozed/deployadactyl/controller/deployer/prechecker"
	"github.com/compozed/deployadactyl/mocks"
	S "github.com/compozed/deployadactyl/structs"

	I "github.com/compozed/deployadactyl/interfaces"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"reflect"
)

var _ = Describe("Prechecker", func() {
	Describe("AssertAllFoundationsUp", func() {
		var (
			httpStatus     int
			foundationURls []string
			prechecker     Prechecker
			eventManager   *mocks.EventManager
			testServer     *httptest.Server
			environment    S.Environment
			event          I.Event
		)

		BeforeEach(func() {
			foundationURls = []string{}

			eventManager = &mocks.EventManager{}
			prechecker = Prechecker{EventManager: eventManager}

			testServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				foundationURls = append(foundationURls, r.URL.Path)
				w.WriteHeader(httpStatus)
			}))

			environment = S.Environment{
				Foundations: []string{testServer.URL},
			}
		})

		AfterEach(func() {
			testServer.Close()
		})

		Context("when no foundations are given", func() {
			It("returns an error and emits an event", func() {
				environment.Foundations = nil

				event = I.Event{
					Type: "validate.foundationsUnavailable",
					Data: S.PrecheckerEventData{
						Environment: environment,
						Description: "no foundations configured",
					},
				}
				eventManager.EmitCall.Returns.Error = append(eventManager.EmitCall.Returns.Error, nil)

				err := prechecker.AssertAllFoundationsUp(environment)
				Expect(err).To(MatchError(NoFoundationsConfiguredError{}))

				Expect(eventManager.EmitCall.Received.Events[0]).To(Equal(event))
			})
			It("calls EmitEvent", func() {
				environment.Foundations = nil

				err := prechecker.AssertAllFoundationsUp(environment)
				Expect(err).To(MatchError(NoFoundationsConfiguredError{}))

				Expect(reflect.TypeOf(eventManager.EmitEventCall.Received.Events[0])).To(Equal(reflect.TypeOf(FoundationsUnavailableEvent{})))
			})
			It("calls provides environment to EmitEvent", func() {
				environment.Foundations = nil

				err := prechecker.AssertAllFoundationsUp(environment)
				Expect(err).To(MatchError(NoFoundationsConfiguredError{}))

				ievent := eventManager.EmitEventCall.Received.Events[0].(FoundationsUnavailableEvent)
				Expect(ievent.Environment).To(Equal(environment))
			})
			It("calls provides description to EmitEvent", func() {
				environment.Foundations = nil

				err := prechecker.AssertAllFoundationsUp(environment)
				Expect(err).To(MatchError(NoFoundationsConfiguredError{}))

				ievent := eventManager.EmitEventCall.Received.Events[0].(FoundationsUnavailableEvent)
				Expect(ievent.Description).To(ContainSubstring("no foundations configured"))
			})
		})

		Context("when the client returns an error", func() {
			It("returns an error and emits an event", func() {
				environment.Foundations = []string{"bork"}

				event = I.Event{
					Type: "validate.foundationsUnavailable",
					Data: S.PrecheckerEventData{
						Environment: environment,
						Description: "no foundations configured",
					},
				}
				eventManager.EmitCall.Returns.Error = append(eventManager.EmitCall.Returns.Error, nil)

				err := prechecker.AssertAllFoundationsUp(environment)

				Expect(err.Error()).To(ContainSubstring(InvalidGetRequestError{"bork", errors.New("")}.Error()))
			})
		})

		Context("when all foundations return a 200 OK", func() {
			It("returns a nil error", func() {
				httpStatus = http.StatusOK

				Expect(prechecker.AssertAllFoundationsUp(environment)).To(Succeed())

				Expect(foundationURls).To(ConsistOf("/v2/info"))
			})
		})

		Context("when a foundation returns a 500 internal server error", func() {
			It("returns an error and emits an event", func() {
				event = I.Event{
					Type: "validate.foundationsUnavailable",
					Data: S.PrecheckerEventData{
						Environment: environment,
						Description: "deploy aborted: one or more CF foundations unavailable",
					},
				}
				eventManager.EmitCall.Returns.Error = append(eventManager.EmitCall.Returns.Error, nil)

				httpStatus = http.StatusInternalServerError

				Expect(prechecker.AssertAllFoundationsUp(environment)).ToNot(Succeed())

				Expect(foundationURls).To(ConsistOf("/v2/info"))
				Expect(eventManager.EmitCall.Received.Events[0]).ToNot(BeNil())
			})
			It("calls EmitEvent", func() {
				httpStatus = http.StatusInternalServerError

				prechecker.AssertAllFoundationsUp(environment)

				Expect(reflect.TypeOf(eventManager.EmitEventCall.Received.Events[0])).To(Equal(reflect.TypeOf(FoundationsUnavailableEvent{})))
			})
			It("provides environment to EmitEvent", func() {
				httpStatus = http.StatusInternalServerError

				prechecker.AssertAllFoundationsUp(environment)

				Expect(eventManager.EmitEventCall.Received.Events[0].(FoundationsUnavailableEvent).Environment).To(Equal(environment))
			})
			It("provides description to EmitEvent", func() {
				httpStatus = http.StatusInternalServerError

				prechecker.AssertAllFoundationsUp(environment)

				ievent := eventManager.EmitEventCall.Received.Events[0].(FoundationsUnavailableEvent)
				Expect(ievent.Description).To(ContainSubstring("deploy aborted: one or more CF foundations unavailable: http://127.0.0.1"))
			})
		})

		Context("when a foundation returns a 404 not found", func() {
			It("returns an error and emits an event", func() {
				event = I.Event{
					Type: "validate.foundationsUnavailable",
					Data: S.PrecheckerEventData{
						Environment: environment,
						Description: "deploy aborted: one or more CF foundations unavailable: http://127.0.0.1:51844: 404 Not Found",
					},
				}
				eventManager.EmitCall.Returns.Error = append(eventManager.EmitCall.Returns.Error, nil)

				httpStatus = http.StatusNotFound

				Expect(prechecker.AssertAllFoundationsUp(environment)).ToNot(Succeed())

				Expect(eventManager.EmitCall.Received.Events[0]).ToNot(BeNil())
			})
		})

		Describe("NewFoundationsUnavailableEventBinding", func() {
			Describe("Accept", func() {
				Context("when accept takes a correct event", func() {
					It("should return true", func() {
						binding := NewFoundationsUnavailableEventBinding(nil)

						event := FoundationsUnavailableEvent{}
						Expect(binding.Accepts(event)).Should(Equal(true))
					})
				})
				Context("when accept takes incorrect event", func() {
					It("should return false", func() {
						binding := NewFoundationsUnavailableEventBinding(nil)

						event := I.Event{}
						Expect(binding.Accepts(event)).Should(Equal(false))
					})
				})
			})
			Describe("Emit", func() {
				Context("when emit takes a correct event", func() {
					It("should invoke handler", func() {
						invoked := false
						handler := func(event FoundationsUnavailableEvent) error {
							invoked = true
							return nil
						}
						binding := NewFoundationsUnavailableEventBinding(handler)
						event := FoundationsUnavailableEvent{}
						binding.Emit(event)

						Expect(invoked).Should(Equal(true))
					})
				})
				Context("when emit takes incorrect event", func() {
					It("should return error", func() {
						invoked := false
						handler := func(event FoundationsUnavailableEvent) error {
							invoked = true
							return nil
						}
						binding := NewFoundationsUnavailableEventBinding(handler)
						event := I.Event{}
						err := binding.Emit(event)

						Expect(invoked).Should(Equal(false))
						Expect(err).ShouldNot(BeNil())
						Expect(err.Error()).Should(Equal("invalid event type"))
					})
				})
			})
		})
	})
})
