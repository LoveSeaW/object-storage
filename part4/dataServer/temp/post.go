package temp

import (
	"encoding/json"
	"log"
	"net/http"
	"object-storage/utils"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type tempInfo struct {
	Uuid string
	Name string
	Size int64
}

// 在数据服务上创造临时对象
func post(w http.ResponseWriter, r *http.Request) {
	output, _ := exec.Command("uuidgen").Output()
	uuid := strings.TrimSuffix(string(output), "\n")
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	size, err := strconv.ParseInt(r.Header.Get("size"), 0, 64)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	temp := tempInfo{uuid, name, size}
	err = temp.writeToFile()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	filePath := filepath.Join(utils.STORAGE_ROOT + "/temp/" + temp.Uuid + ".bat")
	os.Create(filePath)
	w.Write([]byte(uuid))

}

// 创建临时文件信息
func (t *tempInfo) writeToFile() error {
	filePath := filepath.Join(utils.STORAGE_ROOT + "/temp/" + t.Uuid)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	bytes, _ := json.Marshal(t)
	file.Write(bytes)
	return nil
}
