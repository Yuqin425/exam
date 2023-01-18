package mysql

import (
	"LanShan/models"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"strings"
)

// CreateProblem 发布题目
func CreateProblem(problem *models.Problem) (err error) {
	sqlStr := `insert into problem(
	problem_id, title, content, author_id, community_id)
	values(?,?,?,?,?)`
	_, err = db.Exec(sqlStr, problem.ProblemID, problem.Title,
		problem.Content, problem.AuthorId, problem.CommunityID)
	if err != nil {
		zap.L().Error("insert problem failed", zap.Error(err))
		err = ErrorInsertFailed
		return
	}
	return
}

func GetProblemByID(pid int64) (problem *models.Problem, err error) {
	problem = new(models.Problem)
	sqlStr := `select problem_id, title, content, author_id, community_id, create_time
	from problem
	where problem_id = ?`
	err = db.Get(problem, sqlStr, pid)
	if err == sql.ErrNoRows {
		err = ErrorInvalidID
		return
	}
	if err != nil {
		zap.L().Error("query problem failed", zap.String("sql", sqlStr), zap.Error(err))
		err = ErrorQueryFailed
		return
	}
	return
}

func GetProblemListByIDs(ids []string) (problemList []*models.Problem, err error) {
	sqlStr := `select problem_id, title, content, author_id, community_id, create_time
	from problem
	where problem_id in (?)
	order by FIND_IN_SET(problem_id, ?)`
	// 动态填充id
	query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
	if err != nil {
		return
	}
	// 使用Rebind()重新绑定
	query = db.Rebind(query)
	err = db.Select(&problemList, query, args...)
	return
}

func GetProblemList(page, size int64) (problems []*models.Problem, err error) {
	sqlStr := `select problem_id, title, content, author_id, community_id, create_time
	from problem
	ORDER BY create_time
	DESC 
	limit ?,?
	`
	problems = make([]*models.Problem, 0, 10)
	err = db.Select(&problems, sqlStr, (page-1)*size, size)
	return
}

func UpdateProblem(newProblem, pastProblem *models.Problem) (problem *models.Problem, err error) {
	if newProblem.Title != pastProblem.Title {
		sqlStr := "update problem set title = ? where problem_id = ?"
		_, err = db.Exec(sqlStr, newProblem.Title, newProblem.ProblemID)
		if err != nil {
			err = ErrorUpdateFailer
		}
	}
	if newProblem.Content != pastProblem.Content {
		sqlStr := "update problem set content = ? where problem_id = ?"
		_, err = db.Exec(sqlStr, newProblem.Content, newProblem.ProblemID)
		if err != nil {
			err = ErrorUpdateFailer
		}
	}
	if newProblem.Input != pastProblem.Input {
		sqlStr := "update problem set input = ? where problem_id = ?"
		_, err = db.Exec(sqlStr, newProblem.Title, newProblem.ProblemID)
		if err != nil {
			err = ErrorUpdateFailer
		}
	}
	if newProblem.Output != pastProblem.Output {
		sqlStr := "update problem set output = ? where problem_id = ?"
		_, err = db.Exec(sqlStr, newProblem.Title, newProblem.ProblemID)
		if err != nil {
			err = ErrorUpdateFailer
		}
	}
	return
}

func DeleteProblem(problemID int64) (err error) {
	sqlStr := "delete from problem WHERE problem_id = ?"
	_, err = db.Exec(sqlStr, problemID)
	if err != nil {
		zap.L().Error("delete problem failed", zap.Error(err))
		err = ErrorInsertFailed
		return
	}
	return
}
