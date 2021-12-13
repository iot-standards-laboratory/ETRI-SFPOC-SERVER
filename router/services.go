package router

import (
	"encoding/json"
	"etrismartfarmpoc/containermgmt"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func GetServices(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	sname := r.Header.Get("sname")

	if len(sname) != 0 {
		sid, err := db.GetSID(sname)
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

	log.Println("Route to", host+url)
	req, err := http.NewRequest(r.Method, "http://"+host+url, r.Body)
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
