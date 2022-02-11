package helpers

import (
	"fmt"
	"hash/fnv"
	"io"
	"testing"
)

func fnvForTest(what string) string {
	h := fnv.New32a()
	io.WriteString(h, what)

	return fmt.Sprintf("%x", h.Sum32())
}

func TestHashSpecificLength(t *testing.T) {
	expected := fnvForTest("testing-hash")[:4]
	actual := Hash("testing-hash", 4)

	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestHashZeroPad(t *testing.T) {
	expected := "00" + fnvForTest("a")[:8]
	actual := Hash("a", 10)

	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}
