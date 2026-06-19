package manufacturing

import (
	"testing"
)

func TestDrawerBodyHeights(t *testing.T) {
	const panel = 18.0

	t.Run("carcass envelope 350mm", func(t *testing.T) {
		const slot = 350.0
		available := drawerInteriorAvailableHeight(slot, panel, true)
		if available != 314 {
			t.Fatalf("available interior = %v, want 314", available)
		}
		side := drawerBodySideHeight(slot, panel, true)
		if side != 267 {
			t.Fatalf("body side height = %v, want 267", side)
		}
		falseFront := drawerFalseFrontHeight(side)
		if falseFront != 227 {
			t.Fatalf("false front height = %v, want 227", falseFront)
		}
	})

	t.Run("continuous drawer 175mm", func(t *testing.T) {
		const slot = 175.0
		available := drawerInteriorAvailableHeight(slot, panel, false)
		if available != 175 {
			t.Fatalf("available interior = %v, want 175", available)
		}
		side := drawerBodySideHeight(slot, panel, false)
		if side != 149 {
			t.Fatalf("body side height = %v, want 149", side)
		}
	})
}

func TestSnapRunnerLengthMm(t *testing.T) {
	tests := []struct {
		available float64
		want      float64
	}{
		{200, 250},
		{572, 550},
		{600, 600},
		{649, 600},
		{650, 650},
		{700, 650},
	}
	for _, tc := range tests {
		got := snapRunnerLengthMm(tc.available)
		if got != tc.want {
			t.Fatalf("snapRunnerLengthMm(%v) = %v, want %v", tc.available, got, tc.want)
		}
	}
}
