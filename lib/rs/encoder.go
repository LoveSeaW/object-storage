package rs

import (
	"github.com/klauspost/reedsolomon"
	"io"
)

// 编码
type encoder struct {
	writers []io.Writer
	encoder reedsolomon.Encoder
	cache   []byte
}

func NewEncoder(writers []io.Writer) *encoder {
	encode, _ := reedsolomon.New(DataShard, ParityShard)
	return &encoder{writers, encode, nil}
}

func (e *encoder) Write(p []byte) (count int, err error) {
	length := len(p)
	current := 0
	for length != 0 {
		next := BlockSie - len(e.cache)
		if next > length {
			next = length
		}
		e.cache = append(e.cache, p[current:current+next]...) //以块的形式放入缓存
		if len(e.cache) == BlockSie {
			e.Flush() //缓存已满，调用Flush写入writers
		}
		current += next //已写入的内容
		length -= next  //数据剩余内容
	}
	return len(p), nil
}

func (e *encoder) Flush() {
	if len(e.cache) == 0 {
		return
	}
	shards, _ := e.encoder.Split(e.cache) //将缓存数据切成四个数据片
	e.encoder.Encode(shards)              //生成两个校验片
	for i := range shards {
		e.writers[i].Write(shards[i]) //写入writer
	}
	e.cache = []byte{} //清空缓存
}
