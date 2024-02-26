package objectStream

import (
	"fmt"
	"io"
	"net/http"
)

// 将http函数调用转换成读写流的形式
type GetStream struct {
	reader io.Reader
}

func newGetStream(url string) (*GetStream, error) {
	result, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if result.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GetStream dataServer return http code %d", result.StatusCode)
	}
	return &GetStream{result.Body}, nil
}

func NewGetStream(server, object string) (*GetStream, error) {
	if server == "" || object == "" {
		return nil, fmt.Errorf("invalid server %s object %s", server, object)
	}
	fmt.Println("http://" + server + "/objects/" + object)
	return newGetStream("http://" + server + "/objects/" + object)
}

func (r *GetStream) Read(p []byte) (n int, err error) {
	return r.reader.Read(p)
}
