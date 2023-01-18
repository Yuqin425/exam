package api

import (
	"LanShan/api/middlewares/jwt"
	"LanShan/dao/mysql"
	"LanShan/models"
	"LanShan/service"
	"LanShan/utils"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

func SignUpHandler(c *gin.Context) {
	// 1.获取请求参数 2.校验数据有效性
	var fo *models.RegisterForm
	if err := c.ShouldBindJSON(&fo); err != nil {
		// 请求参数有误，直接返回响应
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		// 判断err是不是 validator.ValidationErrors类型的errors
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 非validator.ValidationErrors类型错误直接返回
			utils.ResponseError(c, utils.CodeInvalidParams) // 请求参数错误
			return
		}
		// validator.ValidationErrors类型错误则进行翻译
		utils.ResponseErrorWithMsg(c, utils.CodeInvalidParams, removeTopStruct(errs.Translate(trans)))
		return
	}

	// 3.业务处理——注册用户
	if err := service.SignUp(fo); err != nil {
		zap.L().Error("service.signup failed", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserExit) {
			utils.ResponseError(c, utils.CodeUserExist)
			return
		}
		utils.ResponseError(c, utils.CodeServerBusy)
		return
	}
	//返回响应
	utils.ResponseSuccess(c, nil)
}

func LoginHandler(c *gin.Context) {
	// 获取请求参数及参数校验
	var u *models.LoginForm
	if err := c.ShouldBindJSON(&u); err != nil {
		// 请求参数有误，直接返回响应
		zap.L().Error("Login with invalid param", zap.Error(err))
		// 判断err是不是 validator.ValidationErrors类型的errors
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 非validator.ValidationErrors类型错误直接返回
			utils.ResponseError(c, utils.CodeInvalidParams) // 请求参数错误
			return
		}
		// validator.ValidationErrors类型错误则进行翻译
		utils.ResponseErrorWithMsg(c, utils.CodeInvalidParams, removeTopStruct(errs.Translate(trans)))
		return
	}
	// 2、业务逻辑处理——登录
	user, err := service.Login(u)
	if err != nil {
		zap.L().Error("service.Login failed", zap.String("username", u.UserName), zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotExit) {
			utils.ResponseError(c, utils.CodeUserNotExist)
			return
		} else if errors.Is(err, mysql.ErrorPasswordWrong) {
			utils.ResponseError(c, utils.CodeInvalidPassword)
			return
		}
		utils.ResponseError(c, utils.CodeServerBusy)
		return
	}
	// 3、返回响应
	utils.ResponseSuccess(c, gin.H{
		"user_id":       fmt.Sprintf("%d", user.UserID),
		"user_name":     user.UserName,
		"access_token":  user.AccessToken,
		"refresh_token": user.RefreshToken,
	})
}

func RefreshTokenHandler(c *gin.Context) {
	rt := c.Query("refresh_token")

	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		utils.ResponseErrorWithMsg(c, utils.CodeInvalidToken, "请求头缺少Auth Token")
		c.Abort()
		return
	}
	// 按空格分割
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		utils.ResponseErrorWithMsg(c, utils.CodeInvalidToken, "Token格式不对")
		c.Abort()
		return
	}
	aToken, rToken, err := jwt.RefreshToken(parts[1], rt)
	fmt.Println(err)
	c.JSON(http.StatusOK, gin.H{
		"access_token":  aToken,
		"refresh_token": rToken,
	})
}
