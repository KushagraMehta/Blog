package main

import (
	"fmt"
	"log"
	"net/http"
)

func homePage(w http.ResponseWriter, r *http.Request) {

	log.Println("Endpoint Hit: homePage")

	fmt.Fprintf(w, "Welcome to the HomePage!")
}
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", homePage)
	log.Fatal(http.ListenAndServe(":8090", mux))
}
