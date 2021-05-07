package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"ngxs.site/ginessential/common"
	"os"
)


func main (){
	InitConfig()
	common.InitDB()

	r := gin.Default()
	r = CollectRoute(r)

	//监听端口
	port := viper.GetString("server.port")
	if port != "" {
		panic(r.Run(":"+port))
	}
	panic(r.Run())
}

func InitConfig(){
	workDir, _ := os.Getwd()
	viper.SetConfigName("application")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir + "/config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}