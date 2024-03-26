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

var _ = Describe("List", func() {
	var (
		stdout, stderr string
		cfg            *store.Configuration
		manager        string
		mode           string

		s = termio.Style()
	)

	subject := func() error {
		var err error
		stdout, stderr = termio_stub.CaptureTermOut(func() {
			err = pkg.ListCmd.RunE(nil, nil)
		})

		return err
	}

	BeforeEach(func() {
		manager = "apt"
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
		pkgmanager.StubFindPackageManager(manager)
		viper.Set("mode", mode)
	})

	AfterEach(func() {
		stdout = ""
		stderr = ""
	})

	It("lists the pacakges from the configuration", func() {
		Expect(subject()).To(Succeed())
		Expect(stdout).To(Equal("apt-some-package@1.2.3\n"))
		Expect(stderr).To(BeEmpty())
	})

	Context("when the configuration cannot be loaded", func() {
		JustBeforeEach(func() {
			store.StubLoadConfigurationError()
		})

		It("returns an error", func() {
			Expect(subject()).To(HaveOccurred())
			Expect(stdout).To(BeEmpty())
			Expect(stderr).To(BeEmpty())
		})
	})

	Context("when the package manager cannot be found", func() {
		JustBeforeEach(func() {
			pkgmanager.StubFindPackageManagerError()
		})

		It("logs a warning and uses the default configurations", func() {
			Expect(subject()).To(Succeed())
			Expect(stdout).To(Equal("some-package@1.2.3\n"))
			Expect(stderr).To(Equal(s.Yellow("WARNING: ") + "Unable to resolve packages for host manager, showing base configuration\n"))
		})
	})

	Context("when the mode manages the system", func() {
		var (
			commandStubs     *run.CommandStubManager
			teardownCmdStubs func(testing.TB)
		)

		BeforeEach(func() {
			mode = "system"
			commandStubs, teardownCmdStubs = run.StubCommand()
		})

		AfterEach(func() {
			teardownCmdStubs(GinkgoTB())
		})

		It("lists the packages from the system", func() {
			commandStubs.Register("apt list --installed", "apt-some-package/now 1.2.3")
			Expect(subject()).To(Succeed())
			Expect(stdout).To(Equal("apt-some-package@1.2.3\n"))
			Expect(stderr).To(BeEmpty())
		})

		Context("when the system packages cannot be listed", func() {
			It("returns an error", func() {
				commandStubs.RegisterError("apt list --installed", 1, "failed to list packages")
				Expect(subject()).To(HaveOccurred())
				Expect(stdout).To(BeEmpty())
				Expect(stderr).To(BeEmpty())
			})
		})

		Context("and the mode manages the configuration", func() {
			BeforeEach(func() {
				mode = "hybrid"
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
						{
							Name:    "config-only-pkg",
							Version: "2.3.4",
						},
					},
				}
			})

			It("lists the packages from the configuration and the system with a symbol indicating configuration status", func() {
				commandStubs.Register("apt list --installed", "apt-some-sys-package/now 1.2.3\napt-some-package/now 1.2.3")
				Expect(subject()).To(Succeed())
				Expect(stdout).To(Equal("  apt-some-package@1.2.3\n- apt-some-sys-package@1.2.3\n+ config-only-pkg@2.3.4\n"))
				Expect(stderr).To(BeEmpty())
			})
		})
	})
})
