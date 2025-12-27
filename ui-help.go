package main

import "fmt"

func formatDuration(seconds int) string {
	minutes := seconds / 60
	secs := seconds % 60
	return fmt.Sprintf("%d:%02d", minutes, secs)
}
