package locate

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	if method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	//从url中获取文件名信息定位文件所在位置
	filename := r.URL.Query().Get("filename")
	fmt.Println(filename)
	//info := Locate(strings.Split(r.URL.EscapedPath(), "/")[2])
	//先用文件名代替，后根据文件元信息修改
	info := Locate(filename)
	fmt.Println("------------")
	fmt.Println(info)
	if len(info) == 0 {
		w.WriteHeader(http.StatusNotFound)
	}
	by, _ := json.Marshal(info)
	w.Write(by)
}
