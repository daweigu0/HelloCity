package oss

type Service interface {
	GetUploadToken(param any) (string, error)
	UploadFile(data []byte) error
}
