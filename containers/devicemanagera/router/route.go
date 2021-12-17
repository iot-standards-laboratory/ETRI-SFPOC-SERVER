package router

import (
	"bytes"
	"devicemanagerb/constants"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func NewRouter() http.Handler {
	mux := mux.NewRouter()
	// mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	w.WriteHeader(http.StatusOK)
	// 	w.Write([]byte("I am devicemanagerB"))
	// })

	mux.HandleFunc("/{id}", PutStatusChangedHandle).Methods("PUT")
	mux.HandleFunc("/{id}", PostStatusChangedHandle).Methods("POST")
	mux.HandleFunc("/{id}", GetStatusHandle).Methods("GET")

	n := negroni.Classic()
	n.UseHandler(mux)
	return n
}

// sensing data per device
var state = map[string]interface{}{}

// status from sensor
func PutStatusChangedHandle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	_status := map[string]interface{}{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&_status)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	did, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	// state[did] = _status
	// log.Println(state[did])

	_, ok = state[did]
	if !ok {
		state[did] = _status
	}

	b, err := json.Marshal(state[did])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(b)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)

	req, err := http.NewRequest("PUT", "http://"+constants.ServerAddr+":3000/device/"+did, bytes.NewReader(b))
	if err != nil {
		return
	}

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
}

func GetStatusHandle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	did, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_status, ok := state[did]
	if !ok {
		_status = map[string]interface{}{
			"msg": "hello world",
		}
	}

	encoder := json.NewEncoder(w)

	err := encoder.Encode(_status)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func PostStatusChangedHandle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	_status := map[string]interface{}{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&_status)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	did, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	state[did] = _status
	log.Println(_status)

	w.WriteHeader(http.StatusOK)
}

// status from user per device
