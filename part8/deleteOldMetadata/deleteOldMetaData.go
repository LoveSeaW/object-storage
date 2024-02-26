package main

import (
	"log"
	"object-storage/lib/es"
)

const MinVersionCount = 5

func main() {
	buckets, err := es.SearchVersionStatus(MinVersionCount + 1)
	if err != nil {
		log.Println(err)
		return
	}
	for i := range buckets {
		bucket := buckets[i]
		for v := 0; v < bucket.Doc_count-MinVersionCount; v++ {
			es.DelMetadata(bucket.Key, v+int(bucket.Min_version.Value))
		}
	}
}
