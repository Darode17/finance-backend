package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Role constants
type Role string

const (
	RoleViewer  Role = "viewer"
	RoleAnalyst Role = "analyst"
	RoleAdmin   Role = "admin"
)

// TransactionType constants
type TransactionType string

const (
	TypeIncome  TransactionType = "income"
	TypeExpense TransactionType = "expense"
)

// User represents a system user
type User struct {
	ID        string         `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"not null"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	Password  string         `json:"-" gorm:"not null"`
	Role      Role           `json:"role" gorm:"not null;default:'viewer'"`
	IsActive  bool           `json:"is_active" gorm:"not null;default:true"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}

// FinancialRecord represents a financial transaction/entry
type FinancialRecord struct {
	ID          string          `json:"id" gorm:"primaryKey"`
	Amount      float64         `json:"amount" gorm:"not null"`
	Type        TransactionType `json:"type" gorm:"not null"`
	Category    string          `json:"category" gorm:"not null"`
	Date        time.Time       `json:"date" gorm:"not null"`
	Description string          `json:"description"`
	CreatedBy   string          `json:"created_by" gorm:"not null"`
	User        User            `json:"user,omitempty" gorm:"foreignKey:CreatedBy"`
	IsDeleted   bool            `json:"-" gorm:"default:false"` // soft delete flag
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

func (r *FinancialRecord) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.New().String()
	}
	return nil
}

// --- Request/Response DTOs ---

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required,min=2"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     Role   `json:"role" binding:"required,oneof=viewer analyst admin"`
}

type UpdateUserRequest struct {
	Name     string `json:"name" binding:"omitempty,min=2"`
	Role     Role   `json:"role" binding:"omitempty,oneof=viewer analyst admin"`
	IsActive *bool  `json:"is_active"`
}

type CreateRecordRequest struct {
	Amount      float64         `json:"amount" binding:"required,gt=0"`
	Type        TransactionType `json:"type" binding:"required,oneof=income expense"`
	Category    string          `json:"category" binding:"required,min=1"`
	Date        string          `json:"date" binding:"required"` // "YYYY-MM-DD"
	Description string          `json:"description"`
}

type UpdateRecordRequest struct {
	Amount      *float64        `json:"amount" binding:"omitempty,gt=0"`
	Type        TransactionType `json:"type" binding:"omitempty,oneof=income expense"`
	Category    string          `json:"category" binding:"omitempty,min=1"`
	Date        string          `json:"date"`
	Description string          `json:"description"`
}

type RecordFilter struct {
	Type      string `form:"type"`
	Category  string `form:"category"`
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
	Page      int    `form:"page,default=1"`
	PageSize  int    `form:"page_size,default=20"`
}

// Dashboard summary types

type CategoryTotal struct {
	Category string  `json:"category"`
	Total    float64 `json:"total"`
	Count    int64   `json:"count"`
}

type MonthlyTrend struct {
	Month   string  `json:"month"` // "YYYY-MM"
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
	Net     float64 `json:"net"`
}

type DashboardSummary struct {
	TotalIncome      float64         `json:"total_income"`
	TotalExpenses    float64         `json:"total_expenses"`
	NetBalance       float64         `json:"net_balance"`
	TotalRecords     int64           `json:"total_records"`
	CategoryTotals   []CategoryTotal `json:"category_totals"`
	MonthlyTrends    []MonthlyTrend  `json:"monthly_trends"`
	RecentActivity   []FinancialRecord `json:"recent_activity"`
}

type PaginatedRecords struct {
	Data       []FinancialRecord `json:"data"`
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
	TotalCount int64             `json:"total_count"`
	TotalPages int               `json:"total_pages"`
}

// APIResponse is a standard wrapper for all API responses
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}
