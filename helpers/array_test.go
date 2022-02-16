package helpers

import "testing"

func TestContainsFindsItems(t *testing.T) {
	s := []string{"a", "b", "c"}

	if !Contains(s, "b") {
		t.Errorf("Expected %s to contain %s", s, "b")
	}
}
