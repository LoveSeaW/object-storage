package main

import (
	"log"
	"net/http"
	"object-storage/part3/apiServer/heartbeat"
	"object-storage/part3/apiServer/locate"
	"object-storage/part3/apiServer/objects"
	"object-storage/part3/apiServer/version"
)

func main() {
	go heartbeat.ListenHeartBeat()
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/locate/", locate.Handler)
	http.HandleFunc("/versions/", version.Handler)
	log.Fatal(http.ListenAndServe(":9090", nil))
}
