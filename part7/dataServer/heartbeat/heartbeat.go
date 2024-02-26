package heartbeat

import (
	"object-storage/lib/rabbitmq"
	"os"
	"time"
)

func StartHeartBeat() {
	queue := rabbitmq.New(rabbitmq.RABBITMQ_SERVER)
	defer queue.Close()
	for {
		queue.Publish("apiServer", os.Getenv("LISTEN_ADDRESS")) //向交换机无限循环发送消息
		time.Sleep(5 * time.Second)
	}
}
