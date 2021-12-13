package constants

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func getIP() string {
	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)

	return addrs[0].String()
}

var ServerAddr string
var MyIP string

func init() {
	var exist bool
	ServerAddr, exist = os.LookupEnv("SERVER_ADDR")
	if !exist {
		fmt.Println("Please set SERVER_ADDR as environment variable")
	}

	MyIP = getIP()
	idx := strings.LastIndex(MyIP, ".")
	ServerAddr = MyIP[:idx+1] + "1"
	fmt.Println(ServerAddr)

}
