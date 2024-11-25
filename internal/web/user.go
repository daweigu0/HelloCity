package web

import (
	"HelloCity/internal/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
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
	ug.GET("hello", u.Hello)
}
func (u *UserHandler) Login(ctx *gin.Context) {
	type Req struct {
		Code string `json:"code"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	us, err := u.UserService.Login(ctx, req.Code)
	if err != nil {
		ctx.String(http.StatusOK, "登录失败")
		return
	}
	log.Println("us:", us)
	ctx.String(http.StatusOK, "登录成功")
}
func (u *UserHandler) Hello(ctx *gin.Context) {
	ctx.String(http.StatusOK, "hello world")
	return
}
