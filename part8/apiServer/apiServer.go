package main

import (
	"log"
	"net/http"
	"object-storage/part8/apiServer/heartbeat"
	"object-storage/part8/apiServer/locate"
	"object-storage/part8/apiServer/objects"
	"object-storage/part8/apiServer/temp"
	version "object-storage/part8/apiServer/versions"
	"os"
)

func main() {
	go heartbeat.ListenHeartBeat()
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/locate/", locate.Handler)
	http.HandleFunc("/versions/", version.Handler)
	http.HandleFunc("/temp/", temp.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
