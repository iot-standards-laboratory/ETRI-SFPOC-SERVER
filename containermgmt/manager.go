package containermgnt

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

type Container struct {
	Image string `json:"image"`
	Name  string `json:"name"`
}

type ContainerReturnKey struct{}

var ReturnKey = ContainerReturnKey{}

func CreateContainer(ctx context.Context, cont *Container) context.Context {

	fmt.Println("container : ", *cont)
	args := strings.Split(fmt.Sprintf("container\\run\\-d\\--name\\%s\\%s", cont.Name, cont.Image), "\\")
	fmt.Println(args)
	bout, err := exec.Command("docker", args...).Output()
	if err != nil {
		log.Fatalln(err)
	}

	ctx = context.WithValue(ctx, ReturnKey, bout)

	return ctx
}

func GetContainers(ctx context.Context) context.Context {
	// cmd := exec.Command("firefox")
	// err := cmd.Run()

	cmd := strings.Split("container\\ls\\--format\\'{{.Image}} {{.Names}}'\\-a", "\\")
	bout, err := exec.Command("docker", cmd...).Output()
	if err != nil {
		log.Fatalln(err)
	}

	sout := strings.Split(string(bout), "\n")

	var list []Container
	for _, e := range sout {
		l := strings.Split(e, " ")

		if len(l) < 2 {
			continue
		}

		container := Container{
			Image: l[0],
			Name:  l[1],
		}

		list = append(list, container)
	}

	b, err := json.Marshal(list)
	if err != nil {
		log.Fatalln("JSON Marshal error!!")
	}

	fmt.Println(string(b))

	ctx = context.WithValue(ctx, ReturnKey, b)

	return ctx
}
