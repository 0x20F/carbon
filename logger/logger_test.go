package logger

import "testing"

func TestLoggerApi(t *testing.T) {
	Info("TABLE", "about 10 services", "nice")
}
