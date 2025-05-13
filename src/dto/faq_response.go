package dto

import "time"

type FaqResponse struct {
	UUID        string    `json:"uuid"`
	QuestionID  string    `json:"question_id"`
	QuestionEn  string    `json:"question_en"`
	AnswerID    string    `json:"answer_id"`
	AnswerEn    string    `json:"answer_en"`
	IsPublished bool      `json:"is_published"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
