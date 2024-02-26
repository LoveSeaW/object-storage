package rs

import (
	"io"
	"object-storage/lib/objectStream"
)

type RSResumableGetStream struct {
	*decoder
}

func NewRSResumableGetStream(dataServers []string, uuids []string, size int64) (*RSResumableGetStream, error) {
	readers := make([]io.Reader, AllSharad)
	for i := 0; i < AllSharad; i++ {
		readers[i], err = objectStream.NewTempGetStream(dataServers[i], uuids[i])
		if err != nil {
			return nil, err
		}
	}
	writers := make([]io.Writer, AllSharad)
	decoder := NewDecoder(readers, writers, size)
	return &RSResumableGetStream{decoder}, nil
}
