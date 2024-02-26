package heartbeat

import "math/rand"

// 随机返回一个数据服务节点
func ChooseRandomDataServer(n int, exclude map[int]string) (dataServers []string) {
	candidates := make([]string, 0)
	reverseExcludeMap := make(map[string]int)
	for id, address := range exclude {
		reverseExcludeMap[address] = id
	}
	servers := GetDataServes() //返回所有数据服务节点
	for i := range servers {
		server := servers[i]
		_, exclude := reverseExcludeMap[server]
		if !exclude {
			candidates = append(candidates, server)
		}
	}
	length := len(candidates)
	if length < n {
		return
	}
	part := rand.Perm(length)
	for i := 0; i < n; i++ {
		dataServers = append(dataServers, candidates[part[i]])
	}
	return
}
