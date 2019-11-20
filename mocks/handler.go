package mocks

import I "github.com/compozed/deployadactyl/interfaces"

// Handler handmade mock for tests.
type Handler struct {
	OnEventCall struct {
		Received struct {
			Event I.Event
		}
		Returns struct {
			Error error
		}
	}
}

// OnEvent mock method.
func (h *Handler) OnEvent(event I.Event) error {
	h.OnEventCall.Received.Event = event

	return h.OnEventCall.Returns.Error
}
