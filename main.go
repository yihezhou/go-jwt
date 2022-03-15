package main

import (
	"log"
	"net/http"
)

func main() {
	// "Signin" and "Welcome" are the handlers that we will implement
	http.HandleFunc("/signin", logging(Signin))
	//http.HandleFunc("/welcome", Welcome)
	//http.HandleFunc("/refresh", Refresh)

	// start the http server on port 8000
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func logging(f http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		log.Println(request.URL.Path)
		f(writer, request)
	}
}