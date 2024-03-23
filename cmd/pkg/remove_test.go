package pkg_test

import (
	"github.com/drew-english/system-configurator/cmd/pkg"
	"github.com/drew-english/system-configurator/internal/model"
	"github.com/drew-english/system-configurator/spec/stub/store"
	termio_stub "github.com/drew-english/system-configurator/spec/stub/termio"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Remove", func() {
	var (
		stdout, stderr string
		cfg            *store.Configuration
		args           []string
	)

	subject := func() error {
		var err error
		stdout, stderr = termio_stub.CaptureTermOut(func() {
			err = pkg.RemoveCmd.RunE(nil, args)
		})

		return err
	}

	BeforeEach(func() {
		args = []string{"some-package"}
		cfg = &store.Configuration{
			Packages: []*model.Package{
				{
					Name:    "some-package",
					Version: "1.2.3",
					Alternates: map[string]*model.Package{
						"apt": {
							Name:    "apt-some-package",
							Version: "1.2.3",
						},
					},
				},
			},
		}
	})

	JustBeforeEach(func() {
		store.StubLoadConfiguration(cfg)
		store.StubWriteConfiguration()
	})

	AfterEach(func() {
		stdout = ""
		stderr = ""
	})

	It("removes the package to the configuration", func() {
		Expect(subject()).To(Succeed())
		Expect(cfg.Packages).To(HaveLen(0))
		Expect(stdout).To(Equal("Successfully removed 1 packages\n"))
	})

	Context("when loading the configuration fails", func() {
		JustBeforeEach(func() {
			store.StubLoadConfigurationError()
		})

		It("returns an error", func() {
			Expect(subject()).To(MatchError("Unable to load configuration: error loading configuration"))
		})
	})

	Context("when removing a package from the configuration fails", func() {
		BeforeEach(func() {
			args = []string{"some-other-package"}
		})

		It("prints a warning", func() {
			Expect(subject()).To(Succeed())
			Expect(stderr).To(ContainSubstring("Failed to remove package `some-other-package`: package does not exist in configuration"))
		})
	})

	Context("when writing the configuration fails", func() {
		JustBeforeEach(func() {
			store.StubWriteConfigurationError()
		})

		It("returns an error", func() {
			Expect(subject()).To(MatchError("Failed to write configuration: error writing configuration"))
		})
	})
})
