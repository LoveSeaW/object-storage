package apiServer

import (
	"log"
	"net/http"
	"object-storage/part7/apiServer/heartbeat"
	"object-storage/part7/apiServer/locate"
	"object-storage/part7/apiServer/temp"
	"object-storage/part7/apiServer/version"

	"object-storage/part7/apiServer/objects"
	"os"
)

func main() {
	go heartbeat.ListenHeartBeat()
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/locate/", locate.Handler)
	http.HandleFunc("/version/", version.Handler)
	http.HandleFunc("/temp/", temp.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
