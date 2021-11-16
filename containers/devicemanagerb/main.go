package main

import (
	"bytes"
	"devicemanagerb/router"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
)

var server_addr string

func getIP() string {
	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)

	return addrs[0].String()
}
func registerToServer() {
	var exist bool
	server_addr, exist = os.LookupEnv("SERVER_ADDR")
	if !exist {
		fmt.Println("Please set SERVER_ADDR as environment variable")
	}

	ip := getIP()
	idx := strings.LastIndex(ip, ".")
	serverAddr := ip[:idx+1] + "1"
	fmt.Println(serverAddr)

	var obj map[string]string = make(map[string]string)
	obj["name"] = "devicemanagerb"
	obj["addr"] = getIP() + ":3000"

	b, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("PUT", "http://"+serverAddr+":3000/services", bytes.NewBuffer(b))
	if err != nil {
		panic(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// http.Post("")
}

func main() {
	registerToServer()
	err := http.ListenAndServe(":3000", router.NewRouter())
	if err != nil {
		panic(err)
	}
}
