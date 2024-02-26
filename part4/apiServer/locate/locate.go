package locate

import (
	"fmt"
	"object-storage/lib/rabbitmq"
	"strconv"
	"time"
)

// 定位文件
func Locate(name string) string {
	queue := rabbitmq.New(rabbitmq.RABBITMQ_SERVER)
	fmt.Printf("要发送的消息内容：%s\n", name)
	queue.Publish("dataServers", name)

	consume := queue.Consume()
	//关闭临时队列
	go func() {
		time.Sleep(time.Second)
		queue.Close()
	}()
	fmt.Println(consume)
	message := <-consume
	fmt.Printf("Received message: %s\n", string(message.Body))
	//获取文件所在的数据服务节点
	str, _ := strconv.Unquote(string(message.Body))
	fmt.Println(str)
	return str
}

func Exist(name string) bool {
	return Locate(name) != ""
}
