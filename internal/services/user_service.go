package services

import (
	"errors"

	"github.com/radhikadarode/finance-backend/internal/database"
	"github.com/radhikadarode/finance-backend/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) CreateUser(req models.CreateUserRequest) (*models.User, error) {
	// Check duplicate email
	var existing models.User
	err := database.DB.Where("email = ?", req.Email).First(&existing).Error
	if err == nil {
		return nil, errors.New("a user with this email already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hash),
		Role:     req.Role,
		IsActive: true,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) GetAllUsers(page, pageSize int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	database.DB.Model(&models.User{}).Count(&total)

	offset := (page - 1) * pageSize
	err := database.DB.Offset(offset).Limit(pageSize).Find(&users).Error
	return users, total, err
}

func (s *UserService) GetUserByID(id string) (*models.User, error) {
	var user models.User
	err := database.DB.First(&user, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}
	return &user, err
}

func (s *UserService) UpdateUser(id string, req models.UpdateUserRequest) (*models.User, error) {
	user, err := s.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Role != "" {
		updates["role"] = req.Role
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if err := database.DB.Model(user).Updates(updates).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) DeleteUser(id string) error {
	result := database.DB.Delete(&models.User{}, "id = ?", id)
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return result.Error
}
