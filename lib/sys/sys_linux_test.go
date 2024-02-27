//go:build !android && linux

package sys_test

import (
	"fmt"
	"os"

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
	"other":  {"apk", "apt", "dnf", "snap", "pacman"},
}

var _ = Describe("Linux", func() {
	Describe("SupportedPackageManagers", func() {
		var (
			osVendor    string
			releasePath = "./tmp/os-release.txt"
		)

		BeforeEach(func() {
			sys.OSReleasePath = releasePath
			osVendor = "debian"
		})

		JustBeforeEach(func() {
			os.MkdirAll("./tmp", 0755)
			f, _ := os.Create(releasePath)
			f.WriteString(fmt.Sprintf("ID=%s", osVendor))
			f.Close()
		})

		AfterEach(func() {
			os.RemoveAll("./tmp")
		})

		for vendor, managers := range expectedManagers {
			Context(fmt.Sprintf("when the vendor is %s", vendor), func() {
				BeforeEach(func() {
					osVendor = vendor
				})

				It("returns the expected managers", func() {
					Expect(sys.SupportedPackageManagers()).To(Equal(managers))
				})
			})
		}
	})
})
