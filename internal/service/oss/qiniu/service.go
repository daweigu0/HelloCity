package qiniu

import (
	"HelloCity/internal/utils"
	"context"
	"errors"
	"fmt"
	"github.com/qiniu/go-sdk/v7/storagev2/credentials"
	"github.com/qiniu/go-sdk/v7/storagev2/http_client"
	"github.com/qiniu/go-sdk/v7/storagev2/uploader"
	"github.com/qiniu/go-sdk/v7/storagev2/uptoken"
	"io"
	"path"
	"time"
)

var (
	ErrParamNotMatch  = errors.New("参数类型不匹配")
	ErrFileUploadFail = errors.New("文件上传失败")
)

type Service struct {
	accessKey string
	secretKey string
}

func (s *Service) UploadFile(reader io.Reader, fileName, fileType string, uid uint64) (string, error) {
	// 自定义返回值结构体
	type MyPutRet struct {
		Key    string
		Hash   string
		Fsize  int
		Bucket string
	}
	mac := credentials.NewCredentials(s.accessKey, s.secretKey)
	bucketName := utils.Config.GetString("oss.qiniu.bucketName")
	key := fmt.Sprintf("%s/%d/%s", fileType, uid, path.Base(fileName))
	uploadManager := uploader.NewUploadManager(&uploader.UploadManagerOptions{
		Options: http_client.Options{
			Credentials: mac,
		},
	})
	putPolicy, err := uptoken.NewPutPolicy(bucketName, time.Now().Add(1*time.Hour))
	if err != nil {
		return "", ErrFileUploadFail
	}
	var ret MyPutRet
	putPolicy.SetReturnBody(`{"key":"$(key)","hash":"$(etag)","fsize":$(fsize),"bucket":"$(bucket)"`)
	err = uploadManager.UploadReader(context.Background(), reader, &uploader.ObjectOptions{
		BucketName: bucketName,
		ObjectName: &key,
		FileName:   fileName,
	}, &ret)
	if err != nil {
		return "", ErrFileUploadFail
	}
	return ret.Key, nil
}

type GetUploadTokenParam struct {
	BucketName       string
	CallBackUrl      string
	CallBackBody     string
	CallBackBodyType string
	SaveKey          string
}

func NewGetUploadTokenParam(bucketName, callBackUrl, callBackBody, callBackBodyType, saveKey string) *GetUploadTokenParam {
	return &GetUploadTokenParam{
		BucketName:       bucketName,
		CallBackUrl:      callBackUrl,
		CallBackBody:     callBackBody,
		CallBackBodyType: callBackBodyType,
		SaveKey:          saveKey,
	}
}
func (s *Service) GetUploadToken(param any) (string, error) {
	p, ok := param.(*GetUploadTokenParam)
	if !ok {
		return "", ErrParamNotMatch
	}
	mac := credentials.NewCredentials(s.accessKey, s.secretKey)
	putPolicy, err := uptoken.NewPutPolicy(p.BucketName, time.Now().Add(1*time.Hour))
	if err != nil {
		return "", err
	}
	//这个地方以后是否可以用反射优化？避免加一个参数就要修改代码
	putPolicy.SetCallbackUrl(p.CallBackUrl).
		SetCallbackBody(p.CallBackBody).
		SetCallbackBodyType(p.CallBackBodyType).
		SetSaveKey(p.SaveKey)
	upToken, err := uptoken.NewSigner(putPolicy, mac).GetUpToken(context.Background())
	if err != nil {
		return "", err
	}
	return upToken, nil
}

func NewService(accessKey, secretKey string) *Service {
	return &Service{
		accessKey: accessKey,
		secretKey: secretKey,
	}
}
