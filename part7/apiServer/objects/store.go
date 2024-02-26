package objects

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"object-storage/part7/apiServer/locate"
	"object-storage/utils"
)

// 使用元数据作为标识符
func storeObject(r io.Reader, hash string, size int64) (int, error) {
	if locate.Exist(url.PathEscape(hash)) {
		return http.StatusOK, nil
	}
	stream, err := putStream(url.PathEscape(hash), size)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	reader := io.TeeReader(r, stream)       //读取r内容同时写入stream
	caculate := utils.CalculateHash(reader) //获取文件内容计算出散列值
	if caculate != hash {
		stream.Commit(false)
		return http.StatusBadRequest, fmt.Errorf("object hash mismatch, caculated=%s, requestd=%s", caculate, hash)
	}
	stream.Commit(true)
	return http.StatusOK, nil
}
