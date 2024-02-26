package heartbeat

import (
	"fmt"
	"object-storage/lib/rabbitmq"
	"strconv"
	"sync"
	"time"
)

var dataServers = make(map[string]time.Time) //缓存所有数据服务节点
var mutex sync.Mutex                         //go语言中map不是并发安全的

func ListenHeartBeat() {
	queue := rabbitmq.New(rabbitmq.RABBITMQ_SERVER) //rabbitmq服务结点，可以通过命令行获得
	defer queue.Close()
	queue.Bind("apiServer")
	consume := queue.Consume() //channel类型，获取队列消息
	go removeExpiredDataServer()
	for message := range consume {
		dataServer, err := strconv.Unquote(string(message.Body))
		if dataServer == "" {
			fmt.Errorf("dataServer is nil")
		}
		if err != nil {
			fmt.Println(err)
		}
		mutex.Lock()
		dataServers[dataServer] = time.Now() //收到消息的时间

		mutex.Unlock()

	}
}

// 移除没有10s没有收到心跳消息的数据服务节点
func removeExpiredDataServer() {
	for {
		time.Sleep(5 * time.Second)
		mutex.Lock()
		for server, timer := range dataServers {
			if timer.Add(10 * time.Second).Before(time.Now()) {
				delete(dataServers, server)
			}
		}
		mutex.Unlock()
	}
}

// 获取所有数据节点
func GetDataServes() []string {
	mutex.Lock()
	defer mutex.Unlock()
	dataServer := make([]string, 0)
	for server := range dataServers {
		dataServer = append(dataServer, server)
		fmt.Println(server)
	}
	return dataServer
}
