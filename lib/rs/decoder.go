package rs

import (
	"github.com/klauspost/reedsolomon"
	"io"
)

// 解码
type decoder struct {
	readers   []io.Reader
	writer    []io.Writer
	encode    reedsolomon.Encoder
	size      int64
	cache     []byte
	cacheSize int
	total     int64
}

func NewDecoder(readers []io.Reader, writers []io.Writer, size int64) *decoder {
	encode, _ := reedsolomon.New(DataShard, ParityShard)
	return &decoder{
		readers:   readers,
		writer:    writers,
		encode:    encode,
		size:      size,
		cache:     nil,
		cacheSize: 0,
		total:     0,
	}
}

func (d *decoder) Read(p []byte) (n int, err error) {
	if d.cacheSize == 0 {
		err := d.getData()
		if err != nil {
			return 0, err
		}
	}
	length := len(p) //计算所需数据的长度
	if d.cacheSize < length {
		length = d.cacheSize
	}
	d.cacheSize -= length
	copy(p, d.cache[:length])
	d.cache = d.cache[length:]
	return length, nil
}

func (d *decoder) getData() error {
	if d.total == d.size {
		return io.EOF
	}
	shards := make([][]byte, AllSharad) //数据分片存储
	repairIds := make([]int, 0)
	for i := range shards {
		if d.readers[i] == nil { //数据丢失
			repairIds = append(repairIds, i)
		} else {
			shards[i] = make([]byte, BlockPerShard)
			count, err := io.ReadFull(d.readers[i], shards[i]) //从reader中读取数据
			if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
				shards[i] = nil
			} else if count != BlockPerShard {
				shards[i] = shards[i][:count]
			}
		}
	}
	err := d.encode.Reconstruct(shards) //恢复数据
	if err != nil {
		return err
	}
	for i := range repairIds {
		id := repairIds[i]
		d.writer[id].Write(shards[id])
	}
	for i := 0; i < DataShard; i++ { //读取四个分片恢复文件
		shardSie := int64(len(shards[i]))
		if d.total+shardSie > d.size {
			shardSie -= d.total + shardSie - d.size
		}
		d.cache = append(d.cache, shards[i][:shardSie]...)
		d.cacheSize += int(shardSie)
		d.total += shardSie
	}
	return nil
}
