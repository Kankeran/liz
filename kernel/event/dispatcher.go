package event

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

func (d *dispatcher) add(eventName string, listener func(*Data)) {
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
	if _, ok := d.listeners[name]; !ok {
		return
	}
	d.eventChannel <- &dispatchData{&Data{name, data}, nil}
}

func (d *dispatcher) dispatchSync(name string, data interface{}) {
	if _, ok := d.listeners[name]; !ok {
		return
	}
	waitGroup := new(sync.WaitGroup)
	waitGroup.Add(1)
	d.eventChannel <- &dispatchData{&Data{name, data}, waitGroup}
	waitGroup.Wait()
}
