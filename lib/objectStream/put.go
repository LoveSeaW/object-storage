package objectStream

import (
	"fmt"
	"io"
	"net/http"
)

type PutStream struct {
	writer *io.PipeWriter
	c      chan error
}

func NewPutStream(server, hash string, size int64) *PutStream {
	reader, writer := io.Pipe() //管道互联，写入writer的内容可以从reader中读出来
	c := make(chan error)
	go func() {
		request, _ := http.NewRequest("PUT", "http://"+server+
			"/objects/"+hash, reader)
		client := http.Client{}
		r, e := client.Do(request)
		if e == nil && r.StatusCode != http.StatusOK {
			e = fmt.Errorf("putStream dataServer return http code %d", r.StatusCode)
		}
		c <- e
	}()
	return &PutStream{writer, c}
}

// 实现io.write接口
func (w *PutStream) Write(p []byte) (n int, err error) {
	return w.writer.Write(p)
}

func (w *PutStream) Close() error {
	w.writer.Close()
	return <-w.c
}
