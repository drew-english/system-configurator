package pkg_test

import (
	"github.com/drew-english/system-configurator/cmd/pkg"
	"github.com/drew-english/system-configurator/internal/model"
	"github.com/drew-english/system-configurator/spec/stub/store"
	termio_stub "github.com/drew-english/system-configurator/spec/stub/termio"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Add", func() {
	var (
		stdout, stderr string
		cfg            *store.Configuration
		args           []string
	)

	subject := func() error {
		var err error
		stdout, stderr = termio_stub.CaptureTermOut(func() {
			err = pkg.AddCmd.RunE(nil, args)
		})

		return err
	}

	BeforeEach(func() {
		args = []string{"some-new-package@1.2.3", "some-other-package"}
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

	It("adds the package to the configuration", func() {
		Expect(subject()).To(Succeed())
		Expect(cfg.Packages).To(HaveLen(3))
		Expect(cfg.Packages[0]).To(BeComparableTo(&model.Package{
			Name:    "some-new-package",
			Version: "1.2.3",
		}))
		Expect(cfg.Packages[1]).To(BeComparableTo(&model.Package{
			Name: "some-other-package",
		}))
		Expect(stdout).To(Equal("Successfully added 2 packages\n"))
	})

	Context("when parsing a package fails", func() {
		BeforeEach(func() {
			args = []string{"invalid-package@"}
		})

		It("returns an error", func() {
			Expect(subject()).To(MatchError("failed to parse package string: invalid-package@"))
		})
	})

	Context("when loading the configuration fails", func() {
		JustBeforeEach(func() {
			store.StubLoadConfigurationError()
		})

		It("returns an error", func() {
			Expect(subject()).To(MatchError("Unable to load configuration: error loading configuration"))
		})
	})

	Context("when adding a package to the configuration fails", func() {
		BeforeEach(func() {
			args = []string{"some-package@1.2.3"}
		})

		It("prints a warning", func() {
			Expect(subject()).To(Succeed())
			Expect(stderr).To(ContainSubstring("Failed to add package `some-package@1.2.3`: package already exists in configuration"))
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
