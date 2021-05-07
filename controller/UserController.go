package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"net/http"
	"ngxs.site/ginessential/common"
	"ngxs.site/ginessential/dto"
	"ngxs.site/ginessential/model"
	"ngxs.site/ginessential/response"
	"ngxs.site/ginessential/utils"
)

func Register (ctx *gin.Context) {
	DB := common.GetDB()

	// 使用map获取请求参数
	//var requestMap = make(map[string]string)
	//json.NewDecoder(ctx.Request.Body).Decode(&requestMap)

	//结构体获取参数
	//var requestUser = model.User{}
	//json.NewDecoder(ctx.Request.Body).Decode(&requestUser)

	// gin的bind获取参数
	var requestUser = model.User{}
	ctx.Bind(&requestUser)

	// 获取参数 名称、手机号和密码
	name := requestUser.Name
	telephone := requestUser.Telephone
	password := requestUser.Password
	fmt.Println(name,telephone,password)
	// 数据验证
	// map[string]interface{} 就是  gin.h
	if len(telephone) != 11 {
		response.Response(ctx, http.StatusUnprocessableEntity,422,nil,"手机号必须为11位!")
		return
	}

	if len(password) < 6 {
		response.Response(ctx, http.StatusUnprocessableEntity,422,nil,"密码不能少于6位!")
		return
	}

	if len(name) == 0 {
		name = utils.RandomString(10)
	}

	log.Println(name, telephone, password)

	// 判断手机号是否存在
	if isTelephoneExist(DB, telephone) {
		response.Response(ctx, http.StatusUnprocessableEntity,422,nil,"用户已经存在!")
		return
	}
	// 创建用户
	// 加密用户密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError,500,nil,"密码加密失败!")
		return
	}
	newUser := model.User{
		Name:      name,
		Telephone: telephone,
		Password:  string(hashedPassword),
	}

	DB.Create(&newUser)

	//返回结果
	response.Success(ctx, nil,"注册成功")
}

func Login (ctx *gin.Context) {
	db := common.GetDB()
	//获取参数
	telephone := ctx.PostForm("telephone")
	password := ctx.PostForm("password")

	// 数据验证
	if len(telephone) != 11 {
		response.Response(ctx, http.StatusUnprocessableEntity,422,nil,"手机号必须为11位!")
		return
	}
	if len(password) <6{
		response.Response(ctx, http.StatusUnprocessableEntity,422,nil,"密码不能小于6位!")
		return
	}
	// 判断手机号是否存在
	var user model.User
	db.Where("telephone = ?",telephone).First(&user)
	if user.ID == 0 {
		//用户不存在
		response.Response(ctx, http.StatusUnprocessableEntity,422,nil,"用户不存在!")
		return
	}

	// 判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password));err != nil {
		response.Response(ctx, http.StatusBadRequest,400,nil,"密码错误!")
		return
	}

	// 密码正确发放token给前端
	token,errtoken := common.ReleaseToken(user)
	if errtoken != nil{
		response.Response(ctx, http.StatusInternalServerError,500,nil,"生成token失败!")
		log.Printf("token generate error: %v", errtoken)
		return
	}

	response.Success(ctx, gin.H{"token":token},"ok")
}

func Info (ctx *gin.Context) {
	user, _ := ctx.Get("user")

	response.Success(ctx, gin.H{"user":dto.ToUserDto(user.(model.User))},"ok")
}

func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user model.User
	if res := db.Where("telephone = ?",telephone).First(&user);res.Error != nil {
		return false
	}
	return true
}




