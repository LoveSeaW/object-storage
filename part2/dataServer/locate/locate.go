package locate

import (
	"fmt"
	"object-storage/lib/rabbitmq"
	"object-storage/part2/dataServer/objects"
	"os"
	filepath2 "path/filepath"
	"strconv"
)

func Locate(name string) bool {

	_, err := os.Stat(name) //访问磁盘对应的文件名
	return !os.IsNotExist(err)
}

func StartLocate() {
	queue := rabbitmq.New(rabbitmq.RABBITMQ_SERVER)
	defer queue.Close()
	queue.Bind("dataServers")
	consume := queue.Consume()
	//接口服务需要定位的名字
	for message := range consume {
		object, err := strconv.Unquote(string(message.Body)) //处理消息中的转义字符""
		if err != nil {
			panic(err)
		}
		//STORAGE_ROOT = "E:\\Go\\object-storage\\file\\test"
		fmt.Println(objects.STORAGE_ROOT + object)
		filepath := filepath2.Join(objects.STORAGE_ROOT, object)
		if Locate(filepath) { //storage_root
			queue.Send(message.ReplyTo, os.Getenv("LISTEN_ADDRESS"))
		}
	}
}
