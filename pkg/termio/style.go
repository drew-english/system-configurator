package termio

import (
	"fmt"

	"github.com/mgutz/ansi"
)

var (
	magenta            = ansi.ColorFunc("magenta")
	cyan               = ansi.ColorFunc("cyan")
	red                = ansi.ColorFunc("red")
	yellow             = ansi.ColorFunc("yellow")
	blue               = ansi.ColorFunc("blue")
	green              = ansi.ColorFunc("green")
	gray               = ansi.ColorFunc("black+h")
	lightGrayUnderline = ansi.ColorFunc("white+du")
	bold               = ansi.ColorFunc("default+b")
	cyanBold           = ansi.ColorFunc("cyan+b")
	greenBold          = ansi.ColorFunc("green+b")

	gray256 = func(t string) string {
		return fmt.Sprintf("\x1b[%d;5;%dm%s\x1b[m", 38, 242, t)
	}
)

type style struct {
	enabled      bool
	is256enabled bool
	hasTrueColor bool
}

func NewStyle(enabled, is256enabled, trueColor bool) *style {
	return &style{
		enabled:      enabled,
		is256enabled: is256enabled,
		hasTrueColor: trueColor,
	}
}

func (s *style) Enabled() bool {
	return s.enabled
}

func (s *style) Fmt(style func(string) string, format string, a ...any) string {
	return style(fmt.Sprintf(format, a...))
}

func (s *style) Bold(t string) string {
	if !s.enabled {
		return t
	}

	return bold(t)
}

func (s *style) Red(t string) string {
	if !s.enabled {
		return t
	}

	return red(t)
}

func (s *style) Yellow(t string) string {
	if !s.enabled {
		return t
	}

	return yellow(t)
}

func (s *style) Green(t string) string {
	if !s.enabled {
		return t
	}

	return green(t)
}

func (s *style) GreenBold(t string) string {
	if !s.enabled {
		return t
	}

	return greenBold(t)
}

func (s *style) Gray(t string) string {
	if !s.enabled {
		return t
	}

	if s.is256enabled {
		return gray256(t)
	}

	return gray(t)
}

func (s *style) LightGrayUnderline(t string) string {
	if !s.enabled {
		return t
	}

	return lightGrayUnderline(t)
}

func (s *style) Magenta(t string) string {
	if !s.enabled {
		return t
	}

	return magenta(t)
}

func (s *style) Cyan(t string) string {
	if !s.enabled {
		return t
	}

	return cyan(t)
}

func (s *style) CyanBold(t string) string {
	if !s.enabled {
		return t
	}

	return cyanBold(t)
}

func (s *style) Blue(t string) string {
	if !s.enabled {
		return t
	}

	return blue(t)
}

func (s *style) SuccessIcon() string {
	return s.Green("✓")
}

func (s *style) WarningIcon() string {
	return s.Yellow("!")
}

func (s *style) FailureIcon() string {
	return s.Red("X")
}
