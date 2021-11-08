package main

import (
	"etrismartfarmpoc/router"
	"net/http"
)

func main() {
	http.ListenAndServe(":3000", router.NewRouter())
	// fmt.Println(string(out))
}
