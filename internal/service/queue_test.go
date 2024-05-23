package service

import (
	"testing"

	"github.com/google/uuid"
	"github.com/thisismz/data-processor/internal/entity"
)

func TestSendToQueue(t *testing.T) {
	QueueServiceUp()
	data := entity.Data{
		UID:       uuid.New(),
		DataQuota: uuid.New().String(),
		UserQuota: uuid.New().String(),
		S3Path:    "s3://bucket/path",
	}
	err := SendToQueue(data)
	if err != nil {
		t.Errorf("send to queue failed")
	}
}

func TestReceiveFromQueue(t *testing.T) {
	QueueServiceUp()
	done := make(chan bool)
	go ReceiveFromQueue(done)
	// <-done
	done <- true
}
