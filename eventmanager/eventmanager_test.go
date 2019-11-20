package eventmanager_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/op/go-logging"

	. "github.com/compozed/deployadactyl/eventmanager"
	I "github.com/compozed/deployadactyl/interfaces"
	"github.com/compozed/deployadactyl/mocks"
	"github.com/compozed/deployadactyl/randomizer"
	"github.com/compozed/deployadactyl/state/stop"
)

var _ = Describe("Events", func() {
	var (
		eventType       string
		eventData       string
		eventHandler    *mocks.Handler
		eventHandlerOne *mocks.Handler
		eventHandlerTwo *mocks.Handler
		eventManager    I.EventManager
		logBuffer       *gbytes.Buffer
		log             I.DeploymentLogger
	)

	BeforeEach(func() {
		eventType = "eventType-" + randomizer.StringRunes(10)
		eventData = "eventData-" + randomizer.StringRunes(10)

		eventHandler = &mocks.Handler{}
		eventHandlerOne = &mocks.Handler{}
		eventHandlerTwo = &mocks.Handler{}

		logBuffer = gbytes.NewBuffer()
		log = I.DeploymentLogger{Log: I.DefaultLogger(logBuffer, logging.DEBUG, "eventmanager_test")}

		eventManager = NewEventManager(log, []I.Binding{})
	})

	Context("when an event handler is registered", func() {
		It("should be successful", func() {
			eventManager := NewEventManager(log, []I.Binding{})

			Expect(eventManager.AddHandler(eventHandler, eventType)).To(Succeed())
		})

		It("should fail if a nil value is passed in as an argument", func() {
			eventManager := NewEventManager(log, []I.Binding{})

			err := eventManager.AddHandler(nil, eventType)

			Expect(err).To(MatchError(InvalidArgumentError{}))
		})
	})

	Context("when an event is emitted", func() {
		It("should call all event handlers", func() {
			eventHandlerOne.OnEventCall.Returns.Error = nil
			eventHandlerTwo.OnEventCall.Returns.Error = nil

			event := I.Event{Type: eventType, Data: eventData}

			eventManager.AddHandler(eventHandlerOne, eventType)
			eventManager.AddHandler(eventHandlerTwo, eventType)

			Expect(eventManager.Emit(event)).To(Succeed())

			Expect(eventHandlerOne.OnEventCall.Received.Event).To(Equal(event))
			Expect(eventHandlerTwo.OnEventCall.Received.Event).To(Equal(event))
		})

		It("should return an error if the handler returns an error", func() {
			eventHandler.OnEventCall.Returns.Error = errors.New("on event error")

			event := I.Event{Type: eventType, Data: eventData}

			eventManager.AddHandler(eventHandler, eventType)

			Expect(eventManager.Emit(event)).To(MatchError("on event error"))
			Expect(eventHandler.OnEventCall.Received.Event).To(Equal(event))
		})

		It("should log that the event is emitted", func() {
			eventHandler.OnEventCall.Returns.Error = nil

			event := I.Event{Type: eventType, Data: eventData}

			eventManager.AddHandler(eventHandler, eventType)

			Expect(eventManager.Emit(event)).To(Succeed())

			Expect(eventHandler.OnEventCall.Received.Event).To(Equal(event))
			//Eventually(logBuffer).Should(gbytes.Say("a %s event has been emitted", eventType))
		})
	})

	Context("when there are handlers registered for two different types of events", func() {
		It("only emits to the specified event", func() {
			eventHandlerOne.OnEventCall.Returns.Error = nil
			eventHandlerTwo.OnEventCall.Returns.Error = nil

			event := I.Event{Type: eventType, Data: eventData}

			eventManager.AddHandler(eventHandlerOne, eventType)
			eventManager.AddHandler(eventHandlerTwo, "anotherEventType-"+randomizer.StringRunes(10))

			Expect(eventManager.Emit(event)).To(Succeed())

			Expect(eventHandlerOne.OnEventCall.Received.Event).To(Equal(event))
			Expect(eventHandlerTwo.OnEventCall.Received.Event).ToNot(Equal(event))
		})
	})

	Context("when events are added to the event manager", func() {
		It("should bind each event", func() {

			binding := &mocks.EventBinding{}

			eventManager.AddBinding(binding)
			Expect(eventManager.(*EventManager).Bindings[0]).To(Equal(binding))
		})

		It("should emit each event", func() {
			binding := &mocks.EventBinding{}
			eventManager.AddBinding(binding)

			stopStartedEvent := stop.StopStartedEvent{}
			binding.AcceptsCall.Returns.Bool = true
			eventManager.EmitEvent(stopStartedEvent)
			Expect(binding.AcceptsCall.Received.Event).To(Equal(stopStartedEvent))
			Expect(binding.EmitCall.Received.Event).To(Equal(stopStartedEvent))
		})

		Context("when event is an incorrect type", func() {
			It("should not emit", func() {
				binding := &mocks.EventBinding{}
				eventManager.AddBinding(binding)

				stopStartedEvent := stop.StopStartedEvent{}
				binding.AcceptsCall.Returns.Bool = false
				eventManager.EmitEvent(stopStartedEvent)
				Expect(binding.AcceptsCall.Received.Event).To(Equal(stopStartedEvent))
				Expect(binding.EmitCall.Called.Bool).To(Equal(false))
			})
		})

		Context("when binding returns an error", func() {
			It("emit should return error", func() {
				binding := &mocks.EventBinding{}
				eventManager.AddBinding(binding)

				stopStartedEvent := stop.StopStartedEvent{}
				binding.AcceptsCall.Returns.Bool = true
				binding.EmitCall.Returns.Error = errors.New("emit error")
				err := eventManager.EmitEvent(stopStartedEvent)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("emit error"))

			})
		})

		Context("when a panic happens in Emit", func() {
			It("recovers", func() {
				binding := &mocks.EventBinding{}
				eventManager.AddBinding(binding)

				stopStartedEvent := stop.StopStartedEvent{}
				binding.AcceptsCall.Returns.Bool = true
				binding.EmitCall.ShouldPanic = true

				err := eventManager.EmitEvent(stopStartedEvent)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Recovered from a panic: "))
			})
		})
	})
})
