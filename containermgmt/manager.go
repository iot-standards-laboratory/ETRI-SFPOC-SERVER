package containermgmt

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// type Container struct {
// 	ID   string `json:"image"`
// 	Name string `json:"name"`
// 	Addr string `json:addr`
// }

func IsExist(name string) bool {
	cmd := strings.Split("container\\ls\\--format\\'{{.Image}} {{.Names}}'\\-a", "\\")
	bout, err := exec.Command("docker", cmd...).Output()
	if err != nil {
		log.Fatalln(err)
	}

	sout := strings.Split(string(bout), "\n")

	for _, e := range sout {
		l := strings.Split(e, " ")

		if len(l) < 2 {
			continue
		}

		if name == l[0] {
			return true
		}
	}

	return false
}

func CreateContainer(name string) error {

	if IsExist(name) {
		return nil
	}
	args := strings.Split(fmt.Sprintf("container\\run\\-d\\%s", name), "\\")
	fmt.Println(args)
	_, err := exec.Command("docker", args...).Output()
	if err != nil {
		return err
	}

	return nil
}
