package pkg_test

import (
	"testing"

	"github.com/drew-english/system-configurator/cmd/pkg"
	"github.com/drew-english/system-configurator/internal/model"
	"github.com/drew-english/system-configurator/pkg/termio"
	"github.com/drew-english/system-configurator/spec/stub/pkgmanager"
	"github.com/drew-english/system-configurator/spec/stub/run"
	"github.com/drew-english/system-configurator/spec/stub/store"
	termio_stub "github.com/drew-english/system-configurator/spec/stub/termio"
	"github.com/spf13/viper"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Sync", func() {
	var (
		stdout, stderr string
		cfg            *store.Configuration
		manager        string
		mode           string

		s                = termio.Style()
		commandStubs     *run.CommandStubManager
		teardownCmdStubs func(testing.TB)
	)

	subject := func() error {
		var err error
		stdout, stderr = termio_stub.CaptureTermOut(func() {
			err = pkg.SyncCmd.RunE(nil, nil)
		})

		return err
	}

	BeforeEach(func() {
		manager = "apt"
		commandStubs, teardownCmdStubs = run.StubCommand()
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
		pkgmanager.StubFindPackageManager(manager)
		pkgmanager.StubFindPackageManager(manager)
		viper.Set("mode", mode)
	})

	AfterEach(func() {
		stdout = ""
		stderr = ""
		teardownCmdStubs(GinkgoTB())
	})

	It("syncs the system pacakges to the configuration", func() {
		commandStubs.Register("apt list --installed", "apt-some-sys-package/now 1.2.3")
		Expect(subject()).To(Succeed())
		Expect(stdout).To(Equal("[Configuration] Adding package `apt-some-sys-package@1.2.3`\n"))
		Expect(stderr).To(BeEmpty())
		Expect(cfg.Packages).To(ContainElement(&model.Package{
			Name:    "apt-some-sys-package",
			Version: "1.2.3",
		}))
	})

	Context("when the configuration cannot be loaded", func() {
		JustBeforeEach(func() {
			store.StubLoadConfigurationError()
		})

		It("returns an error", func() {
			Expect(subject()).To(MatchError("Unable to load configuration: error loading configuration"))
			Expect(stdout).To(BeEmpty())
			Expect(stderr).To(BeEmpty())
		})
	})

	Context("when the package manager cannot be found", func() {
		JustBeforeEach(func() {
			pkgmanager.StubFindPackageManagerError()
			pkgmanager.StubFindPackageManagerError()
		})

		It("logs a warning and returns an error", func() {
			Expect(subject()).To(MatchError("Failed to find the package manager: unable to find a supported package manager on host system"))
			Expect(stdout).To(BeEmpty())
			Expect(stderr).To(Equal(s.Yellow("WARNING: ") + "Unable to resolve packages for host manager, showing base configuration\n"))
		})
	})

	Context("when the system packages cannot be listed", func() {
		It("returns an error", func() {
			commandStubs.RegisterError("apt list --installed", 1, "failed to list packages")
			Expect(subject()).To(MatchError("Unable to read system packages: failed to list packages\napt: generic error"))
			Expect(stdout).To(BeEmpty())
			Expect(stderr).To(BeEmpty())
		})
	})

	Context("when the configuration cannot be written", func() {
		JustBeforeEach(func() {
			store.StubWriteConfigurationError()
		})

		It("returns an error", func() {
			commandStubs.Register("apt list --installed", "apt-some-sys-package/now 1.2.3")
			Expect(subject()).To(MatchError("Failed to write configuration: error writing configuration"))
			Expect(stdout).To(Equal("[Configuration] Adding package `apt-some-sys-package@1.2.3`\n"))
			Expect(stderr).To(BeEmpty())
		})
	})

	Context("when the mode manages the system", func() {
		BeforeEach(func() {
			mode = "system"
			commandStubs.Register("apt list --installed", "apt-some-sys-package/now 1.2.3")
		})

		It("syncs the packages to the system from the configuration", func() {
			commandStubs.Register("apt install -y apt-some-package=1.2.3", "successfully installed package")
			Expect(subject()).To(Succeed())
			Expect(stdout).To(Equal("[System] Adding package `apt-some-package=1.2.3`\n"))
			Expect(stderr).To(BeEmpty())
		})

		Context("when the pacakge manager fails to add the package", func() {
			It("logs a warning and continues", func() {
				commandStubs.RegisterError("apt install -y apt-some-package=1.2.3", 1, "failed to install package")
				Expect(subject()).To(Succeed())
				Expect(stdout).To(Equal("[System] Adding package `apt-some-package=1.2.3`\n"))
				Expect(stderr).To(Equal(s.Yellow("WARNING: ") + "[System] Failed to add package `apt-some-package=1.2.3`: failed to install package\napt: generic error\n"))
			})
		})

		Context("and the mode manages the configuration", func() {
			BeforeEach(func() {
				mode = "hybrid"
			})

			It("syncs both the system and configuration packages in an addition only manner", func() {
				commandStubs.Register("apt install -y apt-some-package=1.2.3", "successfully installed package")
				Expect(subject()).To(Succeed())
				Expect(stdout).To(Equal("[Configuration] Adding package `apt-some-sys-package@1.2.3`\n[System] Adding package `apt-some-package=1.2.3`\n"))
				Expect(stderr).To(BeEmpty())
			})
		})
	})
})
