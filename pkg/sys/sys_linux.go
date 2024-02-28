//go:build !android && linux

package sys

import (
	"bufio"
	"os"
	"regexp"
	"slices"
	"strings"
)

var (
	OSReleasePath    = "/etc/os-release"
	reVendor         = regexp.MustCompile(`^ID=(.*)$`)
	supportedVendors = []string{"alpine", "arch", "debian", "fedora", "ubuntu"}

	managers = map[string][]string{
		"alpine": {"apk"},
		"arch":   {"pacman"},
		"debian": {"apt", "snap"},
		"fedora": {"dnf"},
		"ubuntu": {"apt", "snap"},
		"other":  {"apk", "apt", "dnf", "snap", "pacman"},
	}
)

func SupportedPackageManagers() []string {
	return managers[linuxArch()]
}

func linuxArch() (vendor string) {
	vendor = "other"

	f, err := os.Open(OSReleasePath)
	if err != nil {
		// TODO: log warning
		return
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() {
		if v := reVendor.FindStringSubmatch(s.Text()); v != nil {
			vendor = strings.Trim(v[1], `"`)
		}
	}

	if !slices.Contains(supportedVendors, vendor) {
		// TODO: log warning
		vendor = "other"
	}

	// TODO: log warning
	return
}
