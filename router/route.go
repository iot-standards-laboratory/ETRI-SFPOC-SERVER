package router

import (
	"bytes"
	"context"
	"encoding/json"
	containermgnt "etrismartfarmpoc/containermgmt"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// var rd *render.Render = render.New()

func NewRouter() http.Handler {
	mux := mux.NewRouter()

	mux.HandleFunc("/service", GetServiceList).Methods("GET")
	mux.HandleFunc("/extension", PostService).Methods("POST")

	// mux.Handle("/",  http.FileServer(http.Dir("public")))
	n := negroni.Classic() // 파일 서버 및 로그기능을 제공함
	n.UseHandler(mux)

	return n
}

func GetServiceList(w http.ResponseWriter, r *http.Request) {
	result := containermgnt.GetContainers(context.Background())
	b := result.Value(containermgnt.ReturnKey).([]byte)
	w.Write(b)
}

func PostService(w http.ResponseWriter, r *http.Request) {
	request := containermgnt.Container{
		Image: "hello-world",
		Name:  "TempContainer",
	}

	s, err := json.Marshal(request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}

	fmt.Println(s)
	cont := new(containermgnt.Container)
	decoder := json.NewDecoder(bytes.NewReader(s))
	decoder.Decode(cont)

	result := containermgnt.CreateContainer(context.Background(), cont)
	b := result.Value(containermgnt.ReturnKey).([]byte)
	w.Write(b)
}
