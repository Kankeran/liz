package services

import (
	"testing"
)

type testService struct {
	asdf string
}

var (
	test1 = new(Item)
	test2 = new(Item)
	test3 = new(Item)
)

func someInvokerInjection(c *Container) (*Item, Invoker) {
	return test1, func() interface{} {
		return &testService{"nightmare"}
	}
}

func someServiceInjection(c *Container) (*Item, interface{}) {
	return test2, &testService{"nightmare"}
}

func someParameterInjection() (*Item, interface{}) {
	return test3, "dupa"
}

func TestAddInjection(t *testing.T) {
	InjectInvoker(someInvokerInjection)

	if len(invokerInjections) != 1 {
		t.Error("one invoker injection should be added")
	}

	InjectService(someServiceInjection)

	if len(serviceInjections) != 1 {
		t.Error("one service injection should be added")
	}

	InjectParameter(someParameterInjection)

	if len(parameterInjections) != 1 {
		t.Error("one parameter injection should be added")
	}

	AddInitializer(func(c *Container) {})

	if len(initializers) != 1 {
		t.Error("one initializer should be added")
	}
}

func TestNewContainer(t *testing.T) {
	var container = NewContainer()

	if container == nil || container.serviceInvokers == nil || container.services == nil {
		t.Error("incorrect container initialization")
	}
}

func TestContainerConsume(t *testing.T) {
	invokerInjections = nil
	invokerInjections = append(invokerInjections, InvokerInjection(someInvokerInjection))
	var container = NewContainer()
	container.Consume()

	if len(container.serviceInvokers) != 1 {
		t.Error("one invoker sould be added")
	}
}

func TestGetService(t *testing.T) {
	var container = NewContainer()
	container.AddServiceInvoker(someInvokerInjection(container))

	if len(container.serviceInvokers) != 1 {
		t.Error("one invoker sould be added")
	}

	var service, ok = container.GetService(test1).(*testService)

	if !ok || service.asdf != "nightmare" {
		t.Error("service should be the same as added by injection")
	}

	var service2 *testService

	service2, ok = container.GetService(test1).(*testService)

	if !ok || service2.asdf != "nightmare" || service != service2 {
		t.Error("service should be the same as previously added")
	}
}

func TestGetContainer(t *testing.T) {
	var container = NewContainer()
	container.Consume()
	var c, ok = container.GetService(ContainerId).(*Container)

	if !ok || container != c {
		t.Error("getting container service should return container")
	}
}

func TestGetParameter(t *testing.T) {
	var container = NewContainer()
	container.AddParameter(someParameterInjection())

	if len(container.parameters) != 1 {
		t.Error("one parameter sould be added")
	}

	var parameter, ok = container.GetParameter(test3).(string)

	if !ok || parameter != "dupa" {
		t.Error("parameter \"test\" should have value \"dupa\"")
	}
}

func TestInjectionAfterConsume(t *testing.T) {
	var container = NewContainer()
	container.Consume()
	InjectInvoker(someInvokerInjection)

	if len(invokerInjections) != 0 {
		t.Error("no one invoker injection should left after consume")
	}

	if len(container.serviceInvokers) != 1 {
		t.Error("one invoker should be added to the container")
	}

	InjectService(someServiceInjection)

	if len(serviceInjections) != 0 {
		t.Error("no one service injection should be left after consume")
	}

	if len(container.serviceInvokers) != 1 {
		t.Error("one service should be added to the container")
	}

	InjectParameter(someParameterInjection)

	if len(parameterInjections) != 0 {
		t.Error("no one parameter injection should be left after consume")
	}

	if len(container.serviceInvokers) != 1 {
		t.Error("one parameter should be added to the container")
	}
}

func TestSecondConsumeShouldPanic(t *testing.T) {
	var container = NewContainer()
	container.Consume()
	assertPanic(t, container.Consume)
}

func TestAddInitializerAfterConsumeShouldPanic(t *testing.T) {
	var container = NewContainer()
	container.Consume()
	assertPanic(t, func() { AddInitializer(func(c *Container) {}) })
}

func assertPanic(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	f()
}
