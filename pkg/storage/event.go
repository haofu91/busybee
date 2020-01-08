package storage

const (
	// InstanceLoadedEvent  load a new instance event
	InstanceLoadedEvent = iota
	// InstanceRemovedEvent instance removed to another node event
	InstanceRemovedEvent
	// InstanceStateLoadedEvent load a new instance state event
	InstanceStateLoadedEvent
	// InstanceStateRemovedEvent instance state removed to another node event
	InstanceStateRemovedEvent
)

// Event the event
type Event struct {
	EventType int
	Data      interface{}
}