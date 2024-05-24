package sql

import (
	"context"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/thisismz/data-processor/internal/entity"
	"github.com/thisismz/data-processor/pkg/env"
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
	if err := r.db.Where("user_quota = ?", userQuota).Last(&user).Error; err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (r *SqlRepository) GetSync(ctx context.Context, date time.Time) ([]entity.User, error) {
	var users []entity.User
	if err := r.db.Where("date <= ? and is_sync =", date, false).Order("date asc").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
func (r *SqlRepository) GetData(ctx context.Context, dataQuota string) (entity.User, error) {
	var user entity.User
	if err := r.db.Where("data_quota = ?", dataQuota).Last(&user).Error; err != nil {
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
	if err := r.db.Where("id = ?", user.ID).Updates(entity.User{
		IsSync: true,
	}).Error; err != nil {
		return err
	}
	return nil
}

func (r *SqlRepository) CheckDuplicate(ctx context.Context, userQuota, dataQuota string) (bool, error) {
	userDataQuota := userQuota + ":" + dataQuota
	timeLimit, err := strconv.Atoi(env.GetEnv("DUPLICATE_TIME_LIMIT", "10"))
	if err != nil {
		log.Err(err).Msg("Error: in converting string to int")
	}
	now := time.Now()
	timeRange := now.Add(-time.Duration(timeLimit) * time.Minute)
	var user entity.User
	if err := r.db.Where("user_data_quota = ? AND created_at > ? ", userDataQuota, timeRange).First(&user).Error; err != nil {
		return false, err
	}
	if user.UID == uuid.Nil {
		return false, nil
	}
	return true, nil
}
func (r *SqlRepository) GetUnSync() ([]entity.User, error) {
	var user []entity.User
	if err := r.db.Where("is_sync = ?", false).First(&user).Error; err != nil {
		return []entity.User{}, err
	}
	return user, nil
}
