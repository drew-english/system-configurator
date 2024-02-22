//go:build !android && linux

package sys_test

import (
	"fmt"

	"github.com/drew-english/system-configurator/lib/sys"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var expectedManagers = map[string][]string{
	"alpine": {"apk"},
	"arch":   {"pacman"},
	"debian": {"apt", "snap"},
	"fedora": {"dnf"},
	"ubuntu": {"apt", "snap"},
	"other":  {"apk", "apt", "dnf", "flatpak", "snap", "pacman"},
}

var _ = Describe("Linux", func() {
	Describe("SupportedPackageManagers", func() {
		var (
			osVendor   string
			releaseDir = "./tmp/os-release"
		)

		BeforeEach(func() {
			sys.OSReleaseDirectory = releaseDir
			osVendor = "debian"
		})

		JustBeforeEach(func() {
			os.WriteFile(releaseDir, fmt.Sprintf("ID=%s", osVendor, 0644))
		})

		for vendor, managers := range expectedManagers {
			Context(fmt.Sprintf("when the vendor is %s", vendor), func() {
				It("returns the expected managers", func() {
					Expect(sys.SupportedPackageManagers()).To(Equal(expectedManagers[vendor]))
				})
			})
		}
	})
})
