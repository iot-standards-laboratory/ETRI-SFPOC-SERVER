package watcher

type DeviceWatcher struct {
	// did            string
	subscriber []chan map[string]interface{}
}

func (dw *DeviceWatcher) Notify(e Event) {
	dw.onChanged(e)
}

func (dw *DeviceWatcher) onChanged(e Event) {
	for _, s := range dw.subscriber {
		s <- e.Param()
	}
}

func (dw *DeviceWatcher) Subscribe(s chan map[string]interface{}) {
	dw.subscriber = append(dw.subscriber, s)
}

func (dw *DeviceWatcher) Desubscribe(s chan map[string]interface{}) int {
	length := len(dw.subscriber)
	for i, e := range dw.subscriber {
		if s == e {
			dw.subscriber[i] = dw.subscriber[length-1]
			dw.subscriber = dw.subscriber[:length-1]
			return length - 1
		}
	}
	return length
}
