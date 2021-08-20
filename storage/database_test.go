package storage

import (
	"testing"
)

func TestUpsMoDbConn(t *testing.T) {
	_, err := UpsMoDBConn()
	if err != nil {
		t.Errorf(err.Error())
	}
}
