package mode_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/drew-english/system-configurator/internal/mode"
	"github.com/drew-english/system-configurator/spec/stub/termio"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestMode(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mode Suite")
}

var _ = Describe("Mode", func() {
	var recognizedModes = map[string]int{
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
				Expect(mode.Parse("invalid")).To(Equal(-1))
			})
		})
	})

	Describe("Set", func() {
		It("sets the env var based on the given mode", func() {
			mode.Set(mode.ModeConfiguration)
			Expect(os.Getenv("SCFG_MODE")).To(Equal("configuration"))
			mode.Set(mode.ModeSystem)
			Expect(os.Getenv("SCFG_MODE")).To(Equal("system"))
			mode.Set(mode.ModeHybrid)
			Expect(os.Getenv("SCFG_MODE")).To(Equal("hybrid"))
			mode.Set(47)
			Expect(os.Getenv("SCFG_MODE")).To(Equal(""))
		})
	})

	Describe("Current", func() {
		It("returns the corresponding mode", func() {
			for modeStr, modeInt := range recognizedModes {
				os.Setenv("SCFG_MODE", modeStr)
				Expect(mode.Current()).To(Equal(modeInt))
			}
		})

		Context("when the env var is not set", func() {
			It("returns the default mode", func() {
				os.Setenv("SCFG_MODE", "")
				Expect(mode.Current()).To(Equal(0))
			})
		})

		Context("when the env var is set to an unrecognized value", func() {
			It("returns the default mode", func() {
				os.Setenv("SCFG_MODE", "invalid")
				var returnedMode int

				_, stderr := termio.CaptureTermOut(func() { returnedMode = mode.Current() })
				Expect(stderr).To(ContainSubstring("Current SCFG_MODE `invalid` is invalid, using default of `configuration`\n"))
				Expect(returnedMode).To(Equal(0))
			})
		})
	})
})
