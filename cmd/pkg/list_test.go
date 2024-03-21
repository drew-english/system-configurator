package pkg_test

import (
	"github.com/drew-english/system-configurator/cmd/pkg"
	"github.com/drew-english/system-configurator/internal/model"
	"github.com/drew-english/system-configurator/pkg/termio"
	"github.com/drew-english/system-configurator/spec/stub/pkgmanager"
	"github.com/drew-english/system-configurator/spec/stub/store"
	termio_stub "github.com/drew-english/system-configurator/spec/stub/termio"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("List", func() {
	var (
		stdout, stderr string
		cfg            *store.Configuration
		manager        string

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
			Expect(stderr).To(Equal(s.Yellow("WARNING: ") + "Unable to resolve packages for host manager, showing raw configuration\n"))
		})
	})
})
