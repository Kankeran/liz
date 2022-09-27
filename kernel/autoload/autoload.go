package autoload

import (
	"liz/kernel/event"
	"liz/kernel/services"
)

func init() {
	services.Build()
	event.RunListener()
}
