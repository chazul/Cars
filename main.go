package main

import (
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	carH := newCarHandler()

	mux.Handle("/cars", carH)  //root level
	mux.Handle("/cars/", carH) //sub levels

	http.ListenAndServe("localhost:8080", mux)
}
