package pkgmanager

import (
	"errors"

	"github.com/drew-english/system-configurator/pkg/sys/pkgmanager"
)

func StubFindPackageManager(mgrName string) {
	wrapFindPackageManager(pkgmanager.Managers[mgrName], nil)
}

func StubFindPackageManagerError() {
	wrapFindPackageManager(nil, errors.New("unable to find a supported package manager on host system"))
}

func wrapFindPackageManager(mgr pkgmanager.PacakgeManager, err error) {
	origianlFindPackageManager := pkgmanager.FindPackageManager

	pkgmanager.FindPackageManager = func() (pkgmanager.PacakgeManager, error) {
		defer func() {
			pkgmanager.FindPackageManager = origianlFindPackageManager
		}()

		return mgr, err
	}
}
