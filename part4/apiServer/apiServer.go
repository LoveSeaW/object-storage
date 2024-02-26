package apiServer

import (
	"log"
	"net/http"
	"object-storage/part4/apiServer/heartbeat"
	"object-storage/part4/apiServer/locate"
	"object-storage/part4/apiServer/version"
	"object-storage/part5/apiServer/objects"
)

func main() {
	go heartbeat.ListenHeartBeat()
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/locate/", locate.Handler)
	http.HandleFunc("/versions/", version.Handler)
	log.Fatal(http.ListenAndServe(":9090", nil))
}
