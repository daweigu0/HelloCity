package qiniu

import (
	"HelloCity/internal/utils"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"testing"
)

func TestService_UploadFile(t *testing.T) {
	accessKey := utils.Config.GetString("oss.qiniu.accessKey")
	secretKey := utils.Config.GetString("oss.qiniu.secretKey")
	fileSvc := NewService(accessKey, secretKey)
	filePath := `C:\Users\l1768\Desktop\图片+视频\IMG_0931(20200127-163154).JPG`
	filePath = filepath.ToSlash(filePath)
	fileName := path.Base(filePath)
	fileType := utils.GetFileType(filePath)
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println(err)
		t.Fatal()
	}
	err = fileSvc.UploadFile(file, fileName, fileType, 4)
	if err != nil {
		fmt.Println(err)
		t.Fatal()
	}
}
