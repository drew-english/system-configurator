package pkgmanager

import (
	"github.com/drew-english/system-configurator/pkg/sys/pkgmanager"
)

func StubFindPackageManager(mgrName string) {
	origianlFindPackageManager := pkgmanager.FindPackageManager

	pkgmanager.FindPackageManager = func() (pkgmanager.PacakgeManager, error) {
		defer func() {
			pkgmanager.FindPackageManager = origianlFindPackageManager
		}()

		return pkgmanager.Managers[mgrName], nil
	}
}
