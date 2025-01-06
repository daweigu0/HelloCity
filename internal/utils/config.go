package utils

import (
	"fmt"
	"github.com/spf13/viper"
	"path"
	"runtime"
)

var configPath = path.Dir(getCurrentPath()+"/../") + "/config/"

func getCurrentPath() string {
	_, filename, _, _ := runtime.Caller(0)
	return path.Dir(filename)
}
func CreateConfig(fileName string) *viper.Viper {
	config := viper.New()
	config.AddConfigPath(configPath)
	config.SetConfigName(fileName)
	config.SetConfigType("yaml")
	configFile := configPath + fileName + ".yaml"
	if err := config.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic(fmt.Errorf("找不到配置文件 %s ", configFile))
		} else {
			panic(fmt.Errorf("解析配置文件 %s 出错 %s ", configFile, err))
		}
	}
	return config
}
