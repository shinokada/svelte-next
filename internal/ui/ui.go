// Package ui provides colored output helpers and banner printing for
// svelte-next. All output goes to os.Stdout / os.Stderr. The fatih/color
// package automatically disables ANSI escape codes when the output is not a
// TTY (e.g. redirected to a file) and on Windows terminals that do not
// support VT sequences.
package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

// Pre-built color printers used across the package.
var (
	blue    = color.New(color.FgBlue, color.Bold)
	green   = color.New(color.FgGreen, color.Bold)
	yellow  = color.New(color.FgYellow, color.Bold)
	red     = color.New(color.FgRed, color.Bold)
	magenta = color.New(color.FgMagenta, color.Bold)
	cyan    = color.New(color.FgCyan, color.Bold)
	white   = color.New(color.FgWhite, color.Bold)
)

// Color names accepted by PrintBanner.
type Color string

const (
	Blue    Color = "blue"
	Green   Color = "green"
	Yellow  Color = "yellow"
	Red     Color = "red"
	Magenta Color = "magenta"
	Cyan    Color = "cyan"
	White   Color = "white"
)

func printerFor(c Color) *color.Color {
	switch c {
	case Blue:
		return blue
	case Green:
		return green
	case Yellow:
		return yellow
	case Red:
		return red
	case Magenta:
		return magenta
	case Cyan:
		return cyan
	default:
		return white
	}
}

// PrintBanner prints a bordered banner block, mirroring newBannerColor() from
// the Bash version. borderChar is typically "*". borderWidth controls the
// length of the top/bottom border line; 0 defaults to 10.
func PrintBanner(msg string, c Color, borderChar string, borderWidth int) {
	if borderWidth <= 0 {
		borderWidth = 10
	}
	if borderChar == "" {
		borderChar = "*"
	}
	border := strings.Repeat(borderChar, borderWidth)
	p := printerFor(c)
	_, _ = fmt.Fprintln(os.Stdout, p.Sprint(border))
	for _, line := range strings.Split(msg, "\n") {
		_, _ = fmt.Fprintln(os.Stdout, p.Sprint(line))
	}
	_, _ = fmt.Fprintln(os.Stdout, p.Sprint(border))
	_, _ = fmt.Fprintln(os.Stdout)
}

// DryRunf prints a dry-run preview line to stdout.
// All lines are prefixed with "[dry-run] " in cyan so they are visually
// distinct from real execution output.
func DryRunf(format string, a ...any) {
	_, _ = fmt.Fprintln(os.Stdout, cyan.Sprint(fmt.Sprintf("[dry-run] "+format, a...)))
}

// Infof prints an informational line to stdout (blue).
func Infof(format string, a ...any) {
	_, _ = fmt.Fprintln(os.Stdout, blue.Sprint(fmt.Sprintf(format, a...)))
}

// Successf prints a success line to stdout (green).
func Successf(format string, a ...any) {
	_, _ = fmt.Fprintln(os.Stdout, green.Sprint(fmt.Sprintf(format, a...)))
}

// Warnf prints a warning line to stdout (yellow).
func Warnf(format string, a ...any) {
	_, _ = fmt.Fprintln(os.Stdout, yellow.Sprint(fmt.Sprintf(format, a...)))
}

// Errorf prints an error line to stderr (red).
func Errorf(format string, a ...any) {
	_, _ = fmt.Fprintln(os.Stderr, red.Sprint(fmt.Sprintf(format, a...)))
}

// Table renders a simple fixed-width table to stdout.
// headers is the column header slice; rows is a slice of row slices.
// All columns are left-aligned and padded to the widest value.
func Table(headers []string, rows [][]string) {
	// Calculate column widths.
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(h)
	}
	for _, row := range rows {
		for i, cell := range row {
			if i < len(widths) && len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	// Build separator line.
	parts := make([]string, len(headers))
	for i, w := range widths {
		parts[i] = strings.Repeat("─", w)
	}
	separator := strings.Join(parts, "  ")

	// Print header.
	headerCells := make([]string, len(headers))
	for i, h := range headers {
		headerCells[i] = padRight(h, widths[i])
	}
	_, _ = fmt.Fprintln(os.Stdout, white.Sprint(strings.Join(headerCells, "  ")))
	_, _ = fmt.Fprintln(os.Stdout, separator)

	// Print rows.
	for _, row := range rows {
		cells := make([]string, len(headers))
		for i := range headers {
			val := ""
			if i < len(row) {
				val = row[i]
			}
			cells[i] = padRight(val, widths[i])
		}
		_, _ = fmt.Fprintln(os.Stdout, strings.Join(cells, "  "))
	}
	_, _ = fmt.Fprintln(os.Stdout)
}

func padRight(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
}
