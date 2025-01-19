package ioc

import (
	"HelloCity/internal/utils"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/miniProgram"
)

func NewWechatService() *miniProgram.MiniProgram {
	mp, err := miniProgram.NewMiniProgram(&miniProgram.UserConfig{
		AppID:     utils.Config.GetString("wechat.appid"),  // 小程序appid
		Secret:    utils.Config.GetString("wechat.secret"), // 小程序app secret
		HttpDebug: true,
		Log: miniProgram.Log{
			Level: "debug",
			// 可以重定向到你的目录下，如果设置File和Error，默认会在当前目录下的wechat文件夹下生成日志
			//File:   "/Users/user/wechat/mini-program/info.log",
			//Error:  "/Users/user/wechat/mini-program/error.log",
			Stdout: true, //  是否打印在终端
		},
	})
	if err != nil {
		panic(err)
	}
	return mp
}
