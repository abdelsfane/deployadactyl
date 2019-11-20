package delete_test

import (
	"github.com/compozed/deployadactyl/interfaces"
	"github.com/compozed/deployadactyl/state/delete"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("event binding", func() {
	Describe("DeleteStartedEventBinding", func() {
		Describe("Accept", func() {
			Context("when accept takes a correct event", func() {
				It("should return true", func() {
					deleteBind := delete.NewDeleteStartedEventBinding(nil)

					deleteEvent := delete.DeleteStartedEvent{}
					Expect(deleteBind.Accepts(deleteEvent)).Should(Equal(true))
				})
			})
			Context("when accept takes incorrect event", func() {
				It("should return false", func() {
					deleteBind := delete.NewDeleteStartedEventBinding(nil)

					event := interfaces.Event{}
					Expect(deleteBind.Accepts(event)).Should(Equal(false))
				})
			})
		})
		Describe("Emit", func() {
			Context("when emit takes a correct event", func() {
				It("should invoke handler", func() {
					invoked := false
					deleteFunc := func(event delete.DeleteStartedEvent) error {
						invoked = true
						return nil
					}
					deleteBind := delete.NewDeleteStartedEventBinding(deleteFunc)
					deleteEvent := delete.DeleteStartedEvent{}
					deleteBind.Emit(deleteEvent)

					Expect(invoked).Should(Equal(true))
				})
			})
			Context("when emit takes incorrect event", func() {
				It("should return error", func() {
					invoked := false
					deleteFunc := func(event delete.DeleteStartedEvent) error {
						invoked = true
						return nil
					}
					deleteBind := delete.NewDeleteStartedEventBinding(deleteFunc)
					event := interfaces.Event{}
					err := deleteBind.Emit(event)

					Expect(invoked).Should(Equal(false))
					Expect(err).ShouldNot(BeNil())
					Expect(err.Error()).Should(Equal("invalid event type"))
				})
			})
		})

	})
	Describe("DeleteSuccessEventBinding", func() {
		Describe("Accept", func() {
			Context("when accept takes a correct event", func() {
				It("should return true", func() {
					deleteSuccessBind := delete.NewDeleteSuccessEventBinding(nil)

					deleteEvent := delete.DeleteSuccessEvent{}
					Expect(deleteSuccessBind.Accepts(deleteEvent)).Should(Equal(true))
				})
			})
			Context("when accept takes incorrect event", func() {
				It("should return false", func() {
					deleteSuccessBind := delete.NewDeleteSuccessEventBinding(nil)

					event := interfaces.Event{}
					Expect(deleteSuccessBind.Accepts(event)).Should(Equal(false))
				})
			})
		})
		Describe("Emit", func() {
			Context("when emit takes a correct event", func() {
				It("should invoke handler", func() {
					invoked := false
					deleteFunc := func(event delete.DeleteSuccessEvent) error {
						invoked = true
						return nil
					}
					deleteSuccessBind := delete.NewDeleteSuccessEventBinding(deleteFunc)
					deleteEvent := delete.DeleteSuccessEvent{}
					deleteSuccessBind.Emit(deleteEvent)

					Expect(invoked).Should(Equal(true))
				})
			})
			Context("when emit takes incorrect event", func() {
				It("should return error", func() {
					invoked := false
					deleteFunc := func(event delete.DeleteSuccessEvent) error {
						invoked = true
						return nil
					}
					deleteSuccessBind := delete.NewDeleteSuccessEventBinding(deleteFunc)
					event := interfaces.Event{}
					err := deleteSuccessBind.Emit(event)

					Expect(invoked).Should(Equal(false))
					Expect(err).ShouldNot(BeNil())
					Expect(err.Error()).Should(Equal("invalid event type"))
				})
			})
		})

	})
	Describe("DeleteFailureEventBinding", func() {
		Describe("Accept", func() {
			Context("when accept takes a correct event", func() {
				It("should return true", func() {
					binding := delete.NewDeleteFailureEventBinding(nil)

					deleteEvent := delete.DeleteFailureEvent{}
					Expect(binding.Accepts(deleteEvent)).Should(Equal(true))
				})
			})
			Context("when accept takes incorrect event", func() {
				It("should return false", func() {
					binding := delete.NewDeleteFailureEventBinding(nil)

					event := interfaces.Event{}
					Expect(binding.Accepts(event)).Should(Equal(false))
				})
			})
		})
		Describe("Emit", func() {
			Context("when emit takes a correct event", func() {
				It("should invoke handler", func() {
					invoked := false
					deleteFunc := func(event delete.DeleteFailureEvent) error {
						invoked = true
						return nil
					}
					binding := delete.NewDeleteFailureEventBinding(deleteFunc)
					deleteEvent := delete.DeleteFailureEvent{}
					binding.Emit(deleteEvent)

					Expect(invoked).Should(Equal(true))
				})
			})
			Context("when emit takes incorrect event", func() {
				It("should return error", func() {
					invoked := false
					deleteFunc := func(event delete.DeleteFailureEvent) error {
						invoked = true
						return nil
					}
					binding := delete.NewDeleteFailureEventBinding(deleteFunc)
					event := interfaces.Event{}
					err := binding.Emit(event)

					Expect(invoked).Should(Equal(false))
					Expect(err).ShouldNot(BeNil())
					Expect(err.Error()).Should(Equal("invalid event type"))
				})
			})
		})

	})
	Describe("DeleteFinishEventBinding", func() {
		Describe("Accept", func() {
			Context("when accept takes a correct event", func() {
				It("should return true", func() {
					binding := delete.NewDeleteFinishedEventBinding(nil)

					event := delete.DeleteFinishedEvent{}
					Expect(binding.Accepts(event)).Should(Equal(true))
				})
			})
			Context("when accept takes incorrect event", func() {
				It("should return false", func() {
					binding := delete.NewDeleteFinishedEventBinding(nil)

					event := interfaces.Event{}
					Expect(binding.Accepts(event)).Should(Equal(false))
				})
			})
		})
		Describe("Emit", func() {
			Context("when emit takes a correct event", func() {
				It("should invoke handler", func() {
					invoked := false
					deleteFunc := func(event delete.DeleteFinishedEvent) error {
						invoked = true
						return nil
					}
					binding := delete.NewDeleteFinishedEventBinding(deleteFunc)
					deleteEvent := delete.DeleteFinishedEvent{}
					binding.Emit(deleteEvent)

					Expect(invoked).Should(Equal(true))
				})
			})
			Context("when emit takes incorrect event", func() {
				It("should return error", func() {
					invoked := false
					deleteFunc := func(event delete.DeleteFinishedEvent) error {
						invoked = true
						return nil
					}
					binding := delete.NewDeleteFinishedEventBinding(deleteFunc)
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
