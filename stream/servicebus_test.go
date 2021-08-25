package stream

import (
	"testing"
)

func TestSendToServiceBus(t *testing.T) {
	//be sure that there is topic with `test` name in the service bus before running the test
	err := SendToServiceBus("test", "test message")
	if err != nil {
		t.Errorf(err.Error())
	}
}
