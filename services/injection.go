package services

type InvokerInjection func(c *Container) (*Item, Invoker)
type ServiceInjection func(c *Container) (*Item, interface{})
type ParameterInjection func() (*Item, interface{})

var (
	invokerInjections   []InvokerInjection
	serviceInjections   []ServiceInjection
	parameterInjections []ParameterInjection
	initializers        []func(c *Container)
)

var (
	injectInvoker   = func(injection InvokerInjection) { invokerInjections = append(invokerInjections, injection) }
	injectService   = func(injection ServiceInjection) { serviceInjections = append(serviceInjections, injection) }
	injectParameter = func(injection ParameterInjection) { parameterInjections = append(parameterInjections, injection) }
	addInitializer  = func(initializer func(c *Container)) { initializers = append(initializers, initializer) }
)

func InjectInvoker(injection InvokerInjection) {
	injectInvoker(injection)
}

func InjectService(injection ServiceInjection) {
	injectService(injection)
}

func InjectParameter(injection ParameterInjection) {
	injectParameter(injection)
}

func AddInitializer(initializer func(c *Container)) {
	addInitializer(initializer)
}
