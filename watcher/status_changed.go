package watcher

type StatusChangedEvent struct {
	Event
	t     interface{}
	param map[string]interface{}
}

func (e *StatusChangedEvent) Type() interface{} {
	return e.t
}

func (e *StatusChangedEvent) Param() map[string]interface{} {
	return e.param
}

// type StatusChangedListener struct {
// 	h func(e Event)
// }

// func (l *StatusChangedListener) Handle(e Event) {
// 	if l.h != nil {
// 		l.h(e)
// 	}
// }

// func (l *StatusChangedListener) AddSubscriber(s chan map[string]interface{}) {
// 	l.subscriber = append(l.subscriber, s)
// }

// func (l *StatusChangedListener) RemoveSubscriber(s chan map[string]interface{}) int {
// 	length := len(l.subscriber)
// 	for i, e := range l.subscriber {
// 		if s == e {
// 			l.subscriber[i] = l.subscriber[length-1]
// 			l.subscriber = l.subscriber[:length-1]
// 			return length - 1
// 		}
// 	}

// 	return length
// }

// func (l *StatusChangedListener) onClose() {}
