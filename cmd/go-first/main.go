package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	log.Print("Hello, World")

	router := mux.NewRouter()
	router.HandleFunc("/", hello)

	go func() {
		err := http.ListenAndServe(":8081", router)
		if err != nil {
			log.Fatal(err.Error())
		}
	}()

	diagnostics1 := diagnostics.NewRouter()
	err := http.ListenAndServe(":8585", diagnostics1)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, http.StatusText(http.StatusOK))
}
