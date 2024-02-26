package temp

import (
	"net/http"
	"object-storage/utils"
	"os"
	"path/filepath"
	"strings"
)

// 删除临时文件和临时文件信息
func del(w http.ResponseWriter, r *http.Request) {
	uuid := strings.Split(r.URL.EscapedPath(), "/")[2]
	infoFile := filepath.Join(utils.STORAGE_ROOT + "/temp/" + uuid)
	dataFile := filepath.Join(infoFile + ".bat")
	os.Remove(infoFile)
	os.Remove(dataFile)
}
