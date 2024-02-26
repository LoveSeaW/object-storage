package main

import (
	"log"
	"net/http"
	"object-storage/part4/dataServer/objects"
	"object-storage/part4/dataServer/temp"

	"object-storage/part4/dataServer/heartbeat"
	"object-storage/part4/dataServer/locate"
	"os"
)

func main() {
	go heartbeat.StartHeartBeat()
	go locate.StartLocate()
	os.Setenv("LISTEN_ADDRESS", "127.0.0.1:8081")
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/temp/", temp.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
