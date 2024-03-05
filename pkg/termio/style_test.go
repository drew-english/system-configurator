package termio_test

import (
	"fmt"

	"github.com/drew-english/system-configurator/pkg/termio"
	"github.com/mgutz/ansi"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	styles = map[string]func(string) string{
		"magenta":            ansi.ColorFunc("magenta"),
		"cyan":               ansi.ColorFunc("cyan"),
		"red":                ansi.ColorFunc("red"),
		"yellow":             ansi.ColorFunc("yellow"),
		"blue":               ansi.ColorFunc("blue"),
		"green":              ansi.ColorFunc("green"),
		"gray":               ansi.ColorFunc("black+h"),
		"lightGrayUnderline": ansi.ColorFunc("white+du"),
		"bold":               ansi.ColorFunc("default+b"),
		"cyanBold":           ansi.ColorFunc("cyan+b"),
		"greenBold":          ansi.ColorFunc("green+b"),
		"gray256": func(t string) string {
			return fmt.Sprintf("\x1b[%d;5;%dm%s\x1b[m", 38, 242, t)
		},
	}
)

var _ = Describe("Style", func() {
	var (
		enabled, is256enabled, hasTrueColor bool

		style = termio.NewStyle(enabled, is256enabled, hasTrueColor)
	)

	BeforeEach(func() {
		enabled = true
		is256enabled = true
		hasTrueColor = true
	})

	JustBeforeEach(func() {
		style = termio.NewStyle(enabled, is256enabled, hasTrueColor)
	})

	Describe("NewStyle", func() {
		It("returns a Style instance", func() {
			style := termio.NewStyle(true, true, true)
			Expect(style).ToNot(BeNil())
		})
	})

	Describe("Enabled", func() {
		It("returns true", func() {
			Expect(style.Enabled()).To(BeTrue())
		})

		Context("when enabled is false", func() {
			BeforeEach(func() {
				enabled = false
			})

			It("returns false", func() {
				Expect(style.Enabled()).To(BeFalse())
			})
		})
	})

	itReturnsStyledStringWhenEnabled := func(method func(string) string, styleName string) {
		It("returns colored string", func() {
			Expect(method("test")).To(Equal(styles[styleName]("test")))
		})

		Context("when enabled is false", func() {
			BeforeEach(func() {
				enabled = false
			})

			It("returns uncolored string", func() {
				Expect(method("test")).To(Equal("test"))
			})
		})
	}

	Describe("Bold", func() {
		itReturnsStyledStringWhenEnabled(func(s string) string {
			return style.Bold(s)
		}, "bold")
	})

	Describe("Red", func() {
		itReturnsStyledStringWhenEnabled(func(s string) string {
			return style.Red(s)
		}, "red")
	})

	Describe("Yellow", func() {
		itReturnsStyledStringWhenEnabled(func(s string) string {
			return style.Yellow(s)
		}, "yellow")
	})

	Describe("Green", func() {
		itReturnsStyledStringWhenEnabled(func(s string) string {
			return style.Green(s)
		}, "green")
	})

	Describe("GreenBold", func() {
		itReturnsStyledStringWhenEnabled(func(s string) string {
			return style.GreenBold(s)
		}, "greenBold")
	})

	Describe("Gray256", func() {
		itReturnsStyledStringWhenEnabled(func(s string) string {
			return style.Gray(s)
		}, "gray256")
	})

	Describe("Gray", func() {
		BeforeEach(func() {
			hasTrueColor = false
			is256enabled = false
		})

		itReturnsStyledStringWhenEnabled(func(s string) string {
			return style.Gray(s)
		}, "gray")
	})

	Describe("LightGrayUnderline", func() {
		itReturnsStyledStringWhenEnabled(func(s string) string {
			return style.LightGrayUnderline(s)
		}, "lightGrayUnderline")
	})

	Describe("Magenta", func() {
		itReturnsStyledStringWhenEnabled(func(s string) string {
			return style.Magenta(s)
		}, "magenta")
	})

	Describe("Cyan", func() {
		itReturnsStyledStringWhenEnabled(func(s string) string {
			return style.Cyan(s)
		}, "cyan")
	})

	Describe("CyanBold", func() {
		itReturnsStyledStringWhenEnabled(func(s string) string {
			return style.CyanBold(s)
		}, "cyanBold")
	})

	Describe("Blue", func() {
		itReturnsStyledStringWhenEnabled(func(s string) string {
			return style.Blue(s)
		}, "blue")
	})

	Describe("Fmt", func() {
		It("returns colored string", func() {
			Expect(style.Fmt(style.Bold, "test %s", "test")).To(Equal(styles["bold"]("test test")))
		})

		Context("when enabled is false", func() {
			BeforeEach(func() {
				enabled = false
			})

			It("returns uncolored string", func() {
				Expect(style.Fmt(style.Bold, "test %s", "test")).To(Equal("test test"))
			})
		})
	})

	Describe("SuccessIcon", func() {
		It("returns colored string", func() {
			Expect(style.SuccessIcon()).To(Equal(styles["green"]("✓")))
		})

		Context("when enabled is false", func() {
			BeforeEach(func() {
				enabled = false
			})

			It("returns uncolored string", func() {
				Expect(style.SuccessIcon()).To(Equal("✓"))
			})
		})
	})

	Describe("WarningIcon", func() {
		It("returns colored string", func() {
			Expect(style.WarningIcon()).To(Equal(styles["yellow"]("!")))
		})

		Context("when enabled is false", func() {
			BeforeEach(func() {
				enabled = false
			})

			It("returns uncolored string", func() {
				Expect(style.WarningIcon()).To(Equal("!"))
			})
		})
	})

	Describe("FailureIcon", func() {
		It("returns colored string", func() {
			Expect(style.FailureIcon()).To(Equal(styles["red"]("X")))
		})

		Context("when enabled is false", func() {
			BeforeEach(func() {
				enabled = false
			})

			It("returns uncolored string", func() {
				Expect(style.FailureIcon()).To(Equal("X"))
			})
		})
	})
})
