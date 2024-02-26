package objects

import (
	"fmt"
	"io"
	"object-storage/lib/objectStream"
	"object-storage/part2/apiServer/locate"
)

func getStream(object string) (io.Reader, error) {
	server := locate.Locate(object) //获取文件存储的数据节点
	if server == "" {
		return nil, fmt.Errorf("object %s locate fail", object)
	}
	return objectStream.NewGetStream(server, object)
}
