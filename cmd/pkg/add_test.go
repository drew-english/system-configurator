package pkg_test

import (
	"testing"

	"github.com/drew-english/system-configurator/cmd/pkg"
	"github.com/drew-english/system-configurator/internal/model"
	"github.com/drew-english/system-configurator/spec/stub/pkgmanager"
	"github.com/drew-english/system-configurator/spec/stub/run"
	"github.com/drew-english/system-configurator/spec/stub/store"
	termio_stub "github.com/drew-english/system-configurator/spec/stub/termio"
	"github.com/spf13/viper"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Add", func() {
	var (
		stdout string
		cfg    *store.Configuration
		args   []string
	)

	subject := func() error {
		var err error
		stdout, _ = termio_stub.CaptureTermOut(func() {
			err = pkg.AddCmd.RunE(nil, args)
		})

		return err
	}

	BeforeEach(func() {
		viper.Set("mode", "configuration")
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

	Context("when in a mode that modifies the system", func() {
		var (
			commandStubs     *run.CommandStubManager
			teardownCmdStubs func(testing.TB)
		)

		BeforeEach(func() {
			viper.Set("mode", "hybrid")
			commandStubs, teardownCmdStubs = run.StubCommand()
			pkgmanager.StubFindPackageManager("apt")
		})

		AfterEach(func() {
			teardownCmdStubs(GinkgoTB())
		})

		It("adds the package to the system", func() {
			commandStubs.Register("apt install -y some-new-package=1.2.3", "package added successfully")
			commandStubs.Register("apt install -y some-other-package", "package added successfully")
			Expect(subject()).To(Succeed())
		})

		Context("when adding a package to the system fails", func() {
			It("returns an error", func() {
				commandStubs.RegisterError("apt install -y some-new-package=1.2.3", 1, "failed to find package")
				Expect(subject()).To(MatchError("Failed to add package `some-new-package@1.2.3`: failed to find package\napt: generic error\n"))
			})
		})

		Context("when the package manager cannot be found", func() {
			JustBeforeEach(func() {
				pkgmanager.StubFindPackageManagerError()
			})

			It("returns an error", func() {
				Expect(subject()).To(MatchError("Failed to resolve a package manager: unable to find a supported package manager on host system"))
			})
		})

		Context("and the mode does not modify the configuration", func() {
			BeforeEach(func() {
				viper.Set("mode", "system")
			})

			It("does not write the configuration", func() {
				commandStubs.Register("apt install -y some-new-package=1.2.3", "package added successfully")
				commandStubs.Register("apt install -y some-other-package", "package added successfully")
				Expect(subject()).To(Succeed())
				Expect(cfg.Packages).To(HaveLen(1))
			})
		})
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
			err := subject()
			Expect(err).To(MatchError("Failed to add package `some-package@1.2.3`: package already exists in configuration\n"))
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
