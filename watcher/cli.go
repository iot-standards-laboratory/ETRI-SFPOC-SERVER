package watcher

func NewStateChangedEvent(t interface{}, param map[string]interface{}) *StatusChangedEvent {
	return &StatusChangedEvent{t: t, param: param}
}
