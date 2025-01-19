package web

import (
	"HelloCity/internal/domain"
	"HelloCity/internal/global/consts"
	"HelloCity/internal/service/oss"
	"HelloCity/internal/service/oss/qiniu"
	"HelloCity/internal/utils"
	"HelloCity/internal/utils/response"
	"fmt"
	"github.com/gin-gonic/gin"
	uuid2 "github.com/google/uuid"
	"log"
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
	fg.POST("upload_token", h.UploadToken)
}

type myPutRet struct {
	Key    string
	Hash   string
	Fsize  int
	Bucket string
}

// CallBack 这个回调函数需在公网服务器上测，需完善
func (h *FileHandler) CallBack(ctx *gin.Context) {
	var ret myPutRet
	if err := ctx.Bind(&ret); err != nil {
		response.ErrorParam(ctx, nil)
		return
	}
	fmt.Println(ret)
}

type reqUploadToken struct {
	FileType string `json:"file_type"`
	FileExt  string `json:"file_ext"`
}

// 文件接口
// @Tags 文件相关接口
// @Summary 文件上传token获取接口
// @Description	文件上传token获取接口
// @Accept json
// @Produce json
// @Param UploadToken body reqUploadToken true "上传文件token获取参数"
// @Success 200 {object} ginx.Result "{"code":xxx,"data":{},"msg":"xxx"}"
// @Router /files/upload_token [post]
func (h *FileHandler) UploadToken(ctx *gin.Context) {
	req := new(reqUploadToken)
	if err := ctx.Bind(req); err != nil {
		response.ErrorParam(ctx, nil)
		return
	}
	fileTypes := []string{"image", "video", "audio"}
	legalType := false
	for _, fileType := range fileTypes {
		if fileType == req.FileType {
			legalType = true
			break
		}
	}
	if !legalType {
		response.Fail(ctx, consts.ValidatorParamsCheckFailCode, "文件类型不对", nil)
		return
	}
	if req.FileExt == "" {
		response.Fail(ctx, consts.ValidatorParamsCheckFailCode, "文件拓展名为空", nil)
		return
	}
	user := ctx.Value("user").(domain.User)
	uuid, err := uuid2.NewUUID()
	if err != nil {
		log.Println(fmt.Sprintf("uuid生成错误 %v", err))
		response.ErrorSystem(ctx, "", nil)
		return
	}
	fileName := uuid.String() + "." + req.FileExt
	saveKey := fmt.Sprintf("%s/%d/%s", req.FileType, user.ID, fileName)
	bucketName := utils.Config.GetString("oss.qiniu.bucketName")
	callBackUrl := utils.Config.GetString("nihaotongcheng.domain") + "files/callback"
	callBackBody := `{"key":"$(key)","hash":"$(etag)","fsize":$(fsize),"bucket":"$(bucket)"}`
	callBackBodyType := "application/json"
	param := qiniu.NewGetUploadTokenParam(bucketName, callBackUrl, callBackBody, callBackBodyType, saveKey)
	token, err := h.svc.GetUploadToken(param)
	if err != nil {
		log.Printf("上传文件生成token错误 %v\n", err)
		response.ErrorSystem(ctx, "", nil)
		return
	}
	response.Success(ctx, consts.CurdStatusOkMsg, gin.H{
		"upload_token": token,
	})
}
