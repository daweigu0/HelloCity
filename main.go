package main

import (
	_ "HelloCity/docs"
	"HelloCity/internal/utils"
	"fmt"
	"log"
)

// @title 你好同城后端API
// @accept json
// @produce	json
// @schemes	http https
func main() {
	server := InitWebServer()
	domain := utils.Config.GetString("nihaotongcheng.domain")
	port := utils.Config.GetString("nihaotongcheng.port")
	err := server.Run(fmt.Sprintf("%s:%s", domain, port))
	if err != nil {
		log.Panicf("服务器启动错误 %v\n", err)
	}
}
