package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", HelloServer)
	http.ListenAndServe(":7777", nil)
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("\n\n\n Got request. Method: %s \n URL: %s \n Header: %s \n Body: %s", r.Method, r.URL, r.Header, r.Body)
	result := MakeGreeting(r.URL.Path[1:])
	fmt.Fprintf(w, result)
}

func MakeGreeting(name string) string {
	return fmt.Sprintf("Hello, %s!\n\n\n\n", name)
}
