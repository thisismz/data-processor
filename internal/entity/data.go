package entity

import "github.com/google/uuid"

type Data struct {
	UID       uuid.UUID
	UserQuota string
	DataQuota string
	S3Path    string
	Size      int
}
