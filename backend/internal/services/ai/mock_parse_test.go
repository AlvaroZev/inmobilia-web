package ai

import "testing"

func TestParseBodyCount(t *testing.T) {	tests := []struct {
		description string
		want        int
	}{
		{"Closet con dos cuerpos", 2},
		{"Ropero de 3 cuerpos", 3},
		{"Armario tres cuerpos empotrado", 3},
		{"Ropero 4-cuerpos", 4},
		{"Closet sin mencionar cuerpos", 2},
		{"Ropero de un cuerpo", 1},
	}

	for _, tt := range tests {
		if got := parseBodyCount(tt.description); got != tt.want {
			t.Fatalf("parseBodyCount(%q) = %d, want %d", tt.description, got, tt.want)
		}
	}
}

func TestParseDeskDrawerIntent(t *testing.T) {
	tests := []struct {
		description string
		enabled     bool
		mode        string
		count       int
	}{
		{"Escritorio sin cajones", false, "", 0},
		{"Escritorio con un cajón lateral", true, "single", 1},
		{"Escritorio con cajonera", true, "tower", 3},
		{"Escritorio con 3 cajones", true, "tower", 3},
	}

	for _, tt := range tests {
		got := parseDeskDrawerIntent(tt.description)
		if got.enabled != tt.enabled || got.mode != tt.mode || got.count != tt.count {
			t.Fatalf("parseDeskDrawerIntent(%q) = %+v, want enabled=%v mode=%q count=%d",
				tt.description, got, tt.enabled, tt.mode, tt.count)
		}
	}
}