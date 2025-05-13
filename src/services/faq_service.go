package services

import (
	"time"

	"github.com/charis16/luminor-golang-be/src/config"
	"github.com/charis16/luminor-golang-be/src/dto"
	"github.com/charis16/luminor-golang-be/src/models"
)

type FaqInput struct {
	QuestionID  string `json:"question_id" validate:"required"`
	QuestionEn  string `json:"question_en" validate:"required"`
	AnswerID    string `json:"answer_id" validate:"required"`
	AnswerEn    string `json:"answer_en" validate:"required"`
	IsPublished bool   `json:"is_published" validate:"required"`
}

func GetAllFaqs(page int, limit int, search string) ([]dto.FaqResponse, int64, error) {
	var faqs []models.Faq
	var total int64

	query := config.DB.Model(&models.Faq{})

	// Apply search filter if search term is provided
	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("question_id LIKE ? or question_en LIKE ? or answer_id LIKE ? or answer_en LIKE ?", searchTerm)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.
		Select("uuid", "question_id", "question_en", "answer_id", "answer_en", "is_published", "created_at", "updated_at").
		Limit(limit).
		Offset(offset).
		Find(&faqs).Error; err != nil {
		return nil, 0, err
	}

	// Mapping ke response DTO
	response := make([]dto.FaqResponse, len(faqs))
	for i, faq := range faqs {
		response[i] = dto.FaqResponse{
			UUID:        faq.UUID,
			AnswerID:    faq.AnswerID,
			AnswerEn:    faq.AnswerEn,
			QuestionID:  faq.QuestionID,
			QuestionEn:  faq.QuestionEn,
			IsPublished: faq.IsPublished,
			CreatedAt:   faq.CreatedAt,
			UpdatedAt:   faq.UpdatedAt,
		}
	}

	return response, total, nil
}

func CreateFaq(input FaqInput) (*models.Faq, error) {
	tx := config.DB.Begin() // Mulai transaksi

	faq := models.Faq{
		AnswerEn:    input.AnswerEn,
		AnswerID:    input.AnswerID,
		QuestionEn:  input.QuestionEn,
		QuestionID:  input.QuestionID,
		IsPublished: input.IsPublished,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := tx.Create(&faq).Error; err != nil {
		tx.Rollback() // rollback jika error
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err // commit gagal
	}

	return &faq, nil
}

func GetFaqByUUID(uuid string) (models.Faq, error) {
	var faq models.Faq
	if err := config.DB.Where("uuid = ?", uuid).First(&faq).Error; err != nil {
		return models.Faq{}, err
	}
	return faq, nil
}

func UpdateFaq(uuid string, input FaqInput) (models.Faq, error) {
	tx := config.DB.Begin()
	var faq models.Faq

	if err := tx.Where("uuid = ?", uuid).First(&faq).Error; err != nil {
		tx.Rollback()
		return models.Faq{}, err
	}

	faq.AnswerEn = input.AnswerEn
	faq.AnswerID = input.AnswerID
	faq.QuestionEn = input.QuestionEn
	faq.QuestionID = input.QuestionID
	faq.IsPublished = input.IsPublished
	faq.UpdatedAt = time.Now()

	if err := tx.Save(&faq).Error; err != nil {
		tx.Rollback()
		return models.Faq{}, err
	}

	if err := tx.Commit().Error; err != nil {
		return models.Faq{}, err
	}

	return faq, nil
}

func DeleteFaq(uuid string) error {
	if err := config.DB.Where("uuid = ?", uuid).Delete(&models.Faq{}).Error; err != nil {
		return err
	}
	return nil
}
