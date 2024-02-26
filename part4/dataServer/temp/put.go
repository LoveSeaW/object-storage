package temp

import (
	"log"
	"net/http"
	"object-storage/utils"
	"os"
	"path/filepath"
	"strings"
)

// 将临时文件转正
func put(w http.ResponseWriter, r *http.Request) {
	uuid := strings.Split(r.URL.EscapedPath(), "/")[2]
	temp_info, err := readFromFile(uuid)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	infoFile := filepath.Join(utils.STORAGE_ROOT + "/temp/" + uuid)
	dataFile := filepath.Join(infoFile + ".bat")
	file, err := os.Open(dataFile)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	actual := info.Size()
	os.Remove(infoFile)
	if actual != temp_info.Size {
		os.Remove(dataFile)
		log.Println("actual size mismatch, expect ", temp_info.Size, "actual", actual)
		return
	}
	commitTempObject(dataFile, temp_info)
}
