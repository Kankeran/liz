package generators

import (
	"fmt"

	"liz/elements"
)

// Service generator struct
type Service struct{}

// Generate builds code
func (s *Service) Generate(service *elements.Service) string {
	var code = "func() interface{}{\n"
	if len(service.Constructor) != 0 {
		code += "service := " + service.Constructor + "("
		code += s.addArguments(service.Arguments)
		code += ")\n"
	} else {
		code += "service := &" + service.StructName + "{"
		code += s.addArguments(service.Arguments)
		code += "}\n"
	}

	for _, val := range service.Calls {
		code += "service." + val.Method + "("
		code += s.addArguments(val.Arguments)
		code += ")\n"
	}
	code += "\nreturn " + fmt.Sprint(service.Returns) + "\n}"

	return code
}

func (s *Service) addArguments(arguments []interface{}) (code string) {
	if len(arguments) > 0 {
		code += "\n"
	}
	for _, val := range arguments {
		code += fmt.Sprint(val) + ",\n"
	}

	return code
}
