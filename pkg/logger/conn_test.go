package logger

import (
	"testing"
)

func TestConn(t *testing.T) {
	log := NewLogger()
	_ = log.SetLogger("conn", `{"net":"tcp","addr":"10.1.55.10:1024"}`)
	log.Info("this is informational to net")
}
