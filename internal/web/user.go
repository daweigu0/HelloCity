package web

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/go-querystring/query"
	"net/http"
)

type UserHandler struct {
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
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
	return
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
