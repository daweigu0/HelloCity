package qiniu

import (
	"HelloCity/internal/utils"
	"bytes"
	"context"
	"errors"
	"github.com/qiniu/go-sdk/v7/storagev2/credentials"
	"github.com/qiniu/go-sdk/v7/storagev2/http_client"
	"github.com/qiniu/go-sdk/v7/storagev2/uploader"
	"github.com/qiniu/go-sdk/v7/storagev2/uptoken"
	"time"
)

var (
	ErrParamNotMatch = errors.New("参数类型不匹配")
)

type Service struct {
	accessKey string
	secretKey string
}

func (s *Service) UploadFile(data []byte) error {
	mac := credentials.NewCredentials(s.accessKey, s.secretKey)
	bucket := utils.Config.GetString("oss.qiniu.bucketName")
	key := "github-x.png"
	reader := bytes.NewReader(data)
	uploadManager := uploader.NewUploadManager(&uploader.UploadManagerOptions{
		Options: http_client.Options{
			Credentials: mac,
		},
	})
	err := uploadManager.UploadReader(context.Background(), reader, &uploader.ObjectOptions{
		BucketName: bucket,
		ObjectName: &key,
		CustomVars: map[string]string{
			"name": "github logo",
		},
		FileName: "",
	}, nil)
	if err != nil {
		return err
	}
	return nil
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
