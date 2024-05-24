package valueobject

import (
	"github.com/goccy/go-json"

	"github.com/google/uuid"
)

type Data struct {
	UID       uuid.UUID
	UserQuota string
	DataQuota string
	S3Path    string
	Size      int64
	Hash      string
}

func (d *Data) MarshalBinary() ([]byte, error) {
	return json.Marshal(d)
}

func (d *Data) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, d)
}
