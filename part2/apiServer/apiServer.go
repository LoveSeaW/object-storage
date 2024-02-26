package main

import (
	"log"
	"net/http"
	"object-storage/part2/apiServer/heartbeat"
	"object-storage/part2/apiServer/locate"
	"object-storage/part2/apiServer/objects"
)

func main() {
	go heartbeat.ListenHeartBeat()
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/locate/", locate.Handler)
	log.Fatal(http.ListenAndServe(":9090", nil))
}
