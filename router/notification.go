package router

import (
	"etrismartfarmpoc/watcher"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var notifications []chan *Notification
var notiMutex sync.Mutex

func sendNotification(noti *Notification) {
	// notiMutex.Lock()
	// defer notiMutex.Unlock()
	for _, ch := range notifications {
		ch <- noti
	}
}

func removeNotification(noti chan *Notification) {
	notiMutex.Lock()
	defer notiMutex.Unlock()
	for i, e := range notifications {
		if e == noti {
			discoveredDevices[i] = discoveredDevices[len(discoveredDevices)-1]
			discoveredDevices = discoveredDevices[:len(discoveredDevices)-1]
		}
	}
}

func GetNotification(w http.ResponseWriter, r *http.Request) {
	notiChan := make(chan *Notification, 1)
	notiMutex.Lock()
	notifications = append(notifications, notiChan)
	notiMutex.Unlock()
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	for {
		notification := <-notiChan
		// fmt.Println("Write!!", discoveredDevices)
		if conn.WriteJSON(notification) != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			removeNotification(notiChan)
			notiChan = nil
			return
		}
	}
}

var watchers = map[string]*watcher.DeviceWatcher{}

func GetDeviceWatch(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	did, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("not found did"))
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	// payload := map[string]interface{}{}
	// decoder := json.NewDecoder(r.Body)
	// err = decoder.Decode(&payload)
	// if err != nil {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	w.Write([]byte(err.Error()))
	// 	return
	// }

	if !checkDid(did) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("not exist did"))
		return
	}

	s := make(chan map[string]interface{})
	ec := make(chan error)

	defer close(s)
	watch, ok := watchers[did]

	if !ok {
		watch = &watcher.DeviceWatcher{}
		watchers[did] = watch
	}

	watch.Subscribe(s)

	defer func() {
		cnt := watch.Desubscribe(s)
		if cnt == 0 {
			delete(watchers, did)
			log.Println("delete watchers: ", watchers)
		}
	}()

	go func() {
		_, _, err := conn.ReadMessage()
		if err != nil {
			ec <- err
		}
	}()
	for {
		select {
		case param := <-s:
			if conn.WriteJSON(param) != nil {
				log.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

		case <-r.Context().Done():
			fmt.Println("see you later~")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("See you later~"))
			return
		case <-ec:
			return
		}

	}
}

func checkDid(did string) bool {
	return true
}
