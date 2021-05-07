package controller

import (
	"github.com/gin-gonic/gin"
	"ngxs.site/ginessential/response"
)

type IFileController interface {
	RestController
}

type FileController struct {

}

func (f FileController) Create(ctx *gin.Context) {
	uploadedFile, err := ctx.FormFile("uploadedfile")
	if err != nil {
		panic(err)
	}
	if err := ctx.SaveUploadedFile(uploadedFile, uploadedFile.Filename);err != nil{
		panic(err)
	}
	response.Success(ctx,nil,"文件上传成功！")
}

func (f FileController) Update(ctx *gin.Context) {
	panic("implement me")
}

func (f FileController) Show(ctx *gin.Context) {
	panic("implement me")
}

func (f FileController) Delete(ctx *gin.Context) {
	panic("implement me")
}

func NewFileController() IFileController {
	return FileController{}
}
