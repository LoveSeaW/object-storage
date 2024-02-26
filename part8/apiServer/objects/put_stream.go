package objects

import (
	"fmt"
	"object-storage/lib/rs"
	"object-storage/part8/apiServer/heartbeat"
)

func putStream(hash string, size int64) (*rs.RSPutStream, error) {
	servers := heartbeat.ChooseRandomDataServers(rs.AllSharad, nil)
	if len(servers) != rs.AllSharad {
		return nil, fmt.Errorf("cannot find enough dataServer")
	}
	return rs.NewRSPutStream(servers, hash, size)
}
