package power_wechat

import (
	"HelloCity/internal/service/wechat"
	"errors"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/miniProgram"
	"github.com/gin-gonic/gin"
)

type Service struct {
	mp *miniProgram.MiniProgram
}

func NewService(mp *miniProgram.MiniProgram) wechat.Service {
	return &Service{
		mp: mp,
	}
}
func (s *Service) Login(ctx *gin.Context, code string) (string, string, string, error) {
	resp, err := s.mp.Auth.Session(ctx, code)
	if err != nil {
		return "", "", "", err
	}
	if resp.ErrCode != 0 {
		return "", "", "", errors.New(resp.ErrMsg)
	}
	return resp.UnionID, resp.OpenID, resp.SessionKey, nil
}
