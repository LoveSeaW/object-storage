package locate

import (
	"fmt"
	"object-storage/lib/rabbitmq"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

// 存储根目录
const STORAGE_ROOT = "E:\\Go\\object-storage\\file\\test"

var objects = make(map[string]int)
var mutex sync.Mutex

func Locate(name string) bool {
	mutex.Lock()
	_, err := os.Stat(name) //访问磁盘对应的文件名
	mutex.Unlock()
	return !os.IsNotExist(err)
}

func Add(hash string) {
	mutex.Lock()
	objects[hash] = 1
	mutex.Unlock()
}

func Del(hash string) {
	mutex.Lock()
	delete(objects, hash)
	mutex.Unlock()
}

func StartLocate() {
	queue := rabbitmq.New(rabbitmq.RABBITMQ_SERVER)
	defer queue.Close()
	queue.Bind("dataServers")
	consume := queue.Consume()
	//接口服务需要定位的名字
	for message := range consume {
		hash, err := strconv.Unquote(string(message.Body)) //处理消息中的转义字符""
		if err != nil {
			fmt.Println(err)
			continue
		}
		exist := Locate(hash)
		if exist {
			queue.Send(message.ReplyTo, rabbitmq.RABBITMQ_SERVER)
		}
		//STORAGE_ROOT = "E:\\Go\\object-storage\\file\\test"
		//fmt.Println(objects. + hash)
		//filepath := filepath.Join(objects.STORAGE_ROOT, object)
		//if Locate(filepath) { //storage_root
		//	queue.Send(message.ReplyTo, os.Getenv("LISTEN_ADDRESS"))
		//}
	}
}

// 获取所有文件散列值，写入缓存
func CollectObject() {
	files, _ := filepath.Glob(STORAGE_ROOT + "/objects/*")
	for i := range files {
		hash := filepath.Base(files[i])
		objects[hash] = 1
	}
}
