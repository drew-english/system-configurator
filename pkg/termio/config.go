package termio

import (
	"os"
	"strings"
)

type Config struct {
	ColorDisabled    bool
	Color256Enabled  bool
	TrueColorEnabled bool
}

func ConfigFromEnv() *Config {
	return &Config{
		ColorDisabled:    os.Getenv("NO_COLOR") != "",
		Color256Enabled:  color256Supported(),
		TrueColorEnabled: trueColorSupported(),
	}
}

func color256Supported() bool {
	return trueColorSupported() ||
		strings.Contains(os.Getenv("TERM"), "256") ||
		strings.Contains(os.Getenv("COLORTERM"), "256")
}

func trueColorSupported() bool {
	term := os.Getenv("TERM")
	colorterm := os.Getenv("COLORTERM")

	return strings.Contains(term, "24bit") ||
		strings.Contains(term, "truecolor") ||
		strings.Contains(colorterm, "24bit") ||
		strings.Contains(colorterm, "truecolor")
}
