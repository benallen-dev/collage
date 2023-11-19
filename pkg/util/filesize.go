package util

import (
	"fmt"
)

func FormatBytes(size int64) string {
	units := []string{"B", "kB", "MB", "GB", "TB", "PB", "EB"}

	sizeFloat := float64(size)

	// I copied this straight from chatGPT and it looks kinda terrible
	// while also solving the problem of me being too lazy to do this
	// myself, so it's kind of a win I guess.
	unitIndex := 0
	for sizeFloat >= 1024 && unitIndex < len(units)-1 {
		sizeFloat /= 1024
		unitIndex++
	}

	formattedSize := fmt.Sprintf("%.2f %s", sizeFloat, units[unitIndex])
	return formattedSize
}
