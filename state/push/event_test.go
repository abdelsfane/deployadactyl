package push_test

import (
	"github.com/compozed/deployadactyl/interfaces"
	"github.com/compozed/deployadactyl/state/push"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("event binding", func() {
	Describe("DeployStartedEvent", func() {
		Describe("Accept", func() {
			Context("when accept takes a correct event", func() {
				It("should return true", func() {
					binding := push.NewDeployStartEventBinding(nil)

					event := push.DeployStartedEvent{}
					Expect(binding.Accepts(event)).Should(Equal(true))
				})
			})
			Context("when accept takes incorrect event", func() {
				It("should return false", func() {
					binding := push.NewDeployStartEventBinding(nil)

					event := interfaces.Event{}
					Expect(binding.Accepts(event)).Should(Equal(false))
				})
			})
		})
		Describe("Emit", func() {
			Context("when emit takes a correct event", func() {
				It("should invoke handler", func() {
					invoked := false
					handler := func(event push.DeployStartedEvent) error {
						invoked = true
						return nil
					}
					binding := push.NewDeployStartEventBinding(handler)
					event := push.DeployStartedEvent{}
					binding.Emit(event)

					Expect(invoked).Should(Equal(true))
				})
			})
			Context("when emit takes incorrect event", func() {
				It("should return error", func() {
					invoked := false
					handler := func(event push.DeployStartedEvent) error {
						invoked = true
						return nil
					}
					binding := push.NewDeployStartEventBinding(handler)
					event := interfaces.Event{}
					err := binding.Emit(event)

					Expect(invoked).Should(Equal(false))
					Expect(err).ShouldNot(BeNil())
					Expect(err.Error()).Should(Equal("invalid event type"))
				})
			})
		})
	})

	Describe("DeployFinishEvent", func() {
		Describe("Accept", func() {
			Context("when accept takes a correct event", func() {
				It("should return true", func() {
					binding := push.NewDeployFinishedEventBinding(nil)

					event := push.DeployFinishedEvent{}
					Expect(binding.Accepts(event)).Should(Equal(true))
				})
			})
			Context("when accept takes incorrect event", func() {
				It("should return false", func() {
					binding := push.NewDeployFinishedEventBinding(nil)

					event := interfaces.Event{}
					Expect(binding.Accepts(event)).Should(Equal(false))
				})
			})
		})
		Describe("Emit", func() {
			Context("when emit takes a correct event", func() {
				It("should invoke handler", func() {
					invoked := false
					pushFunc := func(event push.DeployFinishedEvent) error {
						invoked = true
						return nil
					}
					binding := push.NewDeployFinishedEventBinding(pushFunc)
					event := push.DeployFinishedEvent{}
					binding.Emit(event)

					Expect(invoked).Should(Equal(true))
				})
			})
			Context("when emit takes incorrect event", func() {
				It("should return error", func() {
					invoked := false
					pushFunc := func(event push.DeployFinishedEvent) error {
						invoked = true
						return nil
					}
					binding := push.NewDeployFinishedEventBinding(pushFunc)
					event := interfaces.Event{}
					err := binding.Emit(event)

					Expect(invoked).Should(Equal(false))
					Expect(err).ShouldNot(BeNil())
					Expect(err.Error()).Should(Equal("invalid event type"))
				})
			})
		})
	})

	Describe("DeploySuccessEvent", func() {
		Describe("Accept", func() {
			Context("when accept takes a correct event", func() {
				It("should return true", func() {
					binding := push.NewDeploySuccessEventBinding(nil)

					event := push.DeploySuccessEvent{}
					Expect(binding.Accepts(event)).Should(Equal(true))
				})
			})
			Context("when accept takes incorrect event", func() {
				It("should return false", func() {
					binding := push.NewDeploySuccessEventBinding(nil)

					event := interfaces.Event{}
					Expect(binding.Accepts(event)).Should(Equal(false))
				})
			})
		})
		Describe("Emit", func() {
			Context("when emit takes a correct event", func() {
				It("should invoke handler", func() {
					invoked := false
					handler := func(event push.DeploySuccessEvent) error {
						invoked = true
						return nil
					}
					binding := push.NewDeploySuccessEventBinding(handler)
					event := push.DeploySuccessEvent{}
					binding.Emit(event)

					Expect(invoked).Should(Equal(true))
				})
			})
			Context("when emit takes incorrect event", func() {
				It("should return error", func() {
					invoked := false
					handler := func(event push.DeploySuccessEvent) error {
						invoked = true
						return nil
					}
					binding := push.NewDeploySuccessEventBinding(handler)
					event := interfaces.Event{}
					err := binding.Emit(event)

					Expect(invoked).Should(Equal(false))
					Expect(err).ShouldNot(BeNil())
					Expect(err.Error()).Should(Equal("invalid event type"))
				})
			})
		})
	})

	Describe("DeployFailureEvent", func() {
		Describe("Accept", func() {
			Context("when accept takes a correct event", func() {
				It("should return true", func() {
					binding := push.NewDeployFailureEventBinding(nil)

					event := push.DeployFailureEvent{}
					Expect(binding.Accepts(event)).Should(Equal(true))
				})
			})
			Context("when accept takes incorrect event", func() {
				It("should return false", func() {
					binding := push.NewDeployFailureEventBinding(nil)

					event := interfaces.Event{}
					Expect(binding.Accepts(event)).Should(Equal(false))
				})
			})
		})
		Describe("Emit", func() {
			Context("when emit takes a correct event", func() {
				It("should invoke handler", func() {
					invoked := false
					handler := func(event push.DeployFailureEvent) error {
						invoked = true
						return nil
					}
					binding := push.NewDeployFailureEventBinding(handler)
					event := push.DeployFailureEvent{}
					binding.Emit(event)

					Expect(invoked).Should(Equal(true))
				})
			})
			Context("when emit takes incorrect event", func() {
				It("should return error", func() {
					invoked := false
					handler := func(event push.DeployFailureEvent) error {
						invoked = true
						return nil
					}
					binding := push.NewDeployFailureEventBinding(handler)
					event := interfaces.Event{}
					err := binding.Emit(event)

					Expect(invoked).Should(Equal(false))
					Expect(err).ShouldNot(BeNil())
					Expect(err.Error()).Should(Equal("invalid event type"))
				})
			})
		})
	})

	Describe("PushStartedEvent", func() {
		Describe("Accept", func() {
			Context("when accept takes a correct event", func() {
				It("should return true", func() {
					binding := push.NewPushStartedEventBinding(nil)

					event := push.PushStartedEvent{}
					Expect(binding.Accepts(event)).Should(Equal(true))
				})
			})
			Context("when accept takes incorrect event", func() {
				It("should return false", func() {
					binding := push.NewPushStartedEventBinding(nil)

					event := interfaces.Event{}
					Expect(binding.Accepts(event)).Should(Equal(false))
				})
			})
		})
		Describe("Emit", func() {
			Context("when emit takes a correct event", func() {
				It("should invoke handler", func() {
					invoked := false
					handler := func(event push.PushStartedEvent) error {
						invoked = true
						return nil
					}
					binding := push.NewPushStartedEventBinding(handler)
					event := push.PushStartedEvent{}
					binding.Emit(event)

					Expect(invoked).Should(Equal(true))
				})
			})
			Context("when emit takes incorrect event", func() {
				It("should return error", func() {
					invoked := false
					handler := func(event push.PushStartedEvent) error {
						invoked = true
						return nil
					}
					binding := push.NewPushStartedEventBinding(handler)
					event := interfaces.Event{}
					err := binding.Emit(event)

					Expect(invoked).Should(Equal(false))
					Expect(err).ShouldNot(BeNil())
					Expect(err.Error()).Should(Equal("invalid event type"))
				})
			})
		})
	})

	Describe("PushFinishedEvent", func() {
		Describe("Accept", func() {
			Context("when accept takes a correct event", func() {
				It("should return true", func() {
					binding := push.NewPushFinishedEventBinding(nil)

					event := push.PushFinishedEvent{}
					Expect(binding.Accepts(event)).Should(Equal(true))
				})
			})
			Context("when accept takes incorrect event", func() {
				It("should return false", func() {
					binding := push.NewPushFinishedEventBinding(nil)

					event := interfaces.Event{}
					Expect(binding.Accepts(event)).Should(Equal(false))
				})
			})
		})
		Describe("Emit", func() {
			Context("when emit takes a correct event", func() {
				It("should invoke handler", func() {
					invoked := false
					handler := func(event push.PushFinishedEvent) error {
						invoked = true
						return nil
					}
					binding := push.NewPushFinishedEventBinding(handler)
					event := push.PushFinishedEvent{}
					binding.Emit(event)

					Expect(invoked).Should(Equal(true))
				})
			})
			Context("when emit takes incorrect event", func() {
				It("should return error", func() {
					invoked := false
					handler := func(event push.PushFinishedEvent) error {
						invoked = true
						return nil
					}
					binding := push.NewPushFinishedEventBinding(handler)
					event := interfaces.Event{}
					err := binding.Emit(event)

					Expect(invoked).Should(Equal(false))
					Expect(err).ShouldNot(BeNil())
					Expect(err.Error()).Should(Equal("invalid event type"))
				})
			})
		})
	})

	Describe("ArtifactRetrievalStartEvent", func() {
		Describe("Accept", func() {
			Context("when accept takes a correct event", func() {
				It("should return true", func() {
					binding := push.NewArtifactRetrievalStartEventBinding(nil)

					event := push.ArtifactRetrievalStartEvent{}
					Expect(binding.Accepts(event)).Should(Equal(true))
				})
			})
			Context("when accept takes incorrect event", func() {
				It("should return false", func() {
					binding := push.NewArtifactRetrievalStartEventBinding(nil)

					event := interfaces.Event{}
					Expect(binding.Accepts(event)).Should(Equal(false))
				})
			})
		})
		Describe("Emit", func() {
			Context("when emit takes a correct event", func() {
				It("should invoke handler", func() {
					invoked := false
					handler := func(event push.ArtifactRetrievalStartEvent) error {
						invoked = true
						return nil
					}
					binding := push.NewArtifactRetrievalStartEventBinding(handler)
					event := push.ArtifactRetrievalStartEvent{}
					binding.Emit(event)

					Expect(invoked).Should(Equal(true))
				})
			})
			Context("when emit takes incorrect event", func() {
				It("should return error", func() {
					invoked := false
					handler := func(event push.ArtifactRetrievalStartEvent) error {
						invoked = true
						return nil
					}
					binding := push.NewArtifactRetrievalStartEventBinding(handler)
					event := interfaces.Event{}
					err := binding.Emit(event)

					Expect(invoked).Should(Equal(false))
					Expect(err).ShouldNot(BeNil())
					Expect(err.Error()).Should(Equal("invalid event type"))
				})
			})
		})
	})

	Describe("ArtifactRetrievalFailureEvent", func() {
		Describe("Accept", func() {
			Context("when accept takes a correct event", func() {
				It("should return true", func() {
					binding := push.NewArtifactRetrievalFailureEventBinding(nil)

					event := push.ArtifactRetrievalFailureEvent{}
					Expect(binding.Accepts(event)).Should(Equal(true))
				})
			})
			Context("when accept takes incorrect event", func() {
				It("should return false", func() {
					binding := push.NewArtifactRetrievalFailureEventBinding(nil)

					event := interfaces.Event{}
					Expect(binding.Accepts(event)).Should(Equal(false))
				})
			})
		})
		Describe("Emit", func() {
			Context("when emit takes a correct event", func() {
				It("should invoke handler", func() {
					invoked := false
					handler := func(event push.ArtifactRetrievalFailureEvent) error {
						invoked = true
						return nil
					}
					binding := push.NewArtifactRetrievalFailureEventBinding(handler)
					event := push.ArtifactRetrievalFailureEvent{}
					binding.Emit(event)

					Expect(invoked).Should(Equal(true))
				})
			})
			Context("when emit takes incorrect event", func() {
				It("should return error", func() {
					invoked := false
					handler := func(event push.ArtifactRetrievalFailureEvent) error {
						invoked = true
						return nil
					}
					binding := push.NewArtifactRetrievalFailureEventBinding(handler)
					event := interfaces.Event{}
					err := binding.Emit(event)

					Expect(invoked).Should(Equal(false))
					Expect(err).ShouldNot(BeNil())
					Expect(err.Error()).Should(Equal("invalid event type"))
				})
			})
		})
	})

	Describe("ArtifactRetrievalSuccessEvent", func() {
		Describe("Accept", func() {
			Context("when accept takes a correct event", func() {
				It("should return true", func() {
					binding := push.NewArtifactRetrievalSuccessEventBinding(nil)

					event := push.ArtifactRetrievalSuccessEvent{}
					Expect(binding.Accepts(event)).Should(Equal(true))
				})
			})
			Context("when accept takes incorrect event", func() {
				It("should return false", func() {
					binding := push.NewArtifactRetrievalSuccessEventBinding(nil)

					event := interfaces.Event{}
					Expect(binding.Accepts(event)).Should(Equal(false))
				})
			})
		})
		Describe("Emit", func() {
			Context("when emit takes a correct event", func() {
				It("should invoke handler", func() {
					invoked := false
					handler := func(event push.ArtifactRetrievalSuccessEvent) error {
						invoked = true
						return nil
					}
					binding := push.NewArtifactRetrievalSuccessEventBinding(handler)
					event := push.ArtifactRetrievalSuccessEvent{}
					binding.Emit(event)

					Expect(invoked).Should(Equal(true))
				})
			})
			Context("when emit takes incorrect event", func() {
				It("should return error", func() {
					invoked := false
					handler := func(event push.ArtifactRetrievalSuccessEvent) error {
						invoked = true
						return nil
					}
					binding := push.NewArtifactRetrievalSuccessEventBinding(handler)
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
