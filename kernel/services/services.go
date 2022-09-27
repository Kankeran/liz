package services

import (
	"liz/builder"
	"liz/domain"
	"liz/generators"
	"liz/kernel/container"
	"liz/parsers"
)

// Build building container container
func Build() {

	container.Set("code_formatter", func() interface{} {
		service := &domain.CodeFormatter{}

		return service
	})

	container.Set("listener_generator", func() interface{} {
		service := &generators.Listener{}

		return service
	})

	container.Set("service_parser", func() interface{} {
		service := &parsers.Service{}

		return service
	})

	container.Set("container_file_writer", func() interface{} {
		service := domain.NewFileWriter(
			"kernel/container",
			"container.go",
		)

		return service
	})

	container.Set("dispatcher_file_writer", func() interface{} {
		service := domain.NewFileWriter(
			"kernel/event",
			"dispatcher.go",
		)

		return service
	})

	container.Set("container_builder", func() interface{} {
		service := builder.NewContainerBuilder(
			container.Get("yaml_file_reader").(*parsers.YamlFileReader),
			container.Get("reference_parser").(*parsers.Reference),
			container.Get("service_parser").(*parsers.Service),
			container.Get("listener_generator").(*generators.Listener),
			container.Get("service_generator").(*generators.Service),
			container.Get("service_file_writer").(*domain.FileWriter),
			container.Get("code_formatter").(*domain.CodeFormatter),
		)

		return service
	})

	container.Set("service_generator", func() interface{} {
		service := &generators.Service{}

		return service
	})

	container.Set("service_file_writer", func() interface{} {
		service := domain.NewFileWriter(
			"kernel/services",
			"services.go",
		)

		return service
	})

	container.Set("autoload_file_writer", func() interface{} {
		service := domain.NewFileWriter(
			"kernel/autoload",
			"autoload.go",
		)

		return service
	})

	container.Set("project_starter_builder", func() interface{} {
		service := builder.NewProjectStarter(
			container.Get("config_yaml_file_writer").(*domain.FileWriter),
			container.Get("service_file_writer").(*domain.FileWriter),
			container.Get("container_file_writer").(*domain.FileWriter),
			container.Get("dispatcher_file_writer").(*domain.FileWriter),
			container.Get("autoload_file_writer").(*domain.FileWriter),
			container.Get("code_formatter").(*domain.CodeFormatter),
		)

		return service
	})

	container.Set("reference_parser", func() interface{} {
		service := parsers.NewReference(
			make(map[string]interface{}),
			container.Get("yaml_file_reader").(*parsers.YamlFileReader),
		)

		return service
	})

	container.Set("yaml_file_reader", func() interface{} {
		service := parsers.NewYamlFileReader(
			make(map[string]interface{}),
		)

		return service
	})

	container.Set("config_yaml_file_writer", func() interface{} {
		service := domain.NewFileWriter(
			"config",
			"services.yaml",
		)

		return service
	})

}
