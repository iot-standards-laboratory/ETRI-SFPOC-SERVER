package router

import (
	"bytes"
	"devicemanagerb/constants"
	"encoding/json"
	"io/ioutil"
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
var s_data = map[string]interface{}{}

// status from sensor
func PutStatusChangedHandle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	status := map[string]interface{}{}
	err = json.Unmarshal(b, &status)
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

	s_data[did] = status
	log.Println(s_data)

	cdata, ok := c_data[did]
	if !ok {
		w.Write([]byte("I am devicemanagerB"))
	} else {
		encoder := json.NewEncoder(w)
		err := encoder.Encode(cdata)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
	}

	req, err := http.NewRequest("PUT", "http://"+constants.ServerAddr+":3000/device/"+did, bytes.NewReader(b))
	if err != nil {
		return
	}

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

}

var c_data = map[string]interface{}{}

func GetStatusHandle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	did, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	status, ok := c_data[did]
	if !ok {
		status = map[string]interface{}{
			"servo": 0,
			"fan":   0,
			"light": 0,
		}
	}

	encoder := json.NewEncoder(w)

	err := encoder.Encode(status)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func PostStatusChangedHandle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	status := map[string]interface{}{}
	err = json.Unmarshal(b, &status)
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

	c_data[did] = status
	log.Println(c_data)

	w.WriteHeader(http.StatusOK)
}

// status from user per device
