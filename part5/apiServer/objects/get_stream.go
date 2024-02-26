package objects

import (
	"fmt"
	"object-storage/lib/rs"
	"object-storage/part5/apiServer/heartbeat"
	"object-storage/part5/apiServer/locate"
)

func GetStream(hash string, size int64) (*rs.RSGetStream, error) {
	locateInfo := locate.Locate(hash)
	if len(locateInfo) < rs.DataShard {
		return nil, fmt.Errorf("object %s locate fail, result %v", hash, locateInfo)
	}
	dataServers := make([]string, 0)
	if len(locateInfo) != rs.AllSharad {
		dataServers = heartbeat.ChooseRandomDataServer(rs.AllSharad-len(locateInfo), locateInfo)
	}
	return rs.NewRSGetStream(locateInfo, dataServers, hash, size)
}
