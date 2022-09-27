package services

type Invoker func() interface{}
type Item int64

var ContainerId = new(Item)

type Container struct {
	serviceInvokers []Invoker
	services        []interface{}
	parameters      []interface{}
	consumerAssert  func()
}

func NewContainer() *Container {
	return &Container{
		serviceInvokers: []Invoker{},
		services:        []interface{}{},
		parameters:      []interface{}{},
		consumerAssert:  needToConsume,
	}
}

func (c *Container) GetService(serviceId *Item) interface{} {
	if *serviceId < 0 {
		itemId := -(*serviceId)
		if itemId > Item(len(c.serviceInvokers)) {
			panic("Bad service invoker Id")
		}

		service := c.serviceInvokers[itemId-1]()
		*serviceId = Item(len(c.services) + 1)
		c.services = append(c.services, service)

		return service
	}

	if *serviceId > Item(len(c.services)) {
		panic("Bad service Id")
	}

	return c.services[*serviceId-1]
}

func (c *Container) AddServiceInvoker(serviceId *Item, serviceInvoker Invoker) {
	*serviceId = -Item(len(c.serviceInvokers) + 1)
	c.serviceInvokers = append(c.serviceInvokers, serviceInvoker)
}

func (c *Container) AddService(serviceId *Item, service interface{}) {
	*serviceId = Item(len(c.services) + 1)
	c.services = append(c.services, service)
}

func (c *Container) GetParameter(parameterId *Item) interface{} {
	if *parameterId > Item(len(c.parameters)) {
		panic("Bad parameter Id")
	}

	return c.parameters[*parameterId-1]
}

func (c *Container) AddParameter(parameterId *Item, value interface{}) {
	*parameterId = Item(len(c.parameters) + 1)
	c.parameters = append(c.parameters, value)
}

func (c *Container) Consume() {
	c.consumerAssert()

	c.AddService(ContainerId, c)

	for _, injection := range parameterInjections {
		c.AddParameter(injection())
	}

	for _, injection := range invokerInjections {
		c.AddServiceInvoker(injection(c))
	}

	for _, injection := range serviceInjections {
		c.AddService(injection(c))
	}

	for _, initializer := range initializers {
		initializer(c)
	}

	parameterInjections = nil
	invokerInjections = nil
	serviceInjections = nil
	initializers = nil

	injectInvoker = c.injectInvoker
	injectService = c.injectService
	injectParameter = c.injectParameter
	addInitializer = func(initializer func(c *Container)) { panic("initializer cannot be added at runtime") }

	c.consumerAssert = alreadyConsumed
}

func alreadyConsumed() {
	panic("injections can be consumed only once, at runtime it's consuming all automatically")
}

func needToConsume() {}

func (c *Container) injectInvoker(injection InvokerInjection) {
	c.AddServiceInvoker(injection(c))
}

func (c *Container) injectService(injection ServiceInjection) {
	c.AddService(injection(c))
}

func (c *Container) injectParameter(injection ParameterInjection) {
	c.AddParameter(injection())
}
