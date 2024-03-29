package alternate_test

import (
	"github.com/drew-english/system-configurator/cmd/pkg/alternate"
	"github.com/drew-english/system-configurator/internal/model"
	"github.com/drew-english/system-configurator/pkg/sys/pkgmanager"
	pkgmanager_stub "github.com/drew-english/system-configurator/spec/stub/pkgmanager"
	"github.com/drew-english/system-configurator/spec/stub/store"
	termio_stub "github.com/drew-english/system-configurator/spec/stub/termio"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Add", func() {
	var (
		basePkgName, altPkgName, mgrName string
		stdout, stderr                   string
		cfg                              *store.Configuration
		manager                          string
	)

	subject := func() error {
		var err error
		stdout, stderr = termio_stub.CaptureTermOut(func() {
			err = alternate.AddCmd.RunE(nil, []string{basePkgName, altPkgName, mgrName})
		})

		return err
	}

	BeforeEach(func() {
		manager = "apt"
		basePkgName = "some-package"
		altPkgName = "some-brew-package@2.2.2"
		mgrName = "brew"
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
		pkgmanager_stub.StubFindPackageManager(manager)
	})

	AfterEach(func() {
		stdout = ""
		stderr = ""
	})

	It("Adds the package alternate", func() {
		Expect(subject()).To(Succeed())
		Expect(stdout).To(Equal("Successfully added alternate `some-brew-package@2.2.2` to `some-package`\n"))
		Expect(stderr).To(BeEmpty())
		Expect(cfg.Packages[0].Alternates["brew"]).To(Equal(&model.Package{
			Name:    "some-brew-package",
			Version: "2.2.2",
		}))
	})

	Context("when the alternate package is invalid", func() {
		BeforeEach(func() {
			altPkgName = "some-invalid@"
		})

		It("returns an error", func() {
			Expect(subject()).To(MatchError("failed to parse package string: some-invalid@"))
			Expect(stdout).To(BeEmpty())
			Expect(stderr).To(BeEmpty())
		})
	})

	Context("when the manager is not supported", func() {
		BeforeEach(func() {
			mgrName = "invalid"
		})

		It("returns an error", func() {
			err := subject()
			Expect(err.Error()).To(ContainSubstring("Invalid manager `invalid`, valid managers are:\n"))
			for mgrName := range pkgmanager.Managers {
				Expect(err.Error()).To(ContainSubstring(mgrName))
			}

			Expect(stdout).To(BeEmpty())
			Expect(stderr).To(BeEmpty())
		})
	})

	Context("when the base package is not found", func() {
		BeforeEach(func() {
			basePkgName = "invalid"
		})

		It("returns an error", func() {
			Expect(subject()).To(MatchError("Unable to find base package `invalid`"))
			Expect(stdout).To(BeEmpty())
			Expect(stderr).To(BeEmpty())
		})
	})

	Context("when the alternate package already exists", func() {
		BeforeEach(func() {
			altPkgName = "apt-some-package"
			mgrName = "apt"
		})

		It("returns an error", func() {
			Expect(subject()).To(MatchError("Failed to add alternate `apt-some-package`: alternate already exists for `apt`"))
			Expect(stdout).To(BeEmpty())
			Expect(stderr).To(BeEmpty())
		})
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

	Context("when the configuration cannot be written", func() {
		JustBeforeEach(func() {
			store.StubWriteConfigurationError()
		})

		It("returns an error", func() {
			Expect(subject()).To(MatchError("Failed to write configuration: error writing configuration"))
			Expect(stdout).To(BeEmpty())
			Expect(stderr).To(BeEmpty())
		})
	})
})
