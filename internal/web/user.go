package web

import (
	"HelloCity/ginx"
	"HelloCity/internal/errs"
	"HelloCity/internal/service"
	"HelloCity/internal/utils"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/go-querystring/query"
	"log"
	"net/http"
	"time"
)

type UserHandler struct {
	UserService service.UserService
}

func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{
		UserService: svc,
	}
}
func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	ug.POST("login", u.Login)
}

type loginReq struct {
	Code string `json:"code"`
}

// 登录注册接口
// @Tags 用户相关接口
// @Summary 用户登录注册接口
// @Description  登录成功返回的token放在响应的header的x-jwt-token里面，登录之后的后续访问需要带上token，放在请求的header里面的Authorization。
// @Accept       json
// @Produce      json
// @Param code body loginReq true "微信登录的临时登录凭证"
// @Success 200 {object} ginx.Result "登录成功"
// @Failure 401001 {object} ginx.Result "请求数据有误"
// @Failure 501001 {object} ginx.Result "登录失败"
// @Router /users/login [post]
func (u *UserHandler) Login(ctx *gin.Context) {
	var req loginReq
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(errs.UserInvalidInput, ginx.Result{
			Msg: "请求数据有误",
		})
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
	code2SessionResponse := u.code2Session(&code2SessionReqParams)
	if code2SessionResponse == nil || code2SessionResponse.ErrCode != 0 {
		ctx.JSON(errs.UserInternalServerError, ginx.Result{
			Msg: "登录失败",
		})
		if code2SessionResponse != nil {
			log.Println(fmt.Printf("请求微信code2Session接口失败，错误码：%d", code2SessionResponse.ErrCode))
		} else {
			log.Println("请求微信code2Session接口失败")
		}
		return
	}
	us, err := u.UserService.Login(ctx, code2SessionResponse.OpenId)
	if err != nil {
		ctx.JSON(errs.UserInternalServerError, ginx.Result{
			Msg: "登录失败",
		})
		return
	}
	log.Println("us:", us)
	rc := UserClaims{
		Uid:       us.Uid,
		UserAgent: ctx.GetHeader("User-Agent"),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Minute)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES512, rc)
	tokenString, err := token.SignedString(JWTKey)
	if err != nil {
		ctx.JSON(errs.UserInternalServerError, ginx.Result{
			Msg: "系统异常",
		})
		return
	}
	ctx.Header("x-jwt-token", tokenString)
	ctx.JSON(http.StatusOK, ginx.Result{
		Msg: "登录成功",
	})
}

type Code2SessionReqParams struct {
	Appid     string `url:"appid"`
	Secret    string `url:"secret"`
	JsCode    string `url:"js_code"`
	GrantType string `url:"grant_type"`
}
type UserClaims struct {
	jwt.RegisteredClaims
	Uid       uint64
	UserAgent string
}

type Code2SessionResponse struct {
	SessionKey string `json:"session_key"`
	UnionId    string `json:"Unionid"`
	ErrMsg     string `json:"errmsg"`
	OpenId     string `json:"openid"`
	ErrCode    int32  `json:"errcode"`
}

var JWTKey = []byte("k6CswdUm77WKcbM68UQUuxVsHSpTCwgA")

func (u *UserHandler) code2Session(reqParams *Code2SessionReqParams) *Code2SessionResponse {
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
