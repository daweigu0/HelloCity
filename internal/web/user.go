package web

import (
	"HelloCity/internal/service"
	"encoding/json"
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
func (u *UserHandler) Login(ctx *gin.Context) {
	type Req struct {
		Code string `json:"code"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	code2SessionReqParams := Code2SessionReqParams{
		JsCode: req.Code,
	}
	code2SessionResponse := u.code2Session(&code2SessionReqParams)
	//下面逻辑需要继续完善
	if code2SessionResponse.SessionKey == "" {

	}
	us, err := u.UserService.Login(ctx, req.Code)
	if err != nil {
		ctx.String(http.StatusOK, "登录失败")
		return
	}
	log.Println("us:", us)
	//ssid := uuid.New().String()
	rc := UserClaims{
		Uid:       us.Uid,
		UserAgent: ctx.GetHeader("User-Agent"),
		//Ssid: ssid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Minute)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES512, rc)
	tokenString, err := token.SignedString(JWTKey)
	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}
	ctx.Header("x-jwt-token", tokenString)
	ctx.String(http.StatusOK, "登录成功")
}
func (u *UserHandler) Hello(ctx *gin.Context) {
	ctx.String(http.StatusOK, "hello world")
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
		return nil
	}
	var data []byte
	_, err = resp.Body.Read(data)
	if err != nil {
		return nil
	}
	code2SessionResponse := new(Code2SessionResponse)
	err = json.Unmarshal(data, code2SessionResponse)
	if err != nil {
		return nil
	}
	return code2SessionResponse
}
