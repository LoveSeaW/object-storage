package temp

import (
	"object-storage/part4/dataServer/locate"
	"object-storage/utils"
	"os"
	"path/filepath"
)

// 文件转正
func commitTempObject(dataFile string, temp_info *tempInfo) {
	filePath := filepath.Join(utils.STORAGE_ROOT + "/objects" + temp_info.Name)
	os.Rename(dataFile, filePath)
	locate.Add(temp_info.Name)
}
