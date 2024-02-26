package objects

import (
	"log"
	"net/url"
	"object-storage/part4/dataServer/locate"
	"object-storage/utils"
	"os"
	"path/filepath"
)

func getFile(hash string) string {
	filePath := filepath.Join(locate.STORAGE_ROOT + "/objects/" + hash)
	file, _ := os.Open(filePath)
	document := url.PathEscape(utils.CalculateHash(file))
	file.Close()
	if document != hash {
		log.Println("object hash mismatch, remove", file)
		locate.Del(hash)
		os.Remove(filePath)
	}
	return filePath
}
