package ioc

import (
	"HelloCity/internal/service/oss"
	"HelloCity/internal/service/oss/qiniu"
	"HelloCity/internal/utils"
)

func NewOssService() oss.Service {
	accessKey := utils.Config.GetString("oss.qiniu.accessKey")
	secretKey := utils.Config.GetString("oss.qiniu.secretKey")
	return qiniu.NewService(accessKey, secretKey)
}
