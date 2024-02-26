package locate

import (
	"encoding/json"
	"object-storage/lib/rabbitmq"
	"object-storage/lib/rs"
	"object-storage/lib/types"
	"os"
	"time"
)

func Locate(name string) (locateInfo map[int]string) {
	queue := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	queue.Publish("dataServers", name)
	consume := queue.Consume()
	go func() {
		time.Sleep(time.Second)
		queue.Close()
	}()
	locateInfo = make(map[int]string)
	for i := 0; i < rs.AllSharad; i++ {
		message := <-consume
		if len(message.Body) == 0 {
			return locateInfo
		}
		var info types.LocateMessage
		json.Unmarshal(message.Body, &info)
		locateInfo[info.Id] = info.Address
	}
	return
}

func Exist(name string) bool {
	return len(Locate(name)) >= rs.DataShard
}
