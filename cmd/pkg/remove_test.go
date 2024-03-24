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

var _ = Describe("Remove", func() {
	var (
		stdout string
		cfg    *store.Configuration
		args   []string
	)

	subject := func() error {
		var err error
		stdout, _ = termio_stub.CaptureTermOut(func() {
			err = pkg.RemoveCmd.RunE(nil, args)
		})

		return err
	}

	BeforeEach(func() {
		viper.Set("mode", "configuration")
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
	})

	It("removes the package to the configuration", func() {
		Expect(subject()).To(Succeed())
		Expect(cfg.Packages).To(HaveLen(0))
		Expect(stdout).To(Equal("Successfully removed 1 packages\n"))
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

		It("removes the package to the system", func() {
			commandStubs.Register("apt remove some-package", "package removeed successfully")
			Expect(subject()).To(Succeed())
		})

		Context("when removeing a package to the system fails", func() {
			It("returns an error", func() {
				commandStubs.RegisterError("apt remove some-package", 1, "failed to remove package")
				Expect(subject()).To(MatchError("Failed to remove package `some-package`: failed to remove package\napt: generic error\n"))
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
				commandStubs.Register("apt remove some-package", "package removeed successfully")
				Expect(subject()).To(Succeed())
				Expect(cfg.Packages).To(HaveLen(1))
			})
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

	Context("when removing a package from the configuration fails", func() {
		BeforeEach(func() {
			args = []string{"some-other-package"}
		})

		It("prints a warning", func() {
			Expect(subject()).To(MatchError("Failed to remove package `some-other-package`: package does not exist in configuration\n"))
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
