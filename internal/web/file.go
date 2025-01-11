package web

import (
	"HelloCity/internal/service/oss"
	"HelloCity/internal/utils/response"
	"github.com/gin-gonic/gin"
)

type FileHandler struct {
	svc oss.Service
}

func NewFileHandler(svc oss.Service) *FileHandler {
	return &FileHandler{
		svc: svc,
	}
}

func (h *FileHandler) RegisterRoutes(server *gin.Engine) {
	fg := server.Group("/files")
	fg.POST("callback", h.CallBack)
	fg.GET("upload_callback_token", h.UploadCallBackToken)
}

type myPutRet struct {
	Key    string
	Hash   string
	Fsize  int
	Bucket string
	Name   string
}

func (h *FileHandler) CallBack(ctx *gin.Context) {
	var ret myPutRet
	if err := ctx.Bind(&ret); err != nil {
		response.ErrorParam(ctx, nil)
		return
	}

}

func (h *FileHandler) UploadCallBackToken(ctx *gin.Context) {
	//TODO: implement
}
