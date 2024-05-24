package entity

import (
	"time"

	"github.com/goccy/go-json"

	"github.com/google/uuid"
)

type User struct {
	ID                     uint `gorm:"primarykey"`
	UID                    uuid.UUID
	CreateAt               time.Time
	UserQuota              string
	DataQuota              string
	S3Path                 string
	RateLimit              int
	RateLimitExpiration    time.Time
	TrafficLimit           int64
	TrafficLimitExpiration time.Time
}

func (u User) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}
func (u *User) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}
