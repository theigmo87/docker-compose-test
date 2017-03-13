package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type selfInfo struct {
	ContainerID string  `json:"containerId"`
	Collection  string  `json:"collection"`
	Message     float64 `json:"message"`
}

func getSelfInfoHandler(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	hostname, _ := os.Hostname()
	retVal := selfInfo{hostname, params["collection"], rand.Float64()}
	json.NewEncoder(w).Encode(retVal)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/{collection}", getSelfInfoHandler).Methods("GET")
	log.Fatal(http.ListenAndServe(":80", router))
}
