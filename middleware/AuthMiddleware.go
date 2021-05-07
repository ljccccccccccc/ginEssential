package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ngxs.site/ginessential/common"
	"ngxs.site/ginessential/model"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//获取authorization header
		tokenString := ctx.GetHeader("Authorization")

		// validate token formate
		if tokenString =="" || !strings.HasPrefix(tokenString,"Bearer ") {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":401	,
				"msg":"权限不足",
			})
		}

		tokenString = tokenString[7:]

		token,claims,err := common.ParseToken(tokenString)
		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":401	,
				"msg":"权限不足",
			})
			ctx.Abort()
			return
		}

		// 验证通过后获取claim中的userID
		userId := claims.UserId
		DB := common.GetDB()
		var user model.User
		DB.First(&user, userId)

		//用户是否存在
		if user.ID == 0 {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":401	,
				"msg":"权限不足",
			})
			ctx.Abort()
			return
		}
		//用户存在 将user信息写入上下文
		ctx.Set("user",user)
		ctx.Next()
	}
}
