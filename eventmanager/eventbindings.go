package eventmanager

import "github.com/compozed/deployadactyl/interfaces"

type EventBindings struct {
	bindings []interfaces.Binding
}

func (eb *EventBindings) AddBinding(b interfaces.Binding) {
	eb.bindings = append(eb.bindings, b)
}

func (eb EventBindings) GetBindings() []interfaces.Binding {
	return eb.bindings
}
