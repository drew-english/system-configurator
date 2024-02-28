package pkgmanager_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/drew-english/system-configurator/pkg/sys"
	"github.com/drew-english/system-configurator/pkg/sys/pkgmanager"
	"github.com/drew-english/system-configurator/spec/stub/run"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSys(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Sys Pkg Suite")
}

var _ = Describe("Pkg", func() {
	Describe("FindPackageManager", func() {
		It("returns a package manager", func() {
			unregister := run.StubFind(sys.SupportedPackageManagers()[0], nil)
			defer unregister()

			manager, err := pkgmanager.FindPackageManager()
			Expect(err).ToNot(HaveOccurred())
			Expect(manager).ToNot(BeNil())
		})

		Context("when no package manager is found", func() {
			It("returns an error", func() {
				allPkgs := strings.Join(sys.SupportedPackageManagers(), "|")
				unregister := run.StubFind(allPkgs, errors.New("not found"))
				defer unregister()

				manager, err := pkgmanager.FindPackageManager()
				Expect(err).To(HaveOccurred())
				Expect(manager).To(BeNil())
			})
		})
	})
})
