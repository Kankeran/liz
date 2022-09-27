package builder

import (
	"bytes"
	"fmt"
	"os/exec"

	"liz/domain"

	"github.com/go-git/go-git/v5"
)

const configYamlData = "services:\nlisteners:"
const serviceData = "package services\n // Build building container container\n func Build() {\n\n}"
const containerData = `package container

type containerStruct struct {
	services         map[string]interface{}
	servicesCreators map[string]func() interface{}
}

var containerInstance = &containerStruct{
	services:         make(map[string]interface{}),
	servicesCreators: make(map[string]func() interface{}),
}

// Get getting searched service instance
func Get(serviceName string) interface{} {
	service, ok := containerInstance.services[serviceName]
	if !ok {
		containerInstance.services[serviceName] = containerInstance.servicesCreators[serviceName]()
		service = containerInstance.services[serviceName]
	}

	return service
}

// Has check service exists
func Has(serviceName string) bool {
	_, ok := containerInstance.servicesCreators[serviceName]

	return ok
}

// Set sets function to invoking service with specified name
func Set(serviceName string, serviceCreator func() interface{}) {
	containerInstance.servicesCreators[serviceName] = serviceCreator
}
`

const dispatcherData = `package event

import "sync"

// Data holds dispatched event data
type Data struct {
	Name  string
	Value interface{}
}

type dispatchData struct {
	event     *Data
	waitGroup *sync.WaitGroup
}

type dispatcher struct {
	eventChannel chan *dispatchData
	listeners    map[string]func(*Data)
}

var dispatcherInstance = &dispatcher{
make(chan *dispatchData),
make(map[string]func(*Data)),
}

// Dispatch asynchronously call event by name
func Dispatch(name string, data interface{}) {
	dispatcherInstance.dispatch(name, data)
}

// DispatchSync synchronously call event by name
func DispatchSync(name string, data interface{}) {
	dispatcherInstance.dispatchSync(name, data)
}

// RunListener start listening of event call
func RunListener() {
	dispatcherInstance.run()
}

// Add event listener to dispatcher
func Add(eventName string, listener func(*Data)) {
	dispatcherInstance.add(eventName, listener)
}

func (d *dispatcher)add(eventName string, listener func(*Data)) {
	d.listeners[eventName] = listener
}

func (d *dispatcher) run() {
	go func() {
		for {
			go func(dat *dispatchData) {
				d.listeners[dat.event.Name](dat.event)
				if dat.waitGroup != nil {
					dat.waitGroup.Done()
				}
			}(<-d.eventChannel)
		}
	}()
}

func (d *dispatcher) dispatch(name string, data interface{}) {
	if _, ok := d.listeners[name]; !ok {return}
	d.eventChannel <- &dispatchData{&Data{name, data}, nil}
}

func (d *dispatcher) dispatchSync(name string, data interface{}) {
	if _, ok := d.listeners[name]; !ok {return}
	waitGroup := new(sync.WaitGroup)
	waitGroup.Add(1)
	d.eventChannel <- &dispatchData{&Data{name, data}, waitGroup}
	waitGroup.Wait()
}
`

const autoLoadData = `package autoload

func init() {
	services.Build()
	event.RunListener()
}`

const appData = `package main

import (
	_ "%s/kernel/autoload"
)

func main() {
	// write your code here :)
}
`

// ProjectStarter used to generate common files
type ProjectStarter struct {
	configYamlWriter, serviceWriter, containerWriter, dispatcherWriter, autoLoadWriter *domain.FileWriter
	codeFormatter                                                                      *domain.CodeFormatter
}

// NewProjectStarter initialize new ProjectStarter
func NewProjectStarter(
	configYamlWriter *domain.FileWriter,
	serviceWriter *domain.FileWriter,
	containerWriter *domain.FileWriter,
	dispatcherWriter *domain.FileWriter,
	autoLoadWriter *domain.FileWriter,
	codeFormatter *domain.CodeFormatter,
) *ProjectStarter {
	return &ProjectStarter{
		configYamlWriter: configYamlWriter,
		serviceWriter:    serviceWriter,
		containerWriter:  containerWriter,
		dispatcherWriter: dispatcherWriter,
		autoLoadWriter:   autoLoadWriter,
		codeFormatter:    codeFormatter,
	}
}

// Build generates project common files
func (ps *ProjectStarter) Build(projectName string) {
	var err error

	_, err = git.PlainInit("./", false)

	if err != nil && err != git.ErrRepositoryAlreadyExists {
		panic(err)
	}

	cmd := exec.Command("go", "mod", "init", projectName)
	var stdErr bytes.Buffer
	cmd.Stderr = &stdErr
	err = cmd.Start()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = cmd.Wait()
	if err != nil && stdErr.String() != "go mod init: go.mod already exists\n" {
		fmt.Println(err.Error())
		fmt.Println(stdErr.String())
		return
	}

	err = ps.configYamlWriter.Write([]byte(configYamlData))
	if err != nil {
		panic(err)
	}

	var output []byte

	output, err = ps.codeFormatter.Format(containerData)
	if err != nil {
		panic(err)
	}
	err = ps.containerWriter.Write(output)
	if err != nil {
		panic(err)
	}

	output, err = ps.codeFormatter.Format(dispatcherData)
	if err != nil {
		panic(err)
	}
	err = ps.dispatcherWriter.Write(output)
	if err != nil {
		panic(err)
	}

	output, err = ps.codeFormatter.Format(serviceData)
	if err != nil {
		panic(err)
	}
	err = ps.serviceWriter.Write(output)
	if err != nil {
		panic(err)
	}

	output, err = ps.codeFormatter.Format(autoLoadData)
	if err != nil {
		panic(err)
	}
	err = ps.autoLoadWriter.Write(output)
	if err != nil {
		panic(err)
	}

	appWriter := domain.NewFileWriter(".", projectName+".go")

	output, err = ps.codeFormatter.Format(fmt.Sprintf(appData, projectName))
	if err != nil {
		panic(err)
	}

	err = appWriter.Write(output)
	if err != nil {
		panic(err)
	}
}
