package sql

import (
	"context"

	"github.com/google/uuid"
	"github.com/thisismz/data-processor/internal/entity"
	"gorm.io/gorm"
)

type SqlRepository struct {
	db *gorm.DB
}

func NewSqlRepository(dataBaseConnections *gorm.DB) *SqlRepository {
	return &SqlRepository{
		db: dataBaseConnections,
	}
}

// read
func (r *SqlRepository) GetUser(ctx context.Context, userQuota string) (entity.User, error) {
	var user entity.User
	if err := r.db.Where("user_quota = ?", userQuota).First(&user).Error; err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (r *SqlRepository) GetData(ctx context.Context, dataQuota string) (entity.User, error) {
	var user entity.User
	if err := r.db.Where("data_quota = ?", dataQuota).First(&user).Error; err != nil {
		return entity.User{}, err
	}
	return user, nil
}

// write
func (r *SqlRepository) Create(ctx context.Context, user entity.User) error {
	if err := r.db.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func (r *SqlRepository) Update(ctx context.Context, user entity.User) error {
	if err := r.db.Save(&user).Error; err != nil {
		return err
	}
	return nil
}

func (r *SqlRepository) CheckDuplicate(ctx context.Context, userQuota string, dataQuota string) (bool, error) {
	var user entity.User
	if err := r.db.Where("user_quota = ? AND data_quota = ?", userQuota, dataQuota).First(&user).Error; err != nil {
		return false, err
	}
	if user.UID == uuid.Nil {
		return false, nil
	}
	return true, nil
}
