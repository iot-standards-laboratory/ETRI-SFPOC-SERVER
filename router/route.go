package router

import (
	"encoding/json"
	"etrismartfarmpoc/model"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

type Notification struct {
	Msg string `json:"msg"`
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
	mux.HandleFunc("/devices", DeleteDevice).Methods("DELETE")
	mux.HandleFunc("/devices", PutDevice).Methods("PUT", "OPTIONS")
	mux.HandleFunc("/device/{id}", GetDeviceWatch).Methods("GET")
	mux.HandleFunc("/device/{id}", PutDeviceStatus).Methods("PUT")
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

// func mapService(urls []string) {

// }
