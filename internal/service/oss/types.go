package oss

import "io"

type Service interface {
	GetUploadToken(param any) (string, error)
	UploadFile(reader io.Reader, fileName, fileType string, uid uint64) (string, error)
}
