package qiniu

import (
	"context"
	"github.com/qiniu/go-sdk/v7/storagev2/credentials"
	"github.com/qiniu/go-sdk/v7/storagev2/http_client"
	"github.com/qiniu/go-sdk/v7/storagev2/objects"
	"github.com/qiniu/go-sdk/v7/storagev2/uptoken"
	"time"
)

type Service struct {
	accessKey string
	secretKey string
}

func (s *Service) CreateBucket(bucketName string) (any, error) {
	mac := credentials.NewCredentials(s.secretKey, s.secretKey)
	objectsManager := objects.NewObjectsManager(&objects.ObjectsManagerOptions{
		Options: http_client.Options{Credentials: mac},
	})
	bucket := objectsManager.Bucket(bucketName)
	return bucket, nil
}

func (s *Service) GetUploadCallBackToken(bucketName, callBackUrl, callBackBody, callBackBodyType string) (string, error) {
	mac := credentials.NewCredentials(s.accessKey, s.secretKey)
	putPolicy, err := uptoken.NewPutPolicy(bucketName, time.Now().Add(1*time.Hour))
	if err != nil {
		return "", err
	}
	putPolicy.SetCallbackUrl(callBackUrl).
		SetCallbackBody(callBackBody).
		SetCallbackBodyType(callBackBodyType)
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
