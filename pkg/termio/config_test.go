package termio_test

import (
	"os"

	"github.com/drew-english/system-configurator/pkg/termio"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
	Describe("ConfigFromEnv", func() {
		var (
			noColorEnvVar string
			termEnvVar    string
			colorTermVar  string
		)

		BeforeEach(func() {
			noColorEnvVar = ""
			termEnvVar = "xterm-256color"
			colorTermVar = "truecolor"
		})

		JustBeforeEach(func() {
			os.Setenv("NO_COLOR", noColorEnvVar)
			os.Setenv("TERM", termEnvVar)
			os.Setenv("COLORTERM", colorTermVar)
		})

		AfterEach(func() {
			os.Unsetenv("NO_COLOR")
			os.Unsetenv("TERM")
			os.Unsetenv("COLORTERM")
		})

		It("returns a Config instance", func() {
			var cfg *termio.Config = termio.ConfigFromEnv()
			Expect(cfg).ToNot(BeNil())
			Expect(cfg).To(BeComparableTo(&termio.Config{
				ColorDisabled:    false,
				Color256Enabled:  true,
				TrueColorEnabled: true,
			}))
		})

		Context("when NO_COLOR is set", func() {
			BeforeEach(func() {
				noColorEnvVar = "1"
			})

			It("returns a Config instance with ColorDisabled set to true", func() {
				var cfg *termio.Config = termio.ConfigFromEnv()
				Expect(cfg).ToNot(BeNil())
				Expect(cfg).To(BeComparableTo(&termio.Config{
					ColorDisabled:    true,
					Color256Enabled:  true,
					TrueColorEnabled: true,
				}))
			})
		})

		Context("when TERM is not set", func() {
			BeforeEach(func() {
				termEnvVar = ""
			})

			It("returns a Config instance with Color256Enabled and TrueColorEnabled set to true", func() {
				var cfg *termio.Config = termio.ConfigFromEnv()
				Expect(cfg).ToNot(BeNil())
				Expect(cfg).To(BeComparableTo(&termio.Config{
					ColorDisabled:    false,
					Color256Enabled:  true,
					TrueColorEnabled: true,
				}))
			})
		})

		Context("when COLORTERM is not set", func() {
			BeforeEach(func() {
				colorTermVar = ""
			})

			It("returns a Config instance with Color256Enabled true and TrueColorEnabled set to false", func() {
				var cfg *termio.Config = termio.ConfigFromEnv()
				Expect(cfg).ToNot(BeNil())
				Expect(cfg).To(BeComparableTo(&termio.Config{
					ColorDisabled:    false,
					Color256Enabled:  true,
					TrueColorEnabled: false,
				}))
			})
		})

		Context("when TERM and COLORTERM are not set", func() {
			BeforeEach(func() {
				termEnvVar = ""
				colorTermVar = ""
			})

			It("returns a Config instance with Color256Enabled and TrueColorEnabled set to false", func() {
				var cfg *termio.Config = termio.ConfigFromEnv()
				Expect(cfg).ToNot(BeNil())
				Expect(cfg).To(BeComparableTo(&termio.Config{
					ColorDisabled:    false,
					Color256Enabled:  false,
					TrueColorEnabled: false,
				}))
			})
		})
	})
})
