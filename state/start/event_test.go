package start_test

import (
	"github.com/compozed/deployadactyl/interfaces"
	"github.com/compozed/deployadactyl/state/start"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("event binding", func() {
	Describe("StartStartedEventBinding", func() {
		Describe("Accept", func() {
			Context("when accept takes a correct event", func() {
				It("should return true", func() {
					stopBind := start.NewStartStartedEventBinding(nil)

					stopEvent := start.StartStartedEvent{}
					Expect(stopBind.Accepts(stopEvent)).Should(Equal(true))
				})
			})
			Context("when accept takes incorrect event", func() {
				It("should return false", func() {
					stopBind := start.NewStartStartedEventBinding(nil)

					event := interfaces.Event{}
					Expect(stopBind.Accepts(event)).Should(Equal(false))
				})
			})
		})
		Describe("Emit", func() {
			Context("when emit takes a correct event", func() {
				It("should invoke handler", func() {
					invoked := false
					startFunc := func(event start.StartStartedEvent) error {
						invoked = true
						return nil
					}
					stopBind := start.NewStartStartedEventBinding(startFunc)
					stopEvent := start.StartStartedEvent{}
					stopBind.Emit(stopEvent)

					Expect(invoked).Should(Equal(true))
				})
			})
			Context("when emit takes incorrect event", func() {
				It("should return error", func() {
					invoked := false
					startFunc := func(event start.StartStartedEvent) error {
						invoked = true
						return nil
					}
					startBind := start.NewStartStartedEventBinding(startFunc)
					event := interfaces.Event{}
					err := startBind.Emit(event)

					Expect(invoked).Should(Equal(false))
					Expect(err).ShouldNot(BeNil())
					Expect(err.Error()).Should(Equal("invalid event type"))
				})
			})
		})

	})
	Describe("StartSuccessEventBinding", func() {
		Describe("Accept", func() {
			Context("when accept takes a correct event", func() {
				It("should return true", func() {
					startSuccessBind := start.NewStartSuccessEventBinding(nil)

					startEvent := start.StartSuccessEvent{}
					Expect(startSuccessBind.Accepts(startEvent)).Should(Equal(true))
				})
			})
			Context("when accept takes incorrect event", func() {
				It("should return false", func() {
					startSuccessBind := start.NewStartSuccessEventBinding(nil)

					event := interfaces.Event{}
					Expect(startSuccessBind.Accepts(event)).Should(Equal(false))
				})
			})
		})
		Describe("Emit", func() {
			Context("when emit takes a correct event", func() {
				It("should invoke handler", func() {
					invoked := false
					startFunc := func(event start.StartSuccessEvent) error {
						invoked = true
						return nil
					}
					stopSuccessBind := start.NewStartSuccessEventBinding(startFunc)
					stopEvent := start.StartSuccessEvent{}
					stopSuccessBind.Emit(stopEvent)

					Expect(invoked).Should(Equal(true))
				})
			})
			Context("when emit takes incorrect event", func() {
				It("should return error", func() {
					invoked := false
					startFunc := func(event start.StartSuccessEvent) error {
						invoked = true
						return nil
					}
					startSuccessBind := start.NewStartSuccessEventBinding(startFunc)
					event := interfaces.Event{}
					err := startSuccessBind.Emit(event)

					Expect(invoked).Should(Equal(false))
					Expect(err).ShouldNot(BeNil())
					Expect(err.Error()).Should(Equal("invalid event type"))
				})
			})
		})

	})
	Describe("StartFailureEventBinding", func() {
		Describe("Accept", func() {
			Context("when accept takes a correct event", func() {
				It("should return true", func() {
					binding := start.NewStartFailureEventBinding(nil)

					startEvent := start.StartFailureEvent{}
					Expect(binding.Accepts(startEvent)).Should(Equal(true))
				})
			})
			Context("when accept takes incorrect event", func() {
				It("should return false", func() {
					binding := start.NewStartFailureEventBinding(nil)

					event := interfaces.Event{}
					Expect(binding.Accepts(event)).Should(Equal(false))
				})
			})
		})
		Describe("Emit", func() {
			Context("when emit takes a correct event", func() {
				It("should invoke handler", func() {
					invoked := false
					startFunc := func(event start.StartFailureEvent) error {
						invoked = true
						return nil
					}
					binding := start.NewStartFailureEventBinding(startFunc)
					startEvent := start.StartFailureEvent{}
					binding.Emit(startEvent)

					Expect(invoked).Should(Equal(true))
				})
			})
			Context("when emit takes incorrect event", func() {
				It("should return error", func() {
					invoked := false
					startFunc := func(event start.StartFailureEvent) error {
						invoked = true
						return nil
					}
					binding := start.NewStartFailureEventBinding(startFunc)
					event := interfaces.Event{}
					err := binding.Emit(event)

					Expect(invoked).Should(Equal(false))
					Expect(err).ShouldNot(BeNil())
					Expect(err.Error()).Should(Equal("invalid event type"))
				})
			})
		})

	})
	Describe("StartFinishEventBinding", func() {
		Describe("Accept", func() {
			Context("when accept takes a correct event", func() {
				It("should return true", func() {
					binding := start.NewStartFinishedEventBinding(nil)

					event := start.StartFinishedEvent{}
					Expect(binding.Accepts(event)).Should(Equal(true))
				})
			})
			Context("when accept takes incorrect event", func() {
				It("should return false", func() {
					binding := start.NewStartFinishedEventBinding(nil)

					event := interfaces.Event{}
					Expect(binding.Accepts(event)).Should(Equal(false))
				})
			})
		})
		Describe("Emit", func() {
			Context("when emit takes a correct event", func() {
				It("should invoke handler", func() {
					invoked := false
					startFunc := func(event start.StartFinishedEvent) error {
						invoked = true
						return nil
					}
					binding := start.NewStartFinishedEventBinding(startFunc)
					stopEvent := start.StartFinishedEvent{}
					binding.Emit(stopEvent)

					Expect(invoked).Should(Equal(true))
				})
			})
			Context("when emit takes incorrect event", func() {
				It("should return error", func() {
					invoked := false
					startFunc := func(event start.StartFinishedEvent) error {
						invoked = true
						return nil
					}
					binding := start.NewStartFinishedEventBinding(startFunc)
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
