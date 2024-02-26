package apiServer

import (
	"log"
	"net/http"
	"object-storage/part6/apiServer/heartbeat"
	"object-storage/part6/apiServer/locate"
	"object-storage/part6/apiServer/temp"
	"object-storage/part6/apiServer/version"

	"object-storage/part6/apiServer/objects"
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
