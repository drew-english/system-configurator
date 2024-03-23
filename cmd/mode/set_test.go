package mode_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	cmd_mode "github.com/drew-english/system-configurator/cmd/mode"
	"github.com/drew-english/system-configurator/internal/mode"
)

var _ = Describe("Set", func() {
	var modeArg string

	subject := func() error {
		return cmd_mode.SetCmd.RunE(nil, []string{modeArg})
	}

	BeforeEach(func() {
		mode.Set(mode.ModeHybrid)
		modeArg = "system"
	})

	It("sets the desired mode via the env var", func() {
		Expect(subject()).To(Succeed())
		Expect(mode.Current()).To(Equal(mode.ModeSystem))
	})

	Context("when the given mode is invalid", func() {
		BeforeEach(func() {
			modeArg = "invalid-mode"
		})

		It("returns an error", func() {
			Expect(subject()).To(MatchError("Invalid mode `invalid-mode`\n"))
		})
	})
})
