package heartbeat

import "math/rand"

// 随机返回一个数据服务节点
func ChooseRandomDataServer() string {
	dataServers := GetDataServes()
	count := len(dataServers)
	if count == 0 {
		return ""
	}
	return dataServers[rand.Intn(count)]
}
