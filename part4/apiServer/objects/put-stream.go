package objects

import (
	"fmt"
	"object-storage/lib/objectStream"
	"object-storage/part4/apiServer/heartbeat"
)

func putStream(hash string, size int64) (*objectStream.TempPutStream, error) {
	server := heartbeat.ChooseRandomDataServer() //获取数据服务节点
	if server == "" {
		return nil, fmt.Errorf("cannot find any dataServer")
	}
	return objectStream.NewTempPutStream(server, hash, size) //上传文件
}
