package oss

type Service interface {
	CreateBucket(bucketName string) (any, error)
	GetUploadCallBackToken(bucketName, callBackUrl, callBackBody, callBackBodyType string) (string, error)
}
