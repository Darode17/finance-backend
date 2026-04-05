package services

import (
	"errors"
	"math"
	"time"

	"github.com/radhikadarode/finance-backend/internal/database"
	"github.com/radhikadarode/finance-backend/internal/models"
	"gorm.io/gorm"
)

type RecordService struct{}

func NewRecordService() *RecordService {
	return &RecordService{}
}

func (s *RecordService) CreateRecord(req models.CreateRecordRequest, userID string) (*models.FinancialRecord, error) {
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, errors.New("invalid date format, use YYYY-MM-DD")
	}

	record := models.FinancialRecord{
		Amount:      req.Amount,
		Type:        req.Type,
		Category:    req.Category,
		Date:        date,
		Description: req.Description,
		CreatedBy:   userID,
	}

	if err := database.DB.Create(&record).Error; err != nil {
		return nil, err
	}
	return &record, nil
}

func (s *RecordService) GetRecords(filter models.RecordFilter) (*models.PaginatedRecords, error) {
	query := database.DB.Model(&models.FinancialRecord{}).Where("is_deleted = ?", false)

	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}
	if filter.Category != "" {
		query = query.Where("category = ?", filter.Category)
	}
	if filter.StartDate != "" {
		start, err := time.Parse("2006-01-02", filter.StartDate)
		if err == nil {
			query = query.Where("date >= ?", start)
		}
	}
	if filter.EndDate != "" {
		end, err := time.Parse("2006-01-02", filter.EndDate)
		if err == nil {
			query = query.Where("date <= ?", end)
		}
	}

	var total int64
	query.Count(&total)

	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PageSize < 1 || filter.PageSize > 100 {
		filter.PageSize = 20
	}

	offset := (filter.Page - 1) * filter.PageSize
	var records []models.FinancialRecord
	err := query.Preload("User").
		Order("date DESC").
		Offset(offset).
		Limit(filter.PageSize).
		Find(&records).Error

	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(filter.PageSize)))
	return &models.PaginatedRecords{
		Data:       records,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
		TotalCount: total,
		TotalPages: totalPages,
	}, nil
}

func (s *RecordService) GetRecordByID(id string) (*models.FinancialRecord, error) {
	var record models.FinancialRecord
	err := database.DB.Preload("User").
		Where("id = ? AND is_deleted = ?", id, false).
		First(&record).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("record not found")
	}
	return &record, err
}

func (s *RecordService) UpdateRecord(id string, req models.UpdateRecordRequest) (*models.FinancialRecord, error) {
	record, err := s.GetRecordByID(id)
	if err != nil {
		return nil, err
	}

	updates := map[string]interface{}{}
	if req.Amount != nil {
		updates["amount"] = *req.Amount
	}
	if req.Type != "" {
		updates["type"] = req.Type
	}
	if req.Category != "" {
		updates["category"] = req.Category
	}
	if req.Date != "" {
		d, err := time.Parse("2006-01-02", req.Date)
		if err != nil {
			return nil, errors.New("invalid date format, use YYYY-MM-DD")
		}
		updates["date"] = d
	}
	updates["description"] = req.Description

	if err := database.DB.Model(record).Updates(updates).Error; err != nil {
		return nil, err
	}
	return record, nil
}

func (s *RecordService) DeleteRecord(id string) error {
	// Soft delete
	result := database.DB.Model(&models.FinancialRecord{}).
		Where("id = ? AND is_deleted = ?", id, false).
		Update("is_deleted", true)
	if result.RowsAffected == 0 {
		return errors.New("record not found")
	}
	return result.Error
}
