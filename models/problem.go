package models

import (
	"encoding/json"
	"errors"
	"time"
)

// 内存对齐概念 字段类型相同的对齐 缩小变量所占内存大小
type Problem struct {
	ProblemID   uint64    `json:"problem_id,string" db:"problem_id"`
	AuthorId    uint64    `json:"author_id" db:"author_id"`
	CommunityID uint64    `json:"community_id" db:"community_id" binding:"required"`
	Status      int32     `json:"status" db:"status"`
	Title       string    `json:"title" db:"title" binding:"required"`
	Content     string    `json:"content" db:"content" binding:"required"`
	Input       string    `json:"input" db:"input" binding:"required"`
	Output      string    `json:"output" db:"output" binding:"required"`
	CreateTime  time.Time `json:"-" db:"create_time"`
}

// UnmarshalJSON 为Post类型实现自定义的UnmarshalJSON方法
func (p *Problem) UnmarshalJSON(data []byte) (err error) {
	required := struct {
		Title       string `json:"title" db:"title"`
		Content     string `json:"content" db:"content"`
		CommunityID int64  `json:"community_id" db:"community_id"`
		InPut       string `json:"input" db:"input"`
		Output      string `json:"output" db:"output"`
	}{}
	err = json.Unmarshal(data, &required)
	if err != nil {
		return
	} else if len(required.Title) == 0 {
		err = errors.New("标题不能为空")
	} else if len(required.Content) == 0 {
		err = errors.New("内容不能为空")
	} else if required.CommunityID == 0 {
		err = errors.New("未指定版块")
	} else if required.CommunityID == 0 {
		err = errors.New("未指定输入")
	} else if required.CommunityID == 0 {
		err = errors.New("未指定输出")
	} else {
		p.Title = required.Title
		p.Content = required.Content
		p.CommunityID = uint64(required.CommunityID)
		p.Input = required.InPut
		p.Output = required.Output
	}
	return
}

type ApiProblemDetail struct {
	*Problem                            // 嵌入问题结构体
	*CommunityDetail `json:"community"` // 嵌入社区信息
	AuthorName       string             `json:"author_name"`
	VoteNum          int64              `json:"vote_num"`
	//CommunityName string `json:"community_name"`
}
