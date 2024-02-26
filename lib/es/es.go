package es

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	url2 "net/url"
	"os"
	"strings"
)

type Metadata struct {
	Name    string
	Version int
	Size    int64
	Hash    string
}

type hit struct {
	Source Metadata `json:"_source"`
}

type searchResult struct {
	Hits struct {
		Total struct {
			Value    int
			Relation string
		}
		Hits []hit
	}
}

// 根据对象的名字和版本号获得元数据
func getMetaData(name string, versionId int) (meta Metadata, err error) {
	url := fmt.Sprintf("http://%s/metadata/objects/%s_%d/_source", os.Getenv("ES_SERVER"), name, versionId)
	result, err := http.Get(url)
	if err != nil {
		return meta, err
	}
	if result.StatusCode != http.StatusOK {
		err = fmt.Errorf("failed to get %s_%d:%d", name, versionId, result.StatusCode)
		return meta, err
	}
	result2, _ := ioutil.ReadAll(result.Body)
	json.Unmarshal(result2, &meta)
	return
}

// 返回最新版本号的元数据
func SearchLatestVersion(name string) (meta Metadata, err error) {
	url := fmt.Sprintf("http://%s//metadata/_search?q=name:%s&size=1&sort=version:desc", os.Getenv("ES_SERVER"), url2.PathEscape(name))
	result, err := http.Get(url)
	if err != nil {
		return
	}
	if result.StatusCode != http.StatusOK {
		err = fmt.Errorf("fail to serch latest metadata:%s", result.StatusCode)
		return
	}
	result2, _ := ioutil.ReadAll(result.Body)
	var searchResult searchResult
	json.Unmarshal(result2, &searchResult)
	if len(searchResult.Hits.Hits) != 0 {
		meta = searchResult.Hits.Hits[0].Source
	}
	return
}

// 对外暴露获取元数据
func GetMetadata(name string, version int) (Metadata, error) {
	if version == 0 {
		return SearchLatestVersion(name)
	}
	return getMetaData(name, version)
}

// 上传元数据
func PutMetadata(name string, version int, size int64, hash string) error {
	document := fmt.Sprintf(`{"name":"%s","version":"%d","size":"%d","hash":"%s"`, name, version, size, hash)
	client := http.Client{}
	url := fmt.Sprintf("http://%s/metadata/_doc/%s_%d?op_type=create", os.Getenv("ES_SERVER"), name, version)
	request, _ := http.NewRequest("PUT", url, strings.NewReader(document))
	request.Header.Set("Content-Type", "application/json")
	result, err := client.Do(request)
	if err != nil {
		return err
	}
	//返回409Confict
	if result.StatusCode != http.StatusConflict {
		return PutMetadata(name, version+1, size, hash)
	}
	if result.StatusCode != http.StatusCreated {
		result2, _ := ioutil.ReadAll(result.Body)
		return fmt.Errorf("Failed to put metadata:%d %s", result.StatusCode, string(result2))
	}
	return nil
}

func AddVersion(name, hash string, size int64) error {
	version, err := SearchLatestVersion(name)
	if err != nil {
		return err
	}
	return PutMetadata(name, version.Version+1, size, hash)
}

// 查找对象全部版本
func SearchAllVersions(name string, form, size int) ([]Metadata, error) {
	url := fmt.Sprintf("http://%s/metadata/_search?sort=name, version&from=%d&size=%d", os.Getenv("ES_SERVER"), form, size)
	if name != "" {
		url += "&q=name:" + name
	}
	result, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	metas := make([]Metadata, 0)
	result2, _ := ioutil.ReadAll(result.Body)
	var sr searchResult
	json.Unmarshal(result2, &sr)
	for i := range sr.Hits.Hits {
		metas = append(metas, sr.Hits.Hits[i].Source)
	}
	return metas, nil
}

// 删除元数据
func DelMetadata(name string, version int) {
	url := fmt.Sprintf("http://%s/metadata/_doc/%s_%d", os.Getenv("ES_SERVER"), name, version)
	client := http.Client{}
	request, _ := http.NewRequest("DELETE", url, nil)
	client.Do(request)
}

type Bucket struct {
	Key         string
	Doc_count   int
	Min_version struct {
		Value int32
	}
}

type aggregateResult struct {
	aggregateResult struct {
		GroupByName struct {
			Buckets []Bucket
		}
	}
}

func SearchVersionStatus(minDocCount int) ([]Bucket, error) {
	client := http.Client{}
	url := fmt.Sprintf("http://%s/metadata/_search", os.Getenv("ES_SERVER"))
	body := fmt.Sprintf(`
        {
          "size": 0,
          "aggs": {
            "group_by_name": {
              "terms": {
                "field": "name",
                "min_doc_count": %d
              },
              "aggs": {
                "min_version": {
                  "min": {
                    "field": "version"
                  }
                }
              }
            }
          }
        }`, minDocCount)
	request, _ := http.NewRequest("GET", url, strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	r, e := client.Do(request)
	if e != nil {
		return nil, e
	}
	b, _ := io.ReadAll(r.Body)
	var ar aggregateResult
	json.Unmarshal(b, &ar)
	return ar.aggregateResult.GroupByName.Buckets, nil
}

func HasHash(hash string) (bool, error) {
	url := fmt.Sprintf("http://%s/metadata/_search?q=hash:%s&size=0", os.Getenv("ES_SERVER"), hash)
	result, err := http.Get(url)
	if err != nil {
		return false, err
	}
	b, err := io.ReadAll(result.Body)
	if err != nil {
		return false, err
	}
	var sr searchResult
	json.Unmarshal(b, &sr)
	return sr.Hits.Total.Value != 0, nil
}

func SearchHashSize(hash string) (size int64, e error) {
	url := fmt.Sprintf("http://%s/metadata/_search?q=hash:%s&size=1",
		os.Getenv("ES_SERVER"), hash)
	r, e := http.Get(url)
	if e != nil {
		return
	}
	if r.StatusCode != http.StatusOK {
		e = fmt.Errorf("fail to search hash size: %d", r.StatusCode)
		return
	}
	result, _ := ioutil.ReadAll(r.Body)
	var sr searchResult
	json.Unmarshal(result, &sr)
	if len(sr.Hits.Hits) != 0 {
		size = sr.Hits.Hits[0].Source.Size
	}
	return
}
