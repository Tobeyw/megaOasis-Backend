package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println(http.Dir("./"))

	http.Handle("/", http.FileServer(http.Dir("/home/ubuntu/go/megaOasis-Backend/img/")))
	e := http.ListenAndServe("0.0.0.0:8090", nil)
	if e != nil {
		fmt.Println(e)
	}
}
