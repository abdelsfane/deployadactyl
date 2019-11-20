package mocks

type EventBinding struct {
	AcceptsCall struct {
		Received struct {
			Event interface{}
		}
		Returns struct {
			Bool bool
		}
	}
	EmitCall struct {
		Received struct {
			Event interface{}
		}
		Called struct {
			Bool bool
		}
		Returns struct {
			Error error
		}
		ShouldPanic bool
	}
}

func (b *EventBinding) Accepts(event interface{}) bool {
	//return reflect.TypeOf(event) == reflect.TypeOf(b)
	b.AcceptsCall.Received.Event = event

	return b.AcceptsCall.Returns.Bool
}

func (b *EventBinding) Emit(gevent interface{}) error {
	b.EmitCall.Called.Bool = true
	b.EmitCall.Received.Event = gevent

	if b.EmitCall.ShouldPanic {
		panic("You messed up")
	}

	return b.EmitCall.Returns.Error

}
