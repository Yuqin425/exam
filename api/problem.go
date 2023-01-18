package api

import (
	"LanShan/models"
	"LanShan/service"
	"LanShan/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

func CreateProblemHandler(c *gin.Context) {
	// 1、获取参数及校验参数
	var problem models.Problem
	if err := c.ShouldBindJSON(&problem); err != nil {
		zap.L().Debug("c.ShouldBindJSON(problem) err", zap.Any("err", err))
		zap.L().Error("create problem with invalid param")
		utils.ResponseErrorWithMsg(c, utils.CodeInvalidParams, err.Error())
		return
	}

	// 获取作者ID
	userID, err := getCurrentUserID(c)
	if err != nil {
		zap.L().Error("GetCurrentUserID() failed", zap.Error(err))
		utils.ResponseError(c, utils.CodeNotLogin)
		return
	}
	problem.AuthorId = userID
	// 2、发布问题
	err = service.CreateProblem(&problem)
	if err != nil {
		zap.L().Error("service.CreateProblem failed", zap.Error(err))
		utils.ResponseError(c, utils.CodeServerBusy)
		return
	}
	// 3、返回响应
	utils.ResponseSuccess(c, nil)
}

// ProblemListHandler 问题列表
func ProblemListHandler(c *gin.Context) {
	// 获取分页参数
	page, size := getPageInfo(c)
	// 获取数据
	data, err := service.GetProblemList(page, size)
	if err != nil {
		utils.ResponseError(c, utils.CodeServerBusy)
		return
	}
	utils.ResponseSuccess(c, data)
}

// ProblemDetailHandler 根据Id查询问题
func ProblemDetailHandler(c *gin.Context) {
	// 1、获取参数(从URL中获取id)
	problemIdStr := c.Param("id")
	problemId, err := strconv.ParseInt(problemIdStr, 10, 64)
	if err != nil {
		zap.L().Error("get problem detail with invalid param", zap.Error(err))
		utils.ResponseError(c, utils.CodeInvalidParams)
		return
	}

	// 2、根据id取出id帖子数据(查数据库)
	problem, err := service.GetProblemById(problemId)
	if err != nil {
		zap.L().Error("service.GetProblem(problemID) failed", zap.Error(err))
		utils.ResponseError(c, utils.CodeServerBusy)
		return
	}

	// 3、返回响应
	utils.ResponseSuccess(c, problem)
}

func ProblemUpdateHandler(c *gin.Context) {
	// 获取参数及校验参数
	var newProblem models.Problem
	if err := c.ShouldBindJSON(&newProblem); err != nil {
		zap.L().Debug("c.ShouldBindJSON(problem) err", zap.Any("err", err))
		zap.L().Error("create problem with invalid param")
		utils.ResponseErrorWithMsg(c, utils.CodeInvalidParams, err.Error())
		return
	}

	// 获取参数(从URL中获取id)
	problemIdStr := c.Param("id")
	problemId, err := strconv.ParseInt(problemIdStr, 10, 64)
	pastProblem, err := service.GetProblemById(problemId)
	if err != nil {
		zap.L().Error("get problem detail with invalid param", zap.Error(err))
		utils.ResponseError(c, utils.CodeServerBusy)
	}

	// 获取作者ID
	UserID, err := getCurrentUserID(c)
	if err != nil {
		zap.L().Error("GetCurrentUserID() failed", zap.Error(err))
		utils.ResponseError(c, utils.CodeNotLogin)
		return
	}
	ok := UserID == pastProblem.AuthorId
	if !ok {
		zap.L().Error("update problem with invalid param")
		utils.ResponseError(c, utils.CodeInvalidParams)
		return
	}

	problem, err := service.UpdateProblem(&newProblem, problemId)
	if err != nil {
		zap.L().Error("service.UpdateProblem() failed")
		utils.ResponseError(c, utils.CodeServerBusy)
		return
	}

	// 3、返回响应
	utils.ResponseSuccess(c, problem)
}

func ProblemDeleteHandler(c *gin.Context) {
	// 获取参数(从URL中获取id)
	problemIdStr := c.Param("id")
	problemId, err := strconv.ParseInt(problemIdStr, 10, 64)
	problem, err := service.GetProblemById(problemId)
	if err != nil {
		zap.L().Error("get problem detail with invalid param", zap.Error(err))
		utils.ResponseError(c, utils.CodeServerBusy)
		return
	}

	// 获取作者ID
	UserID, err := getCurrentUserID(c)
	if err != nil {
		zap.L().Error("GetCurrentUserID() failed", zap.Error(err))
		utils.ResponseError(c, utils.CodeNotLogin)
		return
	}
	ok := UserID == problem.AuthorId
	if !ok {
		zap.L().Error("delete problem with invalid param")
		utils.ResponseError(c, utils.CodeInvalidParams)
		return
	}
	service.DeleteProblem(problemId)

	utils.ResponseSuccess(c, nil)
}
