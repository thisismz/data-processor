package entity

import (
	"time"

	"github.com/goccy/go-json"

	"github.com/google/uuid"
)

type User struct {
	// ID is the identifier of the Entity, the ID is shared for all sub domains
	ID           uuid.UUID
	CreateAt     time.Time
	UserQuota    string
	DataQuota    string
	S3Path       string
	RateLImit    int
	TrafficLImit int
}

func (u User) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}
func (u *User) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}
