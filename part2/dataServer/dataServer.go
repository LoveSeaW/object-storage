package main

import (
	"log"
	"net/http"
	"object-storage/part2/dataServer/objects"

	"object-storage/part2/dataServer/heartbeat"
	"object-storage/part2/dataServer/locate"
	"os"
)

func main() {
	go heartbeat.StartHeartBeat()
	go locate.StartLocate()
	os.Setenv("LISTEN_ADDRESS", "127.0.0.1:8081")
	http.HandleFunc("/objects/", objects.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
