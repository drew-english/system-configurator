//go:build !ios && darwin

package sys_test

import (
	"github.com/drew-english/system-configurator/lib/sys"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Darwin", func() {
	Describe("SupportedPackageManagers", func() {
		It("returns a list of supported package managers", func() {
			Expect(sys.SupportedPackageManagers()).To(Equal([]string{"brew"}))
		})
	})
})
