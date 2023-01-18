package mysql

import (
	"LanShan/models"
	"database/sql"
	"go.uber.org/zap"
)

func CreateAnswer(answer *models.Answer) (err error) {
	sqlStr := `insert into answer(
	answer_id, content, problem_id, author_id, parent_id)
	values(?,?,?,?,?)`
	_, err = db.Exec(sqlStr, answer.AnswerID, answer.Content, answer.ProblemID,
		answer.AuthorID, answer.ParentID)
	if err != nil {
		zap.L().Error("insert answer failed", zap.Error(err))
		err = ErrorInsertFailed
		return
	}
	return
}

func GetAnswerList(page, size int64, problemID int64) (answers []*models.Answer, err error) {
	sqlStr := `select answer_id, content, problem_id, author_id, parent_id, create_time
	from answer
	where problem_id = ?
	ORDER BY create_time
	DESC 
	limit ?,?
	`
	answers = make([]*models.Answer, 0, 10) // 0：长度  2：容量
	err = db.Select(&answers, sqlStr, problemID, (page-1)*size, size)
	return
}

func GetAnswerById(answerID int64) (answer *models.Answer, err error) {
	answer = new(models.Answer)
	sqlStr := `select answer_id, content, post_id, author_id, parent_id, create_time
	from answer
	where comment_id = ?`
	err = db.Get(answer, sqlStr, answerID)
	if err == sql.ErrNoRows {
		err = ErrorInvalidID
		return
	}
	if err != nil {
		zap.L().Error("query answer failed", zap.String("sql", sqlStr), zap.Error(err))
		err = ErrorQueryFailed
		return
	}
	return
}

func UpdateAnswer(answer *models.Answer) (data *models.Answer, err error) {
	sqlStr := `update answer set content = ? where answer_id = ?`
	_, err = db.Exec(sqlStr, answer.Content, answer.AnswerID)
	if err != nil {
		err = ErrorUpdateFailer
	}
	data, err = GetAnswerById(int64(answer.AnswerID))
	return
}

func DeleteAnswer(answerId int64) (err error) {
	sqlStr := `delete from answer where answer_id = ?`
	_, err = db.Exec(sqlStr, answerId)
	if err != nil {
		zap.L().Error("delete problem failed", zap.Error(err))
		err = ErrorInsertFailed
		return
	}
	return
}
