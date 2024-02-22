//go:build !ios && darwin

package sys

var (
	supportedPackageManagers = []string{"brew"}
)

func SupportedPackageManagers() []string {
	return supportedPackageManagers
}
