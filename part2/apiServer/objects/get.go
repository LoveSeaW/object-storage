package objects

import (
	"io"
	"log"
	"net/http"
)

func get(w http.ResponseWriter, r *http.Request) {
	//获取url参数中的文件名
	objectName := r.URL.Query().Get("filename")
	if objectName == "" {
		http.Error(w, "not found filename", http.StatusBadRequest)
		return
	}

	stream, err := getStream(objectName)
	if err != nil {
		return
	}
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	io.Copy(w, stream)
	////构造文件路径
	//filepath := filepath2.Join(STORAGE_ROOT, objectName)
	////打开文件
	//file, err := os.Open(filepath)
	//if err != nil {
	//	http.Error(w, "File not Found", http.StatusNotFound)
	//	log.Println(err)
	//	return
	//}
	//defer file.Close()
	//
	////将文件内容拷贝到响应中
	//_, err = io.Copy(w, file)
	//if err != nil {
	//	log.Println(err)
	//}

}
