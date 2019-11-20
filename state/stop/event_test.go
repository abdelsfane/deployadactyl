package stop_test

import (
	"github.com/compozed/deployadactyl/interfaces"
	"github.com/compozed/deployadactyl/state/stop"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("event binding", func() {
	Describe("StopStartedEventBinding", func() {
		Describe("Accept", func() {
			Context("when accept takes a correct event", func() {
				It("should return true", func() {
					stopBind := stop.NewStopStartedEventBinding(nil)

					stopEvent := stop.StopStartedEvent{}
					Expect(stopBind.Accepts(stopEvent)).Should(Equal(true))
				})
			})
			Context("when accept takes incorrect event", func() {
				It("should return false", func() {
					stopBind := stop.NewStopStartedEventBinding(nil)

					event := interfaces.Event{}
					Expect(stopBind.Accepts(event)).Should(Equal(false))
				})
			})
		})
		Describe("Emit", func() {
			Context("when emit takes a correct event", func() {
				It("should invoke handler", func() {
					invoked := false
					stopFunc := func(event stop.StopStartedEvent) error {
						invoked = true
						return nil
					}
					stopBind := stop.NewStopStartedEventBinding(stopFunc)
					stopEvent := stop.StopStartedEvent{}
					stopBind.Emit(stopEvent)

					Expect(invoked).Should(Equal(true))
				})
			})
			Context("when emit takes incorrect event", func() {
				It("should return error", func() {
					invoked := false
					stopFunc := func(event stop.StopStartedEvent) error {
						invoked = true
						return nil
					}
					stopBind := stop.NewStopStartedEventBinding(stopFunc)
					event := interfaces.Event{}
					err := stopBind.Emit(event)

					Expect(invoked).Should(Equal(false))
					Expect(err).ShouldNot(BeNil())
					Expect(err.Error()).Should(Equal("invalid event type"))
				})
			})
		})

	})
	Describe("StopSuccessEventBinding", func() {
		Describe("Accept", func() {
			Context("when accept takes a correct event", func() {
				It("should return true", func() {
					stopSuccessBind := stop.NewStopSuccessEventBinding(nil)

					stopEvent := stop.StopSuccessEvent{}
					Expect(stopSuccessBind.Accepts(stopEvent)).Should(Equal(true))
				})
			})
			Context("when accept takes incorrect event", func() {
				It("should return false", func() {
					stopSuccessBind := stop.NewStopSuccessEventBinding(nil)

					event := interfaces.Event{}
					Expect(stopSuccessBind.Accepts(event)).Should(Equal(false))
				})
			})
		})
		Describe("Emit", func() {
			Context("when emit takes a correct event", func() {
				It("should invoke handler", func() {
					invoked := false
					stopFunc := func(event stop.StopSuccessEvent) error {
						invoked = true
						return nil
					}
					stopSuccessBind := stop.NewStopSuccessEventBinding(stopFunc)
					stopEvent := stop.StopSuccessEvent{}
					stopSuccessBind.Emit(stopEvent)

					Expect(invoked).Should(Equal(true))
				})
			})
			Context("when emit takes incorrect event", func() {
				It("should return error", func() {
					invoked := false
					stopFunc := func(event stop.StopSuccessEvent) error {
						invoked = true
						return nil
					}
					stopSuccessBind := stop.NewStopSuccessEventBinding(stopFunc)
					event := interfaces.Event{}
					err := stopSuccessBind.Emit(event)

					Expect(invoked).Should(Equal(false))
					Expect(err).ShouldNot(BeNil())
					Expect(err.Error()).Should(Equal("invalid event type"))
				})
			})
		})

	})
	Describe("StopFailureEventBinding", func() {
		Describe("Accept", func() {
			Context("when accept takes a correct event", func() {
				It("should return true", func() {
					binding := stop.NewStopFailureEventBinding(nil)

					stopEvent := stop.StopFailureEvent{}
					Expect(binding.Accepts(stopEvent)).Should(Equal(true))
				})
			})
			Context("when accept takes incorrect event", func() {
				It("should return false", func() {
					binding := stop.NewStopFailureEventBinding(nil)

					event := interfaces.Event{}
					Expect(binding.Accepts(event)).Should(Equal(false))
				})
			})
		})
		Describe("Emit", func() {
			Context("when emit takes a correct event", func() {
				It("should invoke handler", func() {
					invoked := false
					stopFunc := func(event stop.StopFailureEvent) error {
						invoked = true
						return nil
					}
					binding := stop.NewStopFailureEventBinding(stopFunc)
					stopEvent := stop.StopFailureEvent{}
					binding.Emit(stopEvent)

					Expect(invoked).Should(Equal(true))
				})
			})
			Context("when emit takes incorrect event", func() {
				It("should return error", func() {
					invoked := false
					stopFunc := func(event stop.StopFailureEvent) error {
						invoked = true
						return nil
					}
					binding := stop.NewStopFailureEventBinding(stopFunc)
					event := interfaces.Event{}
					err := binding.Emit(event)

					Expect(invoked).Should(Equal(false))
					Expect(err).ShouldNot(BeNil())
					Expect(err.Error()).Should(Equal("invalid event type"))
				})
			})
		})

	})
	Describe("StopFinishEventBinding", func() {
		Describe("Accept", func() {
			Context("when accept takes a correct event", func() {
				It("should return true", func() {
					binding := stop.NewStopFinishedEventBinding(nil)

					event := stop.StopFinishedEvent{}
					Expect(binding.Accepts(event)).Should(Equal(true))
				})
			})
			Context("when accept takes incorrect event", func() {
				It("should return false", func() {
					binding := stop.NewStopFinishedEventBinding(nil)

					event := interfaces.Event{}
					Expect(binding.Accepts(event)).Should(Equal(false))
				})
			})
		})
		Describe("Emit", func() {
			Context("when emit takes a correct event", func() {
				It("should invoke handler", func() {
					invoked := false
					stopFunc := func(event stop.StopFinishedEvent) error {
						invoked = true
						return nil
					}
					binding := stop.NewStopFinishedEventBinding(stopFunc)
					stopEvent := stop.StopFinishedEvent{}
					binding.Emit(stopEvent)

					Expect(invoked).Should(Equal(true))
				})
			})
			Context("when emit takes incorrect event", func() {
				It("should return error", func() {
					invoked := false
					stopFunc := func(event stop.StopFinishedEvent) error {
						invoked = true
						return nil
					}
					binding := stop.NewStopFinishedEventBinding(stopFunc)
					event := interfaces.Event{}
					err := binding.Emit(event)

					Expect(invoked).Should(Equal(false))
					Expect(err).ShouldNot(BeNil())
					Expect(err.Error()).Should(Equal("invalid event type"))
				})
			})
		})

	})
})
