package common

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/url"
	"ngxs.site/ginessential/model"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	//driverName := "mysql"

	/*
		host := "172.16.80.25"
		port := "3306"
		database := "ginessential"
		username := "root"
		password := "root"
		charset := "utf8mb4"
	*/
	host := viper.GetString("datasource.host")
	port := viper.GetString("datasource.port")
	database := viper.GetString("datasource.database")
	username := viper.GetString("datasource.username")
	password := viper.GetString("datasource.password")
	charset := viper.GetString("datasource.charset")
	loc := viper.GetString("datasource.loc")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true&loc=%s",
		username,
		password,
		host,
		port,
		database,
		charset,
		url.QueryEscape(loc),
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database, err: " + err.Error())
	}

	err1 := db.AutoMigrate(&model.User{})
	if err1 != nil {
		return nil
	}
	DB = db
	return db
}

func GetDB() *gorm.DB {
	return DB
}
