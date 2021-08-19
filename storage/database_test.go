package storage

import (
	"testing"
)

func TestUpMoDbConn(t *testing.T) {
	_, err := UpMoDBConn()
	if err != nil {
		t.Errorf(err.Error())
	}
}
