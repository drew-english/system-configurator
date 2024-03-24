package mode_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	cmd_mode "github.com/drew-english/system-configurator/cmd/mode"
	"github.com/drew-english/system-configurator/internal/mode"
	"github.com/drew-english/system-configurator/spec/stub/termio"
)

var _ = Describe("Peak", func() {
	It("prints the currently set mode", func() {
		mode.Set(mode.ModeHybrid)
		stdout, _ := termio.CaptureTermOut(func() { cmd_mode.PeakCmd.RunE(nil, nil) })
		Expect(stdout).To(Equal("hybrid\n"))
	})
})
