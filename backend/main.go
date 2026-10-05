package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "FitMate API is running")
	})

	fmt.Println("FitMate backend started")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}