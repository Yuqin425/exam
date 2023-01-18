package api

import (
	"LanShan/dao/mysql"
	"LanShan/models"
	"LanShan/service"
	"LanShan/utils"
	"LanShan/utils/snowflake"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

func AnswerHandler(c *gin.Context) {
	var answer models.Answer
	if err := c.BindJSON(&answer); err != nil {
		fmt.Println(err)
		utils.ResponseError(c, utils.CodeInvalidParams)
		return
	}
	// 生成ID
	answerID, err := snowflake.GetID()
	if err != nil {
		zap.L().Error("snowflake.GetID() failed", zap.Error(err))
		utils.ResponseError(c, utils.CodeServerBusy)
		return
	}
	// 获取作者ID，当前请求的UserID
	userID, err := getCurrentUserID(c)
	if err != nil {
		zap.L().Error("GetCurrentUserID() failed", zap.Error(err))
		utils.ResponseError(c, utils.CodeNotLogin)
		return
	}
	answer.AnswerID = answerID
	answer.AuthorID = userID

	// 创建题解
	if err := mysql.CreateAnswer(&answer); err != nil {
		zap.L().Error("mysql.CreateAnswer(&answer) failed", zap.Error(err))
		utils.ResponseError(c, utils.CodeServerBusy)
		return
	}
	utils.ResponseSuccess(c, nil)
}

func AnswerDetailHandler(c *gin.Context) {
	// 1、获取参数(从URL中获取id)
	answerIdStr := c.Param("id")
	answerId, err := strconv.ParseInt(answerIdStr, 10, 64)
	if err != nil {
		zap.L().Error("get answer detail with invalid param", zap.Error(err))
		utils.ResponseError(c, utils.CodeInvalidParams)
		return
	}

	// 2、根据id取出数据(查数据库)
	answer, err := mysql.GetAnswerById(answerId)
	if err != nil {
		zap.L().Error("mysql.GetAnswerById(answerId) failed", zap.Error(err))
		utils.ResponseError(c, utils.CodeServerBusy)
		return
	}

	// 3、返回响应
	utils.ResponseSuccess(c, answer)
}

func AnswerListHandler(c *gin.Context) {
	// 获取参数(从URL中获取id)
	problemIdStr := c.Param("id")
	problemId, err := strconv.ParseInt(problemIdStr, 10, 64)
	if err != nil {
		zap.L().Error("get problem detail with invalid param", zap.Error(err))
		utils.ResponseError(c, utils.CodeInvalidParams)
		return
	}
	// 获取分页参数
	page, size := getPageInfo(c)
	// 获取数据
	data, err := mysql.GetAnswerList(page, size, problemId)
	if err != nil {
		utils.ResponseError(c, utils.CodeServerBusy)
		return
	}
	utils.ResponseSuccess(c, data)
}

func AnswerUpdateHandler(c *gin.Context) {
	// 获取参数及校验参数
	var newAnswer models.Answer
	if err := c.ShouldBindJSON(&newAnswer); err != nil {
		zap.L().Debug("c.ShouldBindJSON(newAnswer) err", zap.Any("err", err))
		zap.L().Error("create answer with invalid param")
		utils.ResponseErrorWithMsg(c, utils.CodeInvalidParams, err.Error())
		return
	}

	// 获取参数(从URL中获取id)
	answerIdStr := c.Param("id")
	answerId, err := strconv.ParseInt(answerIdStr, 10, 64)
	pastAnswer, err := mysql.GetAnswerById(answerId)
	if err != nil {
		zap.L().Error("get answer detail with invalid param", zap.Error(err))
		utils.ResponseError(c, utils.CodeServerBusy)
	}

	// 获取作者ID
	UserID, err := getCurrentUserID(c)
	if err != nil {
		zap.L().Error("GetCurrentUserID() failed", zap.Error(err))
		utils.ResponseError(c, utils.CodeNotLogin)
		return
	}
	ok := UserID == pastAnswer.AuthorID
	if !ok {
		zap.L().Error("update problem with invalid param")
		utils.ResponseError(c, utils.CodeInvalidParams)
		return
	}

	answer, err := mysql.UpdateAnswer(&newAnswer)
	if err != nil {
		zap.L().Error("mysql.UpdateAnswer() failed")
		utils.ResponseError(c, utils.CodeServerBusy)
		return
	}

	// 3、返回响应
	utils.ResponseSuccess(c, answer)
}

func AnswerDeleteHandler(c *gin.Context) {
	// 获取参数(从URL中获取id)
	answerIdStr := c.Param("id")
	answerId, err := strconv.ParseInt(answerIdStr, 10, 64)
	answer, err := service.GetProblemById(answerId)
	if err != nil {
		zap.L().Error("get answer detail with invalid param", zap.Error(err))
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
	ok := UserID == answer.AuthorId
	if !ok {
		zap.L().Error("delete problem with invalid param")
		utils.ResponseError(c, utils.CodeInvalidParams)
		return
	}
	err = mysql.DeleteAnswer(answerId)
	if err != nil {
		zap.L().Error("mysql.DeleteAnswer() failed", zap.Error(err))
		utils.ResponseError(c, utils.CodeServerBusy)
		return
	}

	utils.ResponseSuccess(c, nil)
}
