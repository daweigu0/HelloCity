package web

import (
	"HelloCity/internal/domain"
	"HelloCity/internal/global/consts"
	"HelloCity/internal/service"
	"HelloCity/internal/service/oss"
	"HelloCity/internal/service/wechat"
	"HelloCity/internal/utils"
	"HelloCity/internal/utils/check"
	"HelloCity/internal/utils/response"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"time"
)

var (
	prefixSignup = "signup"
)

type UserHandler struct {
	userSvc   service.UserService
	tokenSvc  service.TokenService
	ossSvc    oss.Service
	wechatSvc wechat.Service
}

func NewUserHandler(userSvc service.UserService,
	tokenSvc service.TokenService,
	ossSvc oss.Service,
	wechatSvc wechat.Service,
) *UserHandler {
	return &UserHandler{
		userSvc:   userSvc,
		tokenSvc:  tokenSvc,
		ossSvc:    ossSvc,
		wechatSvc: wechatSvc,
	}
}
func (h *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	ug.POST("login", h.Login)
	ug.POST("signup", h.SignUp)
	ug.POST("profile", h.Profile)
	ug.POST("edit", h.Edit)
}

type loginReq struct {
	Code string `json:"code"`
}

// 登录接口
// @Tags 用户相关接口
// @Summary 用户登录接口
// @Description	登录成功返回的token放在响应的header的x-jwt-token里面，登录之后的后续访问需要带上token，放在请求的header里面的Authorization。
// @Accept json
// @Produce json
// @Param login body loginReq true "微信登录的临时登录凭证"
// @Success 200 {object} ginx.Result "{"code":xxx,"data":{},"msg":"xxx"}"
// @Router /users/login [post]
func (h *UserHandler) Login(ctx *gin.Context) {
	var req loginReq
	if err := ctx.Bind(&req); err != nil {
		response.ErrorParam(ctx, err)
		return
	}
	_, openid, _, err := h.wechatSvc.Login(ctx, req.Code)
	if err != nil {
		log.Printf("访问微信code2session接口错误 %v", err)
		response.Fail(ctx, consts.CurdLoginFailCode, consts.CurdLoginFailMsg, nil)
		return
	}
	us, err := h.userSvc.Login(ctx, openid)
	if errors.Is(err, service.ErrInvalidUser) {
		signupToken := utils.RandStr(16)
		err = h.tokenSvc.Set(ctx, prefixSignup, signupToken, openid)
		if err != nil {
			log.Printf("redis缓存数据错误 %v\n", err)
			response.ErrorSystem(ctx, "", nil)
			return
		}
		response.Fail(ctx, consts.CurdLoginFailCode, "用户不存在，请注册", gin.H{
			"signup_token": signupToken, //这个地方可能有安全问题，后续需要解决
		})
		return
	}
	if err != nil {
		response.Fail(ctx, consts.CurdLoginFailCode, consts.CurdLoginFailMsg, nil)
		return
	}
	log.Println("us:", us)
	uc := utils.UserClaims{
		Uid:       us.ID,
		UserAgent: ctx.GetHeader("User-Agent"),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Minute)),
		},
	}
	tokenString, err := utils.GenerateToken(&uc)
	if err != nil {
		response.ErrorSystem(ctx, "", nil)
		log.Println(err)
		return
	}
	ctx.Header("x-jwt-token", tokenString)
	response.Success(ctx, consts.CurdStatusOkMsg, gin.H{"token": tokenString})
}

type SignUpReq struct {
	Mobile      string `json:"mobile"`
	NickName    string `json:"nick_name"`
	SignupToken string `json:"signup_token"`
}

