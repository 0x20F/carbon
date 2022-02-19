package printer

import (
	"testing"
)

func TestLoggerApi(t *testing.T) {
	beforePrinterTest()

	Info(Cyan, "TABLE", "about 10 services", "nice")
}
