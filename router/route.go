package router

import (
	"encoding/json"
	"etrismartfarmpoc/containermgmt"
	"etrismartfarmpoc/model"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/urfave/negroni"
)

type Notification struct {
	Msg string `json:msg`
}

var db model.DBHandler

func init() {
	var err error
	db, err = model.NewDBHandler("postgres", "dump.db")
	if err != nil {
		panic(err)
	}
}

// var rd *render.Render = render.New()
func NewRouter() http.Handler {
	mux := mux.NewRouter()
	mux.HandleFunc("/noti", GetNotification).Methods("GET")
	mux.HandleFunc("/controllers", GetControllerList).Methods("GET")
	mux.HandleFunc("/controllers", PostController).Methods("POST")
	mux.HandleFunc("/services", GetServices).Methods("GET")
	mux.HandleFunc("/services", PostService).Methods("POST")
	mux.HandleFunc("/services", PutServices).Methods("PUT")
	mux.HandleFunc("/devices", PostDevice).Methods("POST")
	mux.HandleFunc("/devices", GetDevices).Methods("GET")
	mux.HandleFunc("/devices", PutDevice).Methods("PUT", "OPTIONS")
	mux.HandleFunc("/discover", GetDiscoveredDevices).Methods("GET")
	mux.PathPrefix("/services/").HandlerFunc(RouteRequestToService)

	n := negroni.Classic() // 파일 서버 및 로그기능을 제공함
	n.UseHandler(mux)

	return n
}

func EchoRoute(w http.ResponseWriter, r *http.Request) {
	// w.Write([]byte(r.RequestURI))
	// token := strings.Split(r.RequestURI, "/")

	l := len([]rune("/services/"))
	idx := strings.Index(r.RequestURI[l:], "/")
	id := r.RequestURI[l : l+idx]
	url := r.RequestURI[l+idx:]

	m := make(map[string]string)
	m["id"] = id
	m["url"] = url

	encoder := json.NewEncoder(w)
	encoder.Encode(m)
}

func GetControllerList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	list, err := db.GetControllers()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	fmt.Println(list)
	encoder := json.NewEncoder(w)

	err = encoder.Encode(list)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
}

func PostController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	controller, err := db.AddController(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusCreated)
	encoder := json.NewEncoder(w)
	err = encoder.Encode(controller)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	sendNotification(&Notification{Msg: "Added Controller"})
}

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

func GetDiscoveredDevices(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	devices, _, err := db.GetDevices()

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(devices)
}

// func GetDiscoveredDevices(w http.ResponseWriter, r *http.Request) {
// 	noti := make(chan string, 1)
// 	mutex.Lock()
// 	discoveredNotifications = append(discoveredNotifications, noti)
// 	mutex.Unlock()
// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 		w.Write([]byte(err.Error()))
// 		return
// 	}

// 	for {
// 		<-noti
// 		// fmt.Println("Write!!", discoveredDevices)
// 		if conn.WriteJSON(discoveredDevices) != nil {
// 			log.Println(err)
// 			return
// 		}
// 	}
// }

var waitPermission = map[string]chan bool{}
var mutex sync.Mutex
var discoveredDevices []*model.Device
var discoveredNotifications []chan string

func removeDevice(device *model.Device) {
	for i, e := range discoveredDevices {
		if e == device {
			discoveredDevices[i] = discoveredDevices[len(discoveredDevices)-1]
			discoveredDevices = discoveredDevices[:len(discoveredDevices)-1]
		}
	}
}

func GetDevices(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	devices, _, err := db.GetDevices()

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(devices)
}

