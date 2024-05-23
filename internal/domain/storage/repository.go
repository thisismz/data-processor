package storage

import (
	"context"
	"errors"

	"github.com/thisismz/data-processor/internal/entity"
)

var (
	ErrUserNotFound = errors.New("the User was not found in the repository")
	// ErrDataNotFound is returned when a Data is not found.
	ErrDataNotFound = errors.New("the Data was not found in the repository")
	// ErrFailedToAddData is returned when the Data could not be added to the repository.
	ErrFailedToAddData = errors.New("failed to add the Data to the repository")
	// ErrUpdateData is returned when the Data could not be updated in the repository.
	ErrUpdateData = errors.New("failed to update the Data in the repository")
)

type StorageRepository interface {
	Add(ctx context.Context, user entity.User) error
	GetUser(ctx context.Context, userQuota string) (entity.User, error)
	GetData(ctx context.Context, dataQuota string) (entity.User, error)
	Update(ctx context.Context, user entity.User) error
	CheckDuplicate(ctx context.Context, userQuota string, dataQuota string) (bool, error)
}
