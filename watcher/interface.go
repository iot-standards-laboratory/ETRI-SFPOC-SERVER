package watcher

type Watcher interface {
}

type Event interface {
	Type() interface{}
	Param() map[string]interface{}
}

type EventListener interface {
	Handle(e Event)
}
