package controller

import (
	"github.com/gin-gonic/gin"
	"ngxs.site/ginessential/model"
	"ngxs.site/ginessential/repository"
	"ngxs.site/ginessential/response"
	"ngxs.site/ginessential/vo"
	"strconv"
)

type ICategoryController interface {
	RestController
}

func NewCategoryController() ICategoryController {
	repository := repository.NewCategoryRepository()
	repository.DB.AutoMigrate(model.Category{})
	return CategoryController{Repository: repository}
}

type CategoryController struct {
	Repository repository.CategoryRepository
}

func (c CategoryController) Create(ctx *gin.Context) {
	var requestCategory vo.CreateCategoryRequest
	if err := ctx.ShouldBind(&requestCategory); err!=nil {
		response.Fail(ctx, nil, "数据验证错误，分类名称必填")
		return
	}else{
		//category := model.Category{Name: requestCategory.Name}
		category, err := c.Repository.Create(requestCategory.Name);
		if err != nil {
			response.Fail(ctx, nil, "创建失败")
			return
		}
		response.Success(ctx, gin.H{"category":category}, "ok")
	}
}

func (c CategoryController) Update(ctx *gin.Context) {
	//绑定body中的参数
	var requestCategory vo.CreateCategoryRequest
	if err := ctx.ShouldBind(&requestCategory); err!=nil {
		response.Fail(ctx, nil, "数据验证错误，分类名称必填")
		return
	}

	// 获取path中的参数
	categoryId,_ := strconv.Atoi(ctx.Params.ByName("id"))

	//var updateCategory model.Category
	//if err := c.DB.First(&updateCategory, categoryId).Error;err != nil {
	//	if errors.Is(err, gorm.ErrRecordNotFound) {
	//		response.Fail(ctx, nil, "分类不存在")
	//	}
	//}else{
	//	//更新
	//	// 这里可以传入三种类型
	//	// map
	//	// struct
	//	// name value
	//	c.DB.Model(&updateCategory).Update("name",requestCategory.Name)
	//	response.Success(ctx, gin.H{"category":requestCategory},"修改成功！")
	//}
	updateCategory, err := c.Repository.SelectById(categoryId)
	if err != nil {
		response.Fail(ctx,nil,"分类不存在")
	}else{
		category, err := c.Repository.Update(*updateCategory, requestCategory.Name)
		if err != nil {
			panic(err)
		}
		response.Success(ctx, gin.H{"category": category},"修改成功！")
	}

}

func (c CategoryController) Show(ctx *gin.Context) {
	// 获取path中的参数
	categoryId,_ := strconv.Atoi(ctx.Params.ByName("id"))

	category ,err := c.Repository.SelectById(categoryId)
	if err != nil {
		response.Fail(ctx, nil, "没有此分类!")
		return
	}else{
		response.Success(ctx, gin.H{"category":category},"ok")
	}
}

func (c CategoryController) Delete(ctx *gin.Context) {
	categoryId,_ := strconv.Atoi(ctx.Params.ByName("id"))
	if err := c.Repository.DeleteById(categoryId);err != nil {
		response.Fail(ctx, nil,"删除失败！")
	}else{
		response.Success(ctx, nil, "删除成功！")
	}
}

