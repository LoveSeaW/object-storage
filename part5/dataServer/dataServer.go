package dataServer

import (
	"log"
	"net/http"
	"object-storage/part5/dataServer/heartbeat"
	"object-storage/part5/dataServer/locate"
	"object-storage/part5/dataServer/objects"
	"object-storage/part5/dataServer/temp"
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
