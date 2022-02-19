package printer

import "github.com/4khara/replica"

type MockPrinter struct{}

func (e *MockPrinter) Ln(a ...interface{}) {
	replica.MockFn(a...)
}

func beforePrinterTest() {
	WrapStdout(&MockPrinter{})
}
