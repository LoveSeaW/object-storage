package main

import (
	"log"
	"net/http"
	"object-storage/part8/dataServer/heartbeat"
	"object-storage/part8/dataServer/locate"
	"object-storage/part8/dataServer/objects"
	"object-storage/part8/dataServer/temp"
	"os"
)

func main() {
	locate.CollectObjects()
	go heartbeat.StartHeartBeat()
	go locate.StartLocate()
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/temp/", temp.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
