package ioc

import (
	"HelloCity/internal/repository/dao"
	"HelloCity/internal/utils"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	viper := utils.CreateConfig("config")
	prefix := "mysql."
	host := viper.GetString(prefix + "host")
	port := viper.GetInt(prefix + "port")
	user := viper.GetString(prefix + "username")
	passwd := viper.GetString(prefix + "password")
	dbname := viper.GetString(prefix + "database")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, passwd, host, port, dbname)
	//fmt.Println(dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&dao.User{})
	if err != nil {
		panic(err)
	}
	return db
}
