package router

import (
	"encoding/json"
	"etrismartfarmpoc/model"
	"etrismartfarmpoc/watcher"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var waitPermission = map[string]chan bool{}
var mutex sync.Mutex
var discoveredDevices = []*model.Device{}

func GetDiscoveredDevices(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	encoder := json.NewEncoder(w)
	encoder.Encode(discoveredDevices)
}

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

	if !db.IsExistDevice(device.DName) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Already exist device"))
		return
	}

	// 장치 ID 생성 및 탐색된 장치 추가
	device.DID = uuid.NewString()
	mutex.Lock()
	waitPermission[device.DID] = make(chan bool)
	discoveredDevices = append(discoveredDevices, device)

	// 관리자에게 탐색을 알림
	sendNotification(&Notification{Msg: "Add discovered device"})
	mutex.Unlock()

	timer := time.NewTimer(20 * time.Second)

	select {
	case <-r.Context().Done():
		fmt.Println("^^")
		mutex.Lock()
		defer mutex.Unlock()
		delete(waitPermission, device.DID)
		removeDevice(device)
		sendNotification(&Notification{Msg: "Remove discovered device"})
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("This operation is not permitted"))
		return
	case <-timer.C:
		mutex.Lock()
		defer mutex.Unlock()
		delete(waitPermission, device.DID)
		removeDevice(device)
		sendNotification(&Notification{Msg: "Remove discovered device"})
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
			sendNotification(&Notification{Msg: "Add device"})
		} else {
			mutex.Lock()
			defer mutex.Unlock()
			delete(waitPermission, device.DID)
			removeDevice(device)
			sendNotification(&Notification{Msg: "Remove discovered device"})
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("This operation is not permitted"))
		}
	}
}

func DeleteDevice(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	msg := map[string]string{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&msg)

	device, err := db.QueryDevice(msg["dname"])
	db.DeleteDevice(device)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}

	sendNotification(&Notification{Msg: "Delete device"})
	w.WriteHeader(http.StatusOK)
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

// 값 변경 알림
func PutDeviceStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	did, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("not found did"))
		return
	}

	if r.ContentLength == 0 {
		w.WriteHeader(http.StatusOK)
		return
	}

	watch, ok := watchers[did]
	if !ok {
		w.WriteHeader(http.StatusOK)
		return
	}

	param := map[string]interface{}{}
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&param)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		return
	}

	watch.Notify(watcher.NewStateChangedEvent(nil, param))
	w.WriteHeader(http.StatusOK)
}

// func PostDeviceStatus(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	did, ok := vars["id"]
// 	if !ok {
// 		w.WriteHeader(http.StatusBadRequest)
// 		w.Write([]byte("not found did"))
// 		return
// 	}

// 	if r.ContentLength == 0 {
// 		w.WriteHeader(http.StatusOK)
// 		return
// 	}

// 	watch, ok := watchers[did]
// 	if !ok {
// 		w.WriteHeader(http.StatusOK)
// 		return
// 	}

// 	param := map[string]interface{}{}
// 	decoder := json.NewDecoder(r.Body)

// 	err := decoder.Decode(&param)
// 	if err != nil {
// 		w.WriteHeader(http.StatusOK)
// 		return
// 	}

// 	watch.Notify(watcher.NewStateChangedEvent(nil, param))
// }
