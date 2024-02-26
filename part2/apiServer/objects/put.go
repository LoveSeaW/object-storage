package objects

import (
	"log"
	"net/http"
)

// 存储根目录
const STORAGE_ROOT = "E:\\Go\\object-storage\\file\\objects"

func put(w http.ResponseWriter, r *http.Request) {
	file, fileHeader, err := r.FormFile("filename")
	//if err != nil {
	//	http.Error(w, "Failed to retrieve file", http.StatusBadRequest)
	//	return
	//}
	defer file.Close()

	// 获取文件名
	filename := fileHeader.Filename

	c, err := storeObject(file, filename)
	if err != nil {
		log.Println(err)
	}
	w.WriteHeader(c)

	//// 构造文件路径
	//filepath := filepath2.Join(STORAGE_ROOT, filename)
	//
	////创建文件
	//newFile, err := os.Create(filepath)
	//if err != nil {
	//	http.Error(w, "Failed to create file", http.StatusInternalServerError)
	//	log.Println(err)
	//	return
	//}
	//defer newFile.Close()
	//
	////将文件内容拷贝到新文件中
	//_, err = io.Copy(newFile, file)
	//if err != nil {
	//	http.Error(w, "Failed to write file", http.StatusInternalServerError)
	//	log.Println(err)
	//	return
	//}
	//
	////返回成功响应
	//w.WriteHeader(http.StatusCreated)
}
