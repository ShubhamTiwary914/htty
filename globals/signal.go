package globals 

/*
	event/signalling system similar to godot game engine 
	(ex - https://docs.godotengine.org/en/stable/getting_started/step_by_step/signals.html)

	for events across panels, register channels here and listeners for panels to emit/listen to these channels
	examples:
		- side panel load/unloading states of panel I/O
*/

import "sync"

type EventBus struct {
	Locker        sync.RWMutex
	Listeners map[string][]chan any
}

//register a channel to receive data when events come in from publishers onto this event
func (bus *EventBus) Subscribe(event string, channel chan any) {
	bus.Locker.Lock()
	defer bus.Locker.Unlock()
	bus.Listeners[event] = append(bus.Listeners[event], channel)
}

//when publish to this event, all subscribers of this event get the data
func (bus *EventBus) Publish(event string, data any){
	bus.Locker.Lock()
	listeners := bus.Listeners[event]
	bus.Locker.RUnlock()

	//pass data onto the subscribers' registered channels
	for _, channel := range listeners {
		select {
			case channel <- data:
			default: //for registered subs, not ready channels tho
		}
	}
}

//region: event-buses objects
//INFO: registered buses (while single bus can be used for many event, better to seperate em')
//so these event-buses are simply grouping of similar events 

//for filetree load/unloads events 
var StateBus = &EventBus{ Listeners: make(map[string][]chan any) }
const (
	EVENT_STATE_LOAD string = "state.load"
	EVENT_STATE_UNLOAD string = "state.unload"
)
