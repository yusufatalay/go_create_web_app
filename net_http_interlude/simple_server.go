package main

import (
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:]) // prints the path entered after root
}

func main() {
	http.HandleFunc("/", handler)                // handle all requests to the root with handler function
	log.Fatal(http.ListenAndServe(":8080", nil)) // listen the port 8080 listenandserve always returns an error
}
