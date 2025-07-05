package utils

import (
	"fmt"
	"strings"
	"time"
)

func SmoothProgress(start, end int, status string, duration time.Duration, keepVisible bool) {
	steps := 15 // Reduce steps for faster updates
	delay := duration / time.Duration(steps)
	spinChars := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	spinIdx := 0

	for i := 0; i <= steps; i++ {
		progress := start + (end-start)*i/steps

		// Update spinner index
		spinIdx = (spinIdx + 1) % len(spinChars)

		showProgress(progress, status, spinChars[spinIdx])
		time.Sleep(delay)
		if i < steps || !keepVisible {
			fmt.Print("\r")
		}
	}
}

func showProgress(percent int, status string, spinChar string) {
	width := 30

	// Calculate completed width
	completed := width * percent / 100

	// Clear the current line and move to start
	fmt.Print("\033[2K\r")

	// Print spinner and progress bar
	fmt.Printf("%s [", spinChar)
	fmt.Print("\033[36m") // Cyan color
	fmt.Print(strings.Repeat("█", completed))
	if completed < width {
		fmt.Print(strings.Repeat("░", width-completed))
	}
	fmt.Print("\033[0m") // Reset color
	fmt.Printf("] %3d%% %s", percent, status)
}
