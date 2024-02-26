package rs

import (
	"fmt"
	"io"
	"object-storage/lib/objectStream"
)

type RSGetStream struct {
	*decoder
}

var err error

func NewRSGetStream(locateInfo map[int]string, dataServers []string, hash string, size int64) (*RSGetStream, error) {
	if len(locateInfo)+len(dataServers) != AllSharad {
		return nil, fmt.Errorf("dataServers number mismatch")
	}
	readers := make([]io.Reader, AllSharad)
	for i := 0; i < AllSharad; i++ {
		server := locateInfo[i] //查找分片所在数据节点
		if server == "" {
			locateInfo[i] = dataServers[0]
			dataServers = dataServers[1:]
			continue
		}
		reader, err := objectStream.NewGetStream(server, fmt.Sprintf("%s.%d", hash, i))
		if err == nil {
			readers[i] = reader //读取流
		}
	}
	writers := make([]io.Writer, AllSharad)
	perShard := (size + DataShard - 1) / DataShard

	for i := range readers {
		if readers[i] == nil {
			//读取流为空，在writers存在写入流
			writers[i], err = objectStream.NewTempPutStream(locateInfo[i], fmt.Sprintf("%s.%d", hash, i), perShard)
			if err != nil {
				return nil, err
			}
		}
	}
	decoder := NewDecoder(readers, writers, size)
	return &RSGetStream{decoder}, nil
}

func (s *RSGetStream) Close() {
	for i := range s.writer {
		if s.writer[i] != nil {
			s.writer[i].(*objectStream.TempPutStream).Commit(true) //将写入流内容转正
		}
	}
}

// 跳转客户请求位置
func (s *RSGetStream) Seek(offset int64, whence int) (int64, error) {
	if whence != io.SeekCurrent {
		panic("only support SeekCurrent")
	}
	if offset < 0 {
		panic("only support forward seek")
	}
	for offset != 0 {
		length := int64(BlockSie)
		if offset < length {
			length = offset
		}
		buffer := make([]byte, length)
		io.ReadFull(s, buffer)
		offset -= length
	}
	return offset, nil
}
