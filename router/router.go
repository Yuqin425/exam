package router

import (
	"LanShan/api"
	"LanShan/api/middlewares"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetupRouter(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode) // 设置成发布模式
	}
	r := gin.New()

	v1 := r.Group("/api/v1")

	v1.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	v1.POST("/signup", api.SignUpHandler) // 注册业务路由
	v1.POST("/login", api.LoginHandler)   // 登录业务路由

	v1.GET("/community", api.CommunityHandler)           // 获取分类社区列表
	v1.GET("/community/:id", api.CommunityDetailHandler) // 根据ID查找社区详情

	v1.GET("/problem/:id", api.ProblemDetailHandler) // 查询问题详情
	v1.GET("/problems", api.ProblemListHandler)      // 分页展示问题列表

	v1.GET("/answers/:id", api.AnswerListHandler)  // 根据题目获取题解列表
	v1.GET("/answer/:id", api.AnswerDetailHandler) // 获取题解

	v1.Use(middlewares.JWTAuthMiddleware()) // 应用JWT认证中间件
	{

		v1.POST("/problem", api.CreateProblemHandler)            // 发布问题
		v1.GET("/problem/delete/:id", api.ProblemDeleteHandler)  // 删除问题
		v1.POST("/problem/update/:id", api.ProblemUpdateHandler) // 修改问题

		v1.POST("/answer", api.AnswerHandler)                  // 发布题解
		v1.GET("/answer/delete/:id", api.AnswerDeleteHandler)  // 删除题解
		v1.POST("/answer/update/:id", api.AnswerUpdateHandler) //  修改题解

	}

	return r
}
