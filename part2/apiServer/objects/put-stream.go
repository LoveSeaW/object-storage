package objects

import (
	"fmt"
	"object-storage/lib/objectStream"
	"object-storage/part2/apiServer/heartbeat"
)

func putStream(object string) (*objectStream.PutStream, error) {
	server := heartbeat.ChooseRandomDataServer() //获取数据服务节点
	fmt.Println("putStream" + server)
	if server == "" {
		return nil, fmt.Errorf("cannot find any dataServer")
	}
	return objectStream.NewPutStream(server, object), nil //上传文件
}
