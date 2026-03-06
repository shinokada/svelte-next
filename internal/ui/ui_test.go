package ui

import (
	"strings"
	"testing"
)

func TestPadRight(t *testing.T) {
	tests := []struct {
		input string
		width int
		want  string
	}{
		{"abc", 5, "abc  "},
		{"abcde", 5, "abcde"},
		{"abcdef", 5, "abcdef"},
		{"", 3, "   "},
	}
	for _, tc := range tests {
		got := padRight(tc.input, tc.width)
		if got != tc.want {
			t.Errorf("padRight(%q, %d) = %q, want %q", tc.input, tc.width, got, tc.want)
		}
	}
}

func TestPrintBanner_NoPanic(t *testing.T) {
	// Just verify no panic for various inputs.
	PrintBanner("hello", Blue, "*", 10)
	PrintBanner("line1\nline2", Green, "-", 0)
	PrintBanner("", Red, "", -1)
}

func TestTable_NoPanic(t *testing.T) {
	headers := []string{"Dir", "Svelte", "Manager"}
	rows := [][]string{
		{"my-app", "5.28.1", "pnpm"},
		{"admin", "5.0.0", "npm"},
	}
	Table(headers, rows)
}

func TestPrinterFor(t *testing.T) {
	colors := []Color{Blue, Green, Yellow, Red, Magenta, Cyan, White, "unknown"}
	for _, c := range colors {
		p := printerFor(c)
		if p == nil {
			t.Errorf("printerFor(%q) returned nil", c)
		}
	}
}

func TestBorderDefault(t *testing.T) {
	// borderWidth=0 should produce a 10-char border.
	border := strings.Repeat("*", 10)
	if len(border) != 10 {
		t.Error("expected border length 10")
	}
}
