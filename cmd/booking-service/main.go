package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Booking service")

	mux := http.NewServeMux()

	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "server is working")
	})

	http.ListenAndServe(":8080", mux)
}
