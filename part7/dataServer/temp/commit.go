package temp

import (
	"compress/gzip"
	"net/url"
	"object-storage/part7/dataServer/locate"
	"object-storage/utils"
	"os"
	"strconv"
	"strings"
)

func (t *tempInfo) hash() string {
	str := strings.Split(t.Name, ".")
	return str[0]
}

func (t *tempInfo) id() int {
	str := strings.Split(t.Name, ".")
	id, _ := strconv.Atoi(str[1])
	return id
}

func commitTempObject(dataFile string, info *tempInfo) {
	file, _ := os.Open((dataFile))
	defer file.Close()
	calculated := url.PathEscape(utils.CalculateHash(file))
	file.Close()
	path, _ := os.Create(os.Getenv("STORAGE_ROOT") + "/object/" + info.Name + "." + calculated)
	writer := gzip.NewWriter(path)
	writer.Close()
	os.Remove(dataFile)
	locate.Add(info.hash(), info.id())
}
