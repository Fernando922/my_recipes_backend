package main

import (
	"fmt"
	"log"
	"net/http"
	"recipes-backend/server"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/types", server.CreateRecipeType).Methods(http.MethodPost)

	fmt.Println("Server listening on port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))

}
