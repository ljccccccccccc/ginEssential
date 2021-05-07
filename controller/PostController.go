package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"ngxs.site/ginessential/common"
	"ngxs.site/ginessential/model"
	"ngxs.site/ginessential/response"
	"ngxs.site/ginessential/vo"
	"strconv"
)

type IPostController interface {
	RestController
	PageList(ctx *gin.Context)
}

type PostController struct {
	DB *gorm.DB
}

func (p PostController) PageList(ctx *gin.Context) {
	// 获取分页参数
	pageNum,_ := strconv.Atoi(ctx.DefaultQuery("pageNum","1"))
	var pageSize int
	pageSize,_ = strconv.Atoi(ctx.DefaultQuery("pageSize","20"))
	// 分页
	var posts []model.Post
	p.DB.Order("created_at desc").Offset((pageNum -1 ) * pageSize).Limit(pageSize).Find(&posts);

	// 总条数
	var total int64
	p.DB.Model(model.Post{}).Count(&total)
	response.Success(ctx, gin.H{"data":posts,"total":total},"查询成功！")
}

func (p PostController) Create(ctx *gin.Context) {
	var requestPost vo.CreatePostRequest

	// 数据验证
	if err := ctx.ShouldBind(&requestPost); err != nil {
		response.Fail(ctx, nil,"数据验证不通过！")
		return
	}


	// 获取登陆用户user 需要加上中间件auth
	user, _ := ctx.Get("user")

	// 创建post文章
	post := model.Post{
		UserID:     user.(model.User).ID,
		CategoryId: requestPost.CategoryId,
		Title:      requestPost.Title,
		HeadImg:    requestPost.HeadImg,
		Content:    requestPost.Content,
	}

	// 插入数据
	if err:= p.DB.Create(&post).Error; err != nil {
		panic(err)
		return
	}else{
		response.Success(ctx, nil, "创建成功！")
	}

}

func (p PostController) Update(ctx *gin.Context) {
	var requestPost vo.CreatePostRequest
	// 数据验证
	if err := ctx.ShouldBind(&requestPost); err != nil {
		log.Println(err.Error())
		response.Fail(ctx, nil,"数据验证不通过！")
		return
	}

	// 获取path中的postid
	postId := ctx.Params.ByName("id")
	var post model.Post
	if err := p.DB.Where("id = ?",postId).First(&post).Error;err != nil {
		response.Fail(ctx,nil,"文章不存在！")
		return
	}

	// 判断当前用户是否为要修改文章的作者
	user, _ := ctx.Get("user")
	userId := user.(model.User).ID
	if userId != post.UserID {
		response.Fail(ctx, nil,"文章权限错误！请正常操作！")
		return
	}

	//更新文章
	//todo 这里的update参数有问题  先试下updates
	if err := p.DB.Model(&post).Updates(requestPost).Error;err!=nil{
		panic(err)
	}else{
		response.Success(ctx, gin.H{"post":post},"更新文章成功！")
	}

}

func (p PostController) Show(ctx *gin.Context) {
	//var requestPost vo.CreatePostRequest
	// 数据验证
	//if err := ctx.ShouldBind(&requestPost); err != nil {
	//	log.Println(err.Error())
	//	response.Fail(ctx, nil,"数据验证不通过！")
	//}

	// 获取path中的postid
	postId := ctx.Params.ByName("id")
	var post model.Post
    if 	err := p.DB.Preload("Category").Where("id = ?",postId).First(&post).Error;err != nil {
		response.Fail(ctx,nil,"文章不存在！")
		return
	}else{
		response.Success(ctx, gin.H{"post":post},"查询成功！")
	}
}

func (p PostController) Delete(ctx *gin.Context) {
	// 获取query参数
	postId := ctx.Params.ByName("id")
	var post model.Post

	// 判断文章是否存在
	if err := p.DB.Where("id = ?",postId).First(&post).Error;err != nil {
		panic(err)
		return
	}

	// 判断当前操作者是否为文章所有者
	user ,_ := ctx.Get("user")
	if user.(model.User).ID != post.UserID{
		response.Fail(ctx,nil,"权限错误！请合法操作！")
		return
	}

	// 验证成功就删除
	if err:=p.DB.Delete(&post).Error;err!=nil{
		if errors.Is(err,gorm.ErrRecordNotFound){
			response.Fail(ctx, nil, "删除文章失败,没有找到该文章")
			return
		}else{
			panic(err)
			return
		}
	}else {
		response.Success(ctx, nil, "删除文章成功！")
		return
	}
}

func NewPostController() IPostController{
	db := common.GetDB()
	db.AutoMigrate(model.Post{})
	
	return PostController{DB:db}
}
