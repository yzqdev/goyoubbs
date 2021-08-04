package model

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"goyoubbs/system"
	"os"
)

type Database struct {
	*gorm.DB
}

var DB *gorm.DB
var sqlType string

// Init Open mysql 连接
func init() {
	var db *gorm.DB
	g := system.LoadConfig()
	app := &system.Application{}
	app.Init(g, os.Args[0])
	if sqlType == "mysql" {
		//dsn := g.Mysql.User + ":" + g.Mysql.Pass + "@tcp(127.0.0.1:3306)/" + g.Mysql.Name + "?charset=utf8mb4&parseTime=True&loc=Local"
		//db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		//if err != nil {
		//	return
		//}
	} else {
		var err error
		db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
		if err != nil {
			fmt.Println("失败了")
		}
	}

	DB = db
	return
}

func GetDb() *gorm.DB {
	return DB
}
