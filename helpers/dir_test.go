package helpers

import (
	"strings"
	"testing"
)

func TestCarbonComposeDirectory(t *testing.T) {
	dir := ComposeDir()

	// Make sure it ends with .carbon.
	// We only care that we store them in the right place.
	if !strings.HasSuffix(dir, ".carbon") {
		t.Errorf("Expected %s to end with .carbon", dir)
	}
}
