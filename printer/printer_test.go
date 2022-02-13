package printer

import "testing"

func TestLoggerApi(t *testing.T) {
	Info(Cyan, "TABLE", "about 10 services", "nice")
}
