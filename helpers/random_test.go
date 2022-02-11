package helpers

import "testing"

func TestRandomStringLength(t *testing.T) {
	str := RandomAlphaString(50)

	if len(str) != 50 {
		t.Errorf("Expected length of %d, got %d", 50, len(str))
	}
}

func TestUniqueness(t *testing.T) {
	str1 := RandomAlphaString(50)
	str2 := RandomAlphaString(50)

	if str1 == str2 {
		t.Errorf("Expected %s and %s to be different", str1, str2)
	}
}
