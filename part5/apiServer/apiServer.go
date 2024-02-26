package apiServer

import (
	"log"
	"net/http"
	"object-storage/part5/apiServer/heartbeat"
	"object-storage/part5/apiServer/locate"
	"object-storage/part5/apiServer/version"

	"object-storage/part5/apiServer/objects"
	"os"
)

func main() {
	go heartbeat.ListenHeartBeat()
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/locate/", locate.Handler)
	http.HandleFunc("/version/", version.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
