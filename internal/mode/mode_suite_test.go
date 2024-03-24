package mode_test

import (
	"fmt"
	"testing"

	"github.com/drew-english/system-configurator/internal/mode"
	"github.com/drew-english/system-configurator/spec/stub/termio"
	"github.com/spf13/viper"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestMode(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mode Suite")
}

var _ = Describe("Mode", func() {
	var recognizedModes = map[string]mode.Mode{
		"conf":          mode.ModeConfiguration,
		"configuration": mode.ModeConfiguration,
		"system":        mode.ModeSystem,
		"sys":           mode.ModeSystem,
		"hybrid":        mode.ModeHybrid,
		"hyb":           mode.ModeHybrid,
	}

	Describe("Parse", func() {
		for modeStr, modeInt := range recognizedModes {
			Context(fmt.Sprintf("when the mode is %s", modeStr), func() {
				It("returns the corresponding mode", func() {
					Expect(mode.Parse(modeStr)).To(Equal(modeInt))
				})
			})
		}

		Context("when the mode is not recognized", func() {
			It("returns -1", func() {
				Expect(mode.Parse("invalid")).To(Equal(mode.Mode(-1)))
			})
		})
	})

	Describe("Current", func() {
		It("returns the corresponding mode", func() {
			for modeStr, modeInt := range recognizedModes {
				viper.Set("mode", modeStr)
				Expect(mode.Current()).To(Equal(modeInt))
			}
		})

		Context("when the env var is not set", func() {
			It("returns the default mode", func() {
				viper.Set("mode", "")
				Expect(mode.Current()).To(Equal(mode.Mode(0)))
			})
		})

		Context("when the var is set to an unrecognized value", func() {
			It("returns the default mode", func() {
				viper.Set("mode", "invalid")
				var returnedMode mode.Mode

				_, stderr := termio.CaptureTermOut(func() { returnedMode = mode.Current() })
				Expect(stderr).To(ContainSubstring("Current mode `invalid` is invalid, using default of `configuration`\n"))
				Expect(returnedMode).To(Equal(mode.Mode(0)))
			})
		})
	})
})
