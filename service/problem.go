package service

import (
	"LanShan/dao/mysql"
	"LanShan/models"
	"LanShan/utils/snowflake"
	"fmt"
	"go.uber.org/zap"
)

func GetCommunityList() ([]*models.Community, error) {
	// 查数据库 查找到所有的community 并返回
	return mysql.GetCommunityList()
}

func GetCommunityDetailByID(id uint64) (*models.CommunityDetail, error) {
	return mysql.GetCommunityByID(id)
}

func CreateProblem(problem *models.Problem) (err error) {
	// 1、 生成ID
	problemID, err := snowflake.GetID()
	if err != nil {
		zap.L().Error("snowflake.GetID() failed", zap.Error(err))
		return
	}
	problem.ProblemID = problemID
	// 2、创建问题 保存到数据库
	if err := mysql.CreateProblem(problem); err != nil {
		zap.L().Error("mysql.CreateProblem(&problem) failed", zap.Error(err))
		return err
	}

	return
}

func GetProblemById(problemID int64) (data *models.ApiProblemDetail, err error) {
	// 查询信息
	problem, err := mysql.GetProblemByID(problemID)
	if err != nil {
		zap.L().Error("mysql.GetProblemByID(problemID) failed",
			zap.Int64("problemID", problemID),
			zap.Error(err))
		return nil, err
	}
	// 根据作者id查询作者信息
	user, err := mysql.GetUserByID(problem.AuthorId)
	if err != nil {
		zap.L().Error("mysql.GetUserByID() failed",
			zap.Uint64("AuthorID", problem.AuthorId),
			zap.Error(err))
		return
	}
	// 根据社区id查询社区详细信息
	community, err := mysql.GetCommunityByID(problem.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityByID() failed",
			zap.Uint64("community_id", problem.CommunityID),
			zap.Error(err))
		return
	}
	// 接口数据拼接
	data = &models.ApiProblemDetail{
		Problem:         problem,
		CommunityDetail: community,
		AuthorName:      user.UserName,
	}
	return
}

func GetProblemList(page, size int64) (data []*models.ApiProblemDetail, err error) {
	problemList, err := mysql.GetProblemList(page, size)
	if err != nil {
		fmt.Println(err)
		return
	}
	data = make([]*models.ApiProblemDetail, 0, len(problemList)) // data 初始化
	for _, problem := range problemList {
		// 根据作者id查询作者信息
		user, err := mysql.GetUserByID(problem.AuthorId)
		if err != nil {
			zap.L().Error("mysql.GetUserByID() failed",
				zap.Uint64("problemID", problem.AuthorId),
				zap.Error(err))
			continue
		}
		// 根据社区id查询社区详细信息
		community, err := mysql.GetCommunityByID(problem.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityByID() failed",
				zap.Uint64("community_id", problem.CommunityID),
				zap.Error(err))
			continue
		}
		// 接口数据拼接
		problemdetail := &models.ApiProblemDetail{
			Problem:         problem,
			CommunityDetail: community,
			AuthorName:      user.UserName,
		}
		data = append(data, problemdetail)
	}
	return
}

func UpdateProblem(newProblem *models.Problem, pastProblemID int64) (data *models.ApiProblemDetail, err error) {
	// 查询信息
	pastProblem, err := mysql.GetProblemByID(pastProblemID)
	if err != nil {
		zap.L().Error("mysql.GetProblemByID(pastProblemID) failed",
			zap.Int64("pastProblemID", pastProblemID),
			zap.Error(err))
		return nil, err
	}
	problem, err := mysql.UpdateProblem(newProblem, pastProblem)
	if err != nil {
		zap.L().Error("mysql.UpdateProblem() failed", zap.Error(err))
		return nil, err
	}

	data, err = GetProblemById(int64(problem.ProblemID))
	if err != nil {
		zap.L().Error("GetProblemById failed", zap.Error(err))
		return nil, err
	}
	return
}

func DeleteProblem(problemID int64) (err error) {
	mysql.DeleteProblem(problemID)
	if err != nil {
		zap.L().Error("mysql.DeleteProblem() failed", zap.Error(err))
		return
	}
	return
}
