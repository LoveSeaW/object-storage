package objects

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func get(w http.ResponseWriter, r *http.Request) {
	//获取url参数中的文件名
	//http://127.0.0.1:8081/objects/7.jpg
	objectName := strings.Split(r.URL.EscapedPath(), "/")[2]
	fmt.Printf("object is %s", objectName)
	if objectName == "" {
		http.Error(w, "not found filename", http.StatusBadRequest)
		return
	}
	//构造文件路径
	filepath := filepath.Join(STORAGE_ROOT, objectName)
	fmt.Println(filepath)
	//打开文件
	file, err := os.Open(filepath)
	if err != nil {
		http.Error(w, "File not Found", http.StatusNotFound)
		log.Println(err)
		return
	}
	defer file.Close()

	//将文件内容拷贝到响应中
	_, err = io.Copy(w, file)
	if err != nil {
		log.Println(err)
	}

}
