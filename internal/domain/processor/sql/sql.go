package sql

import (
	"github.com/google/uuid"
	"github.com/thisismz/data-processor/internal/entity"
	"gorm.io/gorm"
)

type sqlRepository struct {
	db *gorm.DB
}

func NewSqlRepository(dataBaseConnections *gorm.DB) *sqlRepository {
	return &sqlRepository{
		db: dataBaseConnections,
	}
}

// read
func (r *sqlRepository) CheckDuplicate(userQuota string, dataQuota string) (bool, error) {
	var user entity.User
	if err := r.db.Where("user_quota = ? AND data_quota = ?", userQuota, dataQuota).First(&user).Error; err != nil {
		return false, err
	}
	if user.ID == uuid.Nil {
		return false, nil
	}
	return true, nil
}

func (r *sqlRepository) GetUser(userQuota string) (entity.User, error) {
	var user entity.User
	if err := r.db.Where("user_quota = ?", userQuota).First(&user).Error; err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (r *sqlRepository) GetData(dataQuota string) (entity.User, error) {
	var user entity.User
	if err := r.db.Where("data_quota = ?", dataQuota).First(&user).Error; err != nil {
		return entity.User{}, err
	}
	return user, nil
}

// write
func (r *sqlRepository) Add(user entity.User) error {
	if err := r.db.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func (r *sqlRepository) Update(user entity.User) error {
	if err := r.db.Save(&user).Error; err != nil {
		return err
	}
	return nil
}
