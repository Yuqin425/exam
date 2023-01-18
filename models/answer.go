package models

import "time"

type Answer struct {
	ProblemID  uint64    `json:"problem_id,string" db:"problem_id" binding:"required"`
	ParentID   uint64    `db:"parent_id" json:"parent_id"`
	AnswerID   uint64    `db:"answer_id" json:"answer_id"`
	AuthorID   uint64    `json:"author_id" db:"author_id"`
	Content    string    `json:"content" db:"content" binding:"required"`
	CreateTime time.Time `db:"create_time" json:"create_time"`
}
