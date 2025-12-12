package timeutil

import (
	"fmt"
	"time"
)

// Format formats time with a lightweight pattern compatible with TeaGo.
// Supported tokens: Y (year), m (month), d (day), H (hour), i (minute), s (second), w (weekday 0-6), W (ISO week), t (days in month).
func Format(pattern string, t ...time.Time) string {
	target := time.Now()
	if len(t) > 0 {
		target = t[0]
	}
	return formatPattern(pattern, target)
}

// FormatTime formats a unix timestamp with the same pattern rules.
func FormatTime(pattern string, timestamp int64) string {
	return formatPattern(pattern, time.Unix(timestamp, 0))
}

func formatPattern(pattern string, t time.Time) string {
	var out []rune
	_, week := t.ISOWeek()
	for _, ch := range pattern {
		switch ch {
		case 'Y':
			out = append(out, []rune(fmt.Sprintf("%04d", t.Year()))...)
		case 'm':
			out = append(out, []rune(fmt.Sprintf("%02d", int(t.Month())))...)
		case 'd':
			out = append(out, []rune(fmt.Sprintf("%02d", t.Day()))...)
		case 'H':
			out = append(out, []rune(fmt.Sprintf("%02d", t.Hour()))...)
		case 'i':
			out = append(out, []rune(fmt.Sprintf("%02d", t.Minute()))...)
		case 's':
			out = append(out, []rune(fmt.Sprintf("%02d", t.Second()))...)
		case 'w':
			out = append(out, []rune(fmt.Sprintf("%d", int(t.Weekday())))...)
		case 'W':
			out = append(out, []rune(fmt.Sprintf("%02d", week))...)
		case 't':
			// last day of month
			firstNextMonth := time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, t.Location())
			lastDay := firstNextMonth.AddDate(0, 0, -1).Day()
			out = append(out, []rune(fmt.Sprintf("%02d", lastDay))...)
		default:
			out = append(out, ch)
		}
	}
	return string(out)
}
