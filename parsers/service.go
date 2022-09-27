package parsers

import (
	"os"
	"regexp"
	"strings"
)

// Service used to parse service properties
type Service struct {
	servicesMap map[interface{}]interface{}
}

// SetOriginalServicesMap set source map of yaml to getting original data
func (s *Service) SetOriginalServicesMap(servicesMap map[interface{}]interface{}) {
	s.servicesMap = servicesMap
}

// Parse parses service properties
func (s *Service) Parse(servicesMap interface{}) interface{} {
	switch services := servicesMap.(type) {
	case map[interface{}]interface{}:
		servicesMap = s.parseMap(services)
	case []interface{}:
		servicesMap = s.parseInterface(services)
	}

	return servicesMap
}

func (s *Service) parseInterface(services []interface{}) interface{} {
	for key, value := range services {
		switch v := value.(type) {
		case string:
			services[key] = s.replaceConjunction(s.replaceSelfReference(s.replaceService(s.replaceEnv(v))))
		case map[interface{}]interface{}:
			services[key] = s.parseMap(v)
		case []interface{}:
			services[key] = s.parseInterface(v)
		}
	}

	return services
}

func (s *Service) parseMap(services map[interface{}]interface{}) interface{} {
	for key := range services {
		switch k := key.(type) {
		case string:
			oldKey := key
			key = s.replaceConjunction(s.replaceEnv(k))
			if oldKey != key {
				services[key] = services[oldKey]
				delete(services, oldKey)
			}
		}
	}

	for key, value := range services {
		switch v := value.(type) {
		case string:
			services[key] = s.replaceConjunction(s.replaceSelfReference(s.replaceService(s.replaceEnv(v))))
		case map[interface{}]interface{}:
			services[key] = s.parseMap(v)
		case []interface{}:
			services[key] = s.parseInterface(v)
		}
	}

	return services
}

func (s *Service) replaceService(arg string) string {
	serviceMatcher := regexp.MustCompile(`@\(\S+\)`)
	var replaced string
	var lastId = 0

	for _, index := range serviceMatcher.FindAllStringIndex(arg, -1) {
		replaced += arg[lastId:index[0]]
		replaced += strings.Replace(strings.Replace(arg[index[0]:index[1]], "@(", "container.Get(\"", 1), ")", "\")", 1)
		replaced += ".(*" + s.getStructName(strings.TrimRight(strings.TrimLeft(arg[index[0]:index[1]], "@("), ")")) + ")"
		lastId = index[1]
	}

	if len(replaced) > 0 {
		replaced += arg[lastId:]
		return replaced
	}

	serviceMatcher = regexp.MustCompile(`@\S+`)
	for _, index := range serviceMatcher.FindAllStringIndex(arg, -1) {
		replaced = strings.Replace(arg[index[0]:index[1]], "@", "container.Get(\"", 1) + "\")"
		replaced += ".(*" + s.getStructName(strings.TrimLeft(arg[index[0]:index[1]], "@")) + ")"
		break
	}

	if len(replaced) > 0 {
		return replaced
	}

	return arg
}

func (s *Service) replaceEnv(arg string) string {
	envSettingMatcher := regexp.MustCompile(`\$env\(\S+\)`)
	var replaced string
	var lastId = 0

	for _, index := range envSettingMatcher.FindAllStringIndex(arg, -1) {
		replaced += arg[lastId:index[0]]
		replaced += os.Getenv(strings.Replace(strings.Replace(arg[index[0]:index[1]], "$env(", "", 1), ")", "", 1))
		lastId = index[1]
	}

	if len(replaced) > 0 {
		replaced += arg[lastId:]
		return replaced
	}

	return arg
}

func (s *Service) replaceConjunction(arg string) string {
	return strings.ReplaceAll(strings.ReplaceAll(arg, "'", "\""), "~", "+")
}

func (s *Service) replaceSelfReference(arg string) string {
	return strings.ReplaceAll(arg, "$this.", "service.")
}

func (s *Service) getStructName(serviceName string) string {
	serviceData, ok := s.servicesMap[serviceName].(map[interface{}]interface{})
	if !ok {
		panic("service " + serviceName + " doesn't exists")
	}
	if name, ok := serviceData["struct"]; ok {
		return name.(string)
	}

	panic("Cannot get struct name from service " + serviceName + " reason: field struct is empty")
}
