package web

import (
	"HelloCity/internal/domain"
	"HelloCity/internal/global/consts"
	"HelloCity/internal/service"
	"HelloCity/internal/utils"
	"HelloCity/internal/utils/check"
	"HelloCity/internal/utils/response"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/go-querystring/query"
	"log"
	"net/http"
	"time"
)

type UserHandler struct {
	svc service.UserService
}

func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{
		svc: svc,
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
	viper := utils.CreateConfig("config")
	prefix := "wechat."
	appid := viper.GetString(prefix + "appid")
	secret := viper.GetString(prefix + "secret")
	code2SessionReqParams := Code2SessionReqParams{
		JsCode:    req.Code,
		Appid:     appid,
		Secret:    secret,
		GrantType: "authorization_code",
	}
	code2SessionResponse := h.code2Session(&code2SessionReqParams)
	if code2SessionResponse == nil || code2SessionResponse.ErrCode != 0 {
		response.Fail(ctx, consts.CurdLoginFailCode, consts.CurdLoginFailMsg, nil)
		if code2SessionResponse != nil {
			log.Println(fmt.Printf("请求微信code2Session接口失败，错误码：%d", code2SessionResponse.ErrCode))
		} else {
			log.Println("请求微信code2Session接口失败")
		}
		return
	}
	us, err := h.svc.Login(ctx, code2SessionResponse.OpenId)
	if errors.Is(err, service.ErrInvalidUser) {
		response.Fail(ctx, consts.CurdLoginFailCode, "用户不存在，请注册", gin.H{
			"openid": code2SessionResponse.OpenId, //这个地方可能有安全问题，后续需要解决
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

type Code2SessionReqParams struct {
	Appid     string `url:"appid"`
	Secret    string `url:"secret"`
	JsCode    string `url:"js_code"`
	GrantType string `url:"grant_type"`
}

type Code2SessionResponse struct {
	SessionKey string `json:"session_key"`
	UnionId    string `json:"Unionid"`
	ErrMsg     string `json:"errmsg"`
	OpenId     string `json:"openid"`
	ErrCode    int32  `json:"errcode"`
}

func (h *UserHandler) code2Session(reqParams *Code2SessionReqParams) *Code2SessionResponse {
	url := "https://api.weixin.qq.com/sns/jscode2session?"
	v, err := query.Values(reqParams)
	if err != nil {
		return nil
	}
	url += v.Encode()
	resp, err := http.Get(url)
	if err != nil {
		log.Println(fmt.Sprintf("请求code2Session接口失败，%v", err))
		return nil
	}
	defer resp.Body.Close() // 确保在函数退出时关闭响应体
	var data Code2SessionResponse
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&data)
	if err != nil {
		log.Println(fmt.Sprintf("解析json数据错误，%v", err))
		return nil
	}
	return &data
}

type SignUpReq struct {
	Mobile   string `json:"mobile"`
	NickName string `json:"nick_name"`
	Avatar   string `json:"avatar"`
	OpenId   string `json:"openid"`
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
	if !check.URL(signUpReq.Avatar) {
		response.Fail(ctx, http.StatusBadRequest, "头像url不正确", nil)
		return
	}
	//需要检查openid
	err := h.svc.SignUp(ctx, domain.User{
		Mobile:   signUpReq.Mobile,
		NickName: signUpReq.NickName,
		Avatar:   signUpReq.Avatar,
		OpenID:   signUpReq.OpenId,
	})
	switch err {
	case nil:
		response.Success(ctx, consts.CurdRegisterOkMsg, nil)
	case service.ErrDuplicateMobile:
		response.Fail(ctx, consts.CurdRegisterFailCode, "手机号冲突，请更换一个", nil)
	default:
		response.ErrorSystem(ctx, "", nil)
	}
}

// 用户界面接口
// @Tags 用户相关接口
// @Summary 用户界面接口
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
	user, err1 := h.svc.Profile(ctx, usClaims.Uid)
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
func (h *UserHandler) Edit(ctx *gin.Context) {
	uc, ok := ctx.MustGet("user").(utils.UserClaims)
	if !ok {
		log.Println("ctx 未找到用户")
		return
	}
	type editReq struct {
		Name          string `json:"name"`
		Gender        string `json:"gender"`
		Constellation int8   `json:"constellation"`
		Province      string `json:"province"`
		City          string `json:"city"`
		Signature     string `json:"signature"`
	}
	var req editReq
	if err := ctx.Bind(&req); err != nil {
		response.ErrorParam(ctx, err)
	}
	err := h.svc.Edit(ctx, uc.Uid, domain.User{
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
