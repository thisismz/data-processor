package service

import (
	"testing"
)

func TestSendToQueue(t *testing.T) {
	QueueServiceUp()
	err := SendToQueue()
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
