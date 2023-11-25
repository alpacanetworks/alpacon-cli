package utils

import (
	"github.com/gookit/color"
)

// Green converts a string to green color in the console.
func Green(value string) string {
	return color.FgGreen.Render(value)
}

// Yellow converts a string to yellow color in the console.
func Yellow(value string) string {
	return color.FgYellow.Render(value)
}

// Blue converts a string to blue color in the console.
func Blue(value string) string {
	return color.FgBlue.Render(value)
}

// Red converts a string to red color in the console.
func Red(value string) string {
	return color.FgRed.Render(value)
}
