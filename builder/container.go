package builder

import (
	"liz/domain"
	"liz/elements"
	"liz/generators"
	"liz/parsers"
)

// Container builder struct
type Container struct {
	yamlFileReader    *parsers.YamlFileReader
	referenceParser   *parsers.Reference
	serviceParser     *parsers.Service
	listenerGenerator *generators.Listener
	serviceGenerator  *generators.Service
	fileWriter        *domain.FileWriter
	codeFormatter     *domain.CodeFormatter
}

// NewContainerBuilder initialize Container struct
func NewContainerBuilder(
	yamlFileReader *parsers.YamlFileReader,
	referenceParser *parsers.Reference,
	serviceParser *parsers.Service,
	listenerGenerator *generators.Listener,
	serviceGenerator *generators.Service,
	fileWriter *domain.FileWriter,
	codeFormatter *domain.CodeFormatter,
) *Container {
	return &Container{
		yamlFileReader:    yamlFileReader,
		referenceParser:   referenceParser,
		serviceParser:     serviceParser,
		listenerGenerator: listenerGenerator,
		serviceGenerator:  serviceGenerator,
		fileWriter:        fileWriter,
		codeFormatter:     codeFormatter,
	}
}

// Build builds service and listeners code
func (c *Container) Build() {
	servicesMap, err := c.yamlFileReader.Read("./config/services.yaml")
	if err != nil {
		panic(err)
	}

	servicesMap, err = c.referenceParser.Parse(servicesMap.(map[interface{}]interface{}), "./config/services.yaml")
	if err != nil {
		panic(err)
	}
	listenersMap := servicesMap.(map[interface{}]interface{})["listeners"]
	servicesMap = servicesMap.(map[interface{}]interface{})["services"]

	var code = "package services\n // Build building container container\n func Build() {\n\n"

	if servicesMap != nil {
		c.serviceParser.SetOriginalServicesMap(servicesMap.(map[interface{}]interface{}))
		servicesMap = c.serviceParser.Parse(servicesMap.(map[interface{}]interface{}))
		for serviceName, serviceMap := range servicesMap.(map[interface{}]interface{}) {
			code += "container.Set(\"" + serviceName.(string) + "\", " + c.serviceGenerator.Generate(elements.NewService(serviceMap.(map[interface{}]interface{}))) + ")\n\n"
		}
	}

	if listenersMap != nil {
		listenersMap = c.serviceParser.Parse(listenersMap.(map[interface{}]interface{}))
		for eventName, listenerMap := range listenersMap.(map[interface{}]interface{}) {
			var listeners = make([]*elements.Listener, len(listenerMap.([]interface{})))
			for key, listenerData := range listenerMap.([]interface{}) {
				listeners[key] = elements.NewListener(listenerData.(map[interface{}]interface{}))
			}
			code += "event.Add(\"" + eventName.(string) + "\", " + c.listenerGenerator.Generate(listeners) + ")\n\n"
		}
	}

	code += "}"

	var output []byte
	// println(code)
	output, err = c.codeFormatter.Format(code)
	if err != nil {
		panic(err)
	}

	err = c.fileWriter.Write(output)
	if err != nil {
		panic(err)
	}
}
