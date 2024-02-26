package objects

import (
	"log"
	"net/http"
	"net/url"
	"object-storage/lib/es"
	"object-storage/lib/rs"
	"object-storage/part7/apiServer/heartbeat"
	"object-storage/part7/apiServer/locate"
	"object-storage/utils"
	"strconv"
	"strings"
)

func post(w http.ResponseWriter, r *http.Request) {
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	size, err := strconv.ParseInt(r.Header.Get("size"), 0, 64)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusForbidden)
		return
	}
	hash := utils.GetHashFromHeader(r.Header)
	if hash == "" {
		log.Println("missing object hash in digest header")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if locate.Exist(url.PathEscape(hash)) {
		err = es.AddVersion(name, hash, size)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
		return
	}
	dataServers := heartbeat.ChooseRandomDataServer(rs.AllSharad, nil)
	if len(dataServers) != rs.AllSharad {
		log.Println("cannot find enough dataServer")
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	stream, err := rs.NewRSResumablePutStream(dataServers, name, url.PathEscape(hash), size)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("location", "/temp/"+url.PathEscape(stream.ToToken()))
	w.WriteHeader(http.StatusCreated)
}
