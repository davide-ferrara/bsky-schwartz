package views

import "strconv"

func formatFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', 2, 64)
}
