package rs

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"object-storage/lib/objectStream"
	"object-storage/utils"
)

type resumableToken struct {
	Name    string
	Size    int64
	Hash    string
	Servers []string
	Uuids   []string
}

type RSResumablePutStream struct {
	*RSPutStream
	*resumableToken
}

func NewRSResumablePutStream(dataServers []string, name, hash string, size int64) (*RSResumablePutStream, error) {
	putStream, err := NewRSPutStream(dataServers, hash, size)
	if err != nil {
		return nil, err
	}
	uuids := make([]string, AllSharad)
	for i := range uuids {
		uuids[i] = putStream.writers[i].(*objectStream.TempPutStream).Uuid
	}
	token := &resumableToken{
		Name:    name,
		Size:    size,
		Hash:    hash,
		Servers: dataServers,
		Uuids:   uuids,
	}
	return &RSResumablePutStream{
		RSPutStream:    putStream,
		resumableToken: token,
	}, nil

}

func PutStreamFromToken(inputToken string) (*RSResumablePutStream, error) {
	bytes, err := base64.StdEncoding.DecodeString(inputToken)
	if err != nil {
		return nil, err
	}
	var token resumableToken
	err = json.Unmarshal(bytes, &token)
	if err != nil {
		return nil, err
	}
	writers := make([]io.Writer, AllSharad)
	for i := range writers {
		writers[i] = &objectStream.TempPutStream{token.Servers[i], token.Uuids[i]}
	}
	encoder := NewEncoder(writers)
	return &RSResumablePutStream{&RSPutStream{encoder}, &token}, nil
}

func (s *RSResumablePutStream) ToToken() string {
	bytes, _ := json.Marshal(s)
	return base64.StdEncoding.EncodeToString(bytes)
}

func (s *RSResumablePutStream) CurrentSize() int64 {
	result, err := http.Head(fmt.Sprintf("http://%s/temp/%s", s.Servers[0], s.Uuids[0]))
	if err != nil {
		log.Println(err)
		return -1
	}
	if result.StatusCode != http.StatusOK {
		log.Println(result.StatusCode)
		return -1
	}
	size := utils.GetSizeFromHeader(result.Header) * DataShard
	if size > s.Size {
		size = s.Size
	}
	return size
}
