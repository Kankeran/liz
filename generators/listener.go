package generators

import (
	"liz/elements"
)

// Listener generator struct
type Listener struct{}

// Generate generates listeners code
func (l *Listener) Generate(listeners []*elements.Listener) string {
	var code = "func(d *event.Data) {\n"
	for _, listener := range listeners {
		code += listener.ServiceGetter + "." + listener.Method + "(d)\n"
	}
	code += "}"
	return code
}
