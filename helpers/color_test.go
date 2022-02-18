package helpers

import "testing"

func TestStringToColorReturnsHexColor(t *testing.T) {
	color := StringToColor("test")

	if len(color) != 7 {
		t.Error("StringToColor should return a hex color in the #RRGGBB format")
	}
}

func TestStringToColorIsConsistent(t *testing.T) {
	color1 := StringToColor("test")
	color2 := StringToColor("test")

	if color1 != color2 {
		t.Error("StringToColor should return the same color for the same string")
	}
}

func TestStringToColorColorIsBasedOnInput(t *testing.T) {
	color1 := StringToColor("test")
	color2 := StringToColor("test2")

	if color1 == color2 {
		t.Error("StringToColor should return a different color for different strings")
	}
}