// PostDevice : Handle for Receiving discovery message
func PostDevice(w http.ResponseWriter, r *http.Request) {
	var device = &model.Device{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(device)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 등록된 제어기로 부터 전송된 요청 메시지임을 확인
	if !db.IsExistController(device.CID) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Wrong Controller ID"))
		return
	}

	// 장치 ID 생성 및 탐색된 장치 추가
	device.DID = uuid.NewString()
	mutex.Lock()
	waitPermission[device.DID] = make(chan bool)
	discoveredDevices = append(discoveredDevices, device)

	// 관리자에게 탐색을 알림
	for _, noti := range discoveredNotifications {
		noti <- device.DID
	}
	mutex.Unlock()

	timer := time.NewTimer(20 * time.Second)
	select {
	case <-r.Context().Done():
		fmt.Println("Done!!")
		mutex.Lock()
		defer mutex.Unlock()
		delete(waitPermission, device.DID)
		removeDevice(device)

		for _, noti := range discoveredNotifications {
			noti <- device.DID
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("This operation is not permitted"))
		return
	case <-timer.C:
		mutex.Lock()
		defer mutex.Unlock()
		delete(waitPermission, device.DID)
		removeDevice(device)

		for _, noti := range discoveredNotifications {
			noti <- device.DID
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("This operation is not permitted"))
		return
	case b := <-waitPermission[device.DID]:
		if b {
			// 디바이스 등록 절차 수행
			db.AddDevice(device)
			db.AddService(device.SName)

			// 디바이스 등록 알림
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(device)
			sendNotification(&Notification{Msg: "Added Device"})
			mutex.Lock()
			defer mutex.Unlock()
			delete(waitPermission, device.DID)
			removeDevice(device)

			for _, noti := range discoveredNotifications {
				noti <- device.DID
			}
		} else {
			mutex.Lock()
			defer mutex.Unlock()
			delete(waitPermission, device.DID)
			removeDevice(device)

			for _, noti := range discoveredNotifications {
				noti <- device.DID
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("This operation is not permitted"))
		}
	}
}

func PutDevice(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	msg := map[string]string{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&msg)

	fmt.Println("msg[did]: ", msg["did"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}

	ch, ok := waitPermission[msg["did"]]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Wrong Device ID is Sended"))
		return
	}
	ch <- true
	w.WriteHeader(http.StatusOK)
}

func GetServices(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	if r.ContentLength != 0 {
		obj := map[string]string{}
		err := json.NewDecoder(r.Body).Decode(&obj)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		sid, err := db.GetSID(obj["sname"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Write([]byte(sid))
		return
	}

	l, err := db.GetServices()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(l)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func PutServices(w http.ResponseWriter, r *http.Request) {
	var obj = map[string]string{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&obj)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	s, err := db.UpdateService(obj["name"], obj["addr"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sendNotification(&Notification{Msg: "Update Service"})
	err = json.NewEncoder(w).Encode(s)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

// func GetServiceList(w http.ResponseWriter, r *http.Request) {
// 	result := containermgmt.GetContainers(context.Background())
// 	b := result.Value(containermgmt.ReturnKey).([]byte)
// 	w.Write(b)
// }

func PostService(w http.ResponseWriter, r *http.Request) {

	var obj = map[string]string{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&obj)

	db.IsExistService(obj["name"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	err = containermgmt.CreateContainer(obj["name"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func RouteRequestToService(w http.ResponseWriter, r *http.Request) {
	l := len([]rune("/services/"))
	var id string
	var url string
	idx := strings.Index(r.RequestURI[l:], "/")
	if idx == -1 {
		id = r.RequestURI[l:]
		url = "/"
	} else {
		id = r.RequestURI[l : l+idx]
		if len([]rune(r.RequestURI)) <= l+idx {
			url = "/"
		} else {
			url = r.RequestURI[l+idx:]
		}
	}

	// w.Write([]byte(vars["url"]))
	host, err := db.GetAddr(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	log.Println("Route to", host+"/"+url)
	req, err := http.NewRequest(r.Method, "http://"+host+"/"+url, r.Body)
	if err != nil {
		// 잘못된 메시지 포맷이 전달된 경우
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		// 잘못된 메시지 포맷이 전달된 경우
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	io.Copy(w, resp.Body)
}

// func mapService(urls []string) {

// }
