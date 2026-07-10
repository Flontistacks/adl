package aria2

import (
	"fmt"
	"strings"
)

func FormatBytes(n int64) string {
	const unit = 1024
	if n < unit {
		return fmt.Sprintf("%d B", n)
	}
	div, exp := int64(unit), 0
	for v := n / unit; v >= unit; v /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(n)/float64(div), "KMGTPE"[exp])
}

func FormatSpeed(n int64) string {
	return FormatBytes(n) + "/s"
}

func FormatETA(sec int64) string {
	if sec <= 0 {
		return "--"
	}
	if sec < 60 {
		return fmt.Sprintf("%ds", sec)
	}
	if sec < 3600 {
		return fmt.Sprintf("%dm%ds", sec/60, sec%60)
	}
	return fmt.Sprintf("%dh%dm", sec/3600, (sec%3600)/60)
}

func ProgressBar(completed, total int64, width int) string {
	if width < 10 {
		width = 10
	}
	pct := 0.0
	if total > 0 {
		pct = float64(completed) / float64(total) * 100
	}
	filled := int(pct / 100 * float64(width))
	if filled > width {
		filled = width
	}
	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	return fmt.Sprintf("%s %5.1f%%", bar, pct)
}
