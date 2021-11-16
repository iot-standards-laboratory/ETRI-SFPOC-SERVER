package main

import (
	"etrismartfarmpoc/router"
	"net/http"
)

func main() {
	err := http.ListenAndServe(":3000", router.NewRouter())
	if err != nil {
		panic(err)
	}

	// fmt.Scanln()

	// fmt.Println(uuid.New())
	// bootstrap.RunBootstrapServer()

}
