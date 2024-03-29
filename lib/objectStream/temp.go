package objectStream

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

// TempPutStream 缓存传输,数据校验
type TempPutStream struct {
	Server string
	Uuid   string
}

func NewTempPutStream(server, object string, size int64) (*TempPutStream, error) {
	request, err := http.NewRequest("POST", "http://"+server+"/temp/"+object, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("size", fmt.Sprintf("%d", size))
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	uuid, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return &TempPutStream{server, string(uuid)}, nil
}

// 写入数据上传
func (w *TempPutStream) Write(p []byte) (n int, err error) {
	request, err := http.NewRequest("PATCH", "http://"+w.Server+"/temp/"+w.Uuid, strings.NewReader(string(p)))
	if err != nil {
		return 0, err
	}
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return 0, err
	}
	if response.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("dataServer return http code %d", response.StatusCode)
	}
	return len(p), nil
}

func (w *TempPutStream) Commit(good bool) {
	method := http.MethodDelete
	if good {
		method = http.MethodPut
	}
	requset, _ := http.NewRequest(method, "http://"+w.Server+"/temp/"+w.Uuid, nil)
	client := http.Client{}
	client.Do(requset)
}

func NewTempGetStream(server, uuid string) (*GetStream, error) {
	return newGetStream("http://" + server + "/temp/" + uuid)
}
