package api

import (
	"LanShan/service"
	"LanShan/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

func CommunityHandler(c *gin.Context) {
	// 查询到所有的社区(community_id,community_name)以列表的形式返回
	communityList, err := service.GetCommunityList()
	if err != nil {
		zap.L().Error("service.GetCommunityList() failed", zap.Error(err))
		utils.ResponseError(c, utils.CodeServerBusy) // 不轻易把服务端报错暴露给外面
		return
	}
	utils.ResponseSuccess(c, communityList)
}

// CommunityDetailHandler 社区详情
func CommunityDetailHandler(c *gin.Context) {
	// 1、获取社区ID
	communityIdStr := c.Param("id")                               // 获取URL参数
	communityId, err := strconv.ParseUint(communityIdStr, 10, 64) // id字符串格式转换
	if err != nil {
		utils.ResponseError(c, utils.CodeInvalidParams)
		return
	}

	// 2、根据ID获取社区详情
	communityList, err := service.GetCommunityDetailByID(communityId)
	if err != nil {
		zap.L().Error("service.GetCommunityByID() failed", zap.Error(err))
		utils.ResponseErrorWithMsg(c, utils.CodeSuccess, err.Error())
		return
	}
	utils.ResponseSuccess(c, communityList)
}
