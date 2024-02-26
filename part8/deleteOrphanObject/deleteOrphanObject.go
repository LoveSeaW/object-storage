package main

import (
	"log"
	"net/http"
	"object-storage/lib/es"
	"os"
	"path/filepath"

	"strings"
)

func main() {
	files, _ := filepath.Glob(os.Getenv("STORAGE_ROOT") + "/objects/*")
	for i := range files {
		hash := strings.Split(filepath.Base(files[i]), ".")[0]
		hashInMetaData, err := es.HasHash(hash)
		if err != nil {
			log.Println(err)
			return
		}
		if !hashInMetaData {
			del(hash)
		}
	}
}

func del(hash string) {
	log.Println("delete:", hash)
	url := "http://" + os.Getenv("LISTEN_ADDRESS") + "/objects/" + hash
	request, _ := http.NewRequest("DELETE", url, nil)
	client := http.Client{}
	client.Do(request)
}