// 注册接口
// @Tags 用户相关接口
// @Summary 用户注册接口
// @Description	用户注册接口
// @Accept json
// @Produce json
// @Param signup body SignUpReq true "注册参数"
// @Success 200 {object} ginx.Result "{"code":xxx,"data":{},"msg":"xxx"}"
// @Router /users/signup [post]
func (h *UserHandler) SignUp(ctx *gin.Context) {
	var signUpReq SignUpReq
	if err := ctx.Bind(&signUpReq); err != nil {
		log.Println(err)
		response.ErrorParam(ctx, nil)
		return
	}
	if signUpReq.NickName == "" {
		response.Fail(ctx, http.StatusBadRequest, "昵称不能为空", nil)
		return
	}
	if !check.CNMobile(signUpReq.Mobile) {
		response.Fail(ctx, http.StatusBadRequest, "手机号不正确", nil)
		return
	}
	if signUpReq.SignupToken == "" {
		response.Fail(ctx, http.StatusBadRequest, "signup_token不能为空", nil)
		return
	}
	openid, err := h.tokenSvc.Get(ctx, prefixSignup, signUpReq.SignupToken)
	fmt.Println("openid:", openid)
	if err != nil {
		log.Printf("从redis中获取值错误 %v\n", err)
		response.Fail(ctx, http.StatusBadRequest, "signup_token错误", nil)
		return
	}
	err = h.userSvc.SignUp(ctx, domain.User{
		Mobile:   signUpReq.Mobile,
		NickName: signUpReq.NickName,
		OpenID:   openid,
	})
	//if err == nil { //用户头像的上传是否可以优化？
	//	avatarData, err := utils.Base64Decode(signUpReq.Avatar)
	//	if err != nil {
	//		log.Printf("头像解码错误 %v\n", err)
	//		response.Fail(ctx, http.StatusBadRequest, "头像解码错误", nil)
	//		return
	//	}
	//	uuid, err := uuid2.NewUUID()
	//	if err != nil {
	//		log.Printf("uuid生成错误 %v\n", err)
	//		response.ErrorSystem(ctx, "", nil)
	//		return
	//	}
	//	fileType := utils.GetFileType(avatarData)
	//	fileName := uuid.String() + "/" + fileType
	//	u, err := h.userSvc.FindUserByOpenID(ctx, openid)
	//	if err != nil {
	//		log.Printf("根据id查询用户错误 %v\n", err)
	//		response.Fail(ctx, http.StatusBadRequest, "头像上传失败", nil)
	//		return
	//	}
	//	avatar, err := h.ossSvc.UploadFile(bytes.NewReader(avatarData), fileName, fileType, u.ID)
	//	if err != nil {
	//		log.Printf("头像上传错误 %v\n", err)
	//		response.Fail(ctx, http.StatusBadRequest, "头像上传失败", nil)
	//		return
	//	}
	//	u.Avatar = avatar
	//	err = h.userSvc.UpdateNonSensitiveInfo(ctx, u)
	//	if err != nil {
	//		log.Printf("更新用户头像错误 %v\n", err)
	//		response.Fail(ctx, http.StatusBadRequest, "头像上传失败", nil)
	//		return
	//	}
	//}
	switch err {
	case nil:
		response.Success(ctx, consts.CurdRegisterOkMsg, nil)
	case service.ErrDuplicateMobile:
		response.Fail(ctx, consts.CurdRegisterFailCode, "手机号冲突，请更换一个", nil)
	default:
		response.ErrorSystem(ctx, "", nil)
	}
}

// 用户信息接口
// @Tags 用户相关接口
// @Summary 用户信息接口
// @Description
// @Accept json
// @Produce json
// @Success 200 {object} ginx.Result "{"code":xxx,"data":{},"msg":"xxx"}"
// @Router /users/profile [post]
func (h *UserHandler) Profile(ctx *gin.Context) {
	usClaims, err := ctx.MustGet("user").(utils.UserClaims)
	if err == false {
		log.Println("ctx中未存放user")
		return
	}
	user, err1 := h.userSvc.Profile(ctx, usClaims.Uid)
	if err1 != nil {
		response.Fail(ctx, consts.CurdSelectFailCode, consts.CurdSelectFailMsg, nil)
	}
	type resp struct {
		Id            uint64 `json:"id"`
		UserName      string `json:"username"`
		Avartar       string `json:"avartar"`
		ThumbsCount   int64  `json:"thumbsCount"`
		FansCount     int64  `json:"fansCount"`
		FollowerCount int64  `json:"followerCount"`
		Signature     string `json:"signature"`
		Constellation int8   `json:"constellation"`
		Province      string `json:"province"`
		City          string `json:"city"`
	}
	Re := &resp{
		Id:            user.ID,
		UserName:      user.NickName,
		Avartar:       user.Avatar,
		ThumbsCount:   user.ThumbsCount,
		FansCount:     user.FollowerCount,
		FollowerCount: user.FollowerCount,
		Signature:     user.Signature,
		Constellation: user.Constellation,
		Province:      user.Province,
		City:          user.City,
	}
	response.Success(ctx, consts.CurdStatusOkMsg, Re)
}

type editReq struct {
	Name          string `json:"name"`
	Gender        string `json:"gender"`
	Constellation int8   `json:"constellation"`
	Province      string `json:"province"`
	City          string `json:"city"`
	Signature     string `json:"signature"`
}

// 用户信息更新接口
// @Tags 用户相关接口
// @Summary 用户信息更新接口
// @Description
// @Accept json
// @Produce json
// @Param edit body editReq true "用户信息更新参数"
// @Success 200 {object} ginx.Result "{"code":xxx,"data":{},"msg":"xxx"}"
// @Router /users/profile [post]
func (h *UserHandler) Edit(ctx *gin.Context) {
	uc, ok := ctx.MustGet("user").(utils.UserClaims)
	if !ok {
		log.Println("ctx 未找到用户")
		return
	}
	var req editReq
	if err := ctx.Bind(&req); err != nil {
		response.ErrorParam(ctx, err)
	}
	err := h.userSvc.Edit(ctx, uc.Uid, domain.User{
		NickName:      req.Name,
		Gender:        req.Gender,
		Constellation: req.Constellation,
		Province:      req.Province,
		City:          req.City,
		Signature:     req.Signature,
	})
	if err != nil {
		log.Println(err)
		response.Fail(ctx, consts.CurdUpdateFailCode, consts.CurdUpdateFailMsg, nil)
	}
	response.Success(ctx, consts.CurdStatusOkMsg, nil)
}
