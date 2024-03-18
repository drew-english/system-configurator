// Package management on the host system.
package pkgmanager

import (
	"errors"
	"regexp"
	"strings"
	"text/template"

	"github.com/drew-english/system-configurator/internal/model"
	"github.com/drew-english/system-configurator/pkg/run"
	"github.com/drew-english/system-configurator/pkg/sys"
)

type (
	PacakgeManager interface {
		Name() string
		AddPackage(*model.Package) error
		RemovePackage(*model.Package) error
		ListPackages() ([]*model.Package, error)
	}

	basePackageManager struct {
		BaseCmd          string
		AddCmd           []string
		RemoveCmd        []string
		ListCmd          []string
		listParsePattern *regexp.Regexp
		versionTmpl      *template.Template
	}
)

var (
	cachedManager PacakgeManager
)

var FindPackageManager = func() (PacakgeManager, error) {
	if cachedManager != nil {
		return cachedManager, nil
	}

	possibleManagers := sys.SupportedPackageManagers()
	for _, mgr := range possibleManagers {
		if _, err := run.Find(mgr); err == nil && Managers[mgr] != nil {
			cachedManager = Managers[mgr]
			return Managers[mgr], nil
		}
	}

	return nil, errors.New("unable to find a supported package manager on host system")
}

func ResetCachedManager() {
	cachedManager = nil
}

func (pm *basePackageManager) Name() string {
	return pm.BaseCmd
}

func (pm *basePackageManager) AddPackage(pkg *model.Package) error {
	args := append(pm.AddCmd, pm.fmtPackageVersion(pkg))
	return run.Command(pm.BaseCmd, args...).Run()
}

func (pm *basePackageManager) RemovePackage(pkg *model.Package) error {
	args := append(pm.RemoveCmd, pkg.Name)
	return run.Command(pm.BaseCmd, args...).Run()
}

func (pm *basePackageManager) ListPackages() ([]*model.Package, error) {
	out, err := run.Command(pm.BaseCmd, pm.ListCmd...).Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(out), "\n")
	pkgs := make([]*model.Package, 0, len(lines))

	for _, line := range lines {
		if line == "" {
			continue
		}

		if pkg := pm.parsePgk(line); pkg != nil {
			pkgs = append(pkgs, pkg)
		}
	}

	return pkgs, nil
}

func (pm *basePackageManager) parsePgk(line string) *model.Package {
	matches := pm.listParsePattern.FindStringSubmatch(line)
	if matches == nil {
		return nil
	}

	return &model.Package{
		Name:    matches[1],
		Version: matches[2],
	}
}

func (pm *basePackageManager) fmtPackageVersion(pkg *model.Package) string {
	if pkg.Version == "" {
		return pkg.Name
	}

	var buf strings.Builder
	if err := pm.versionTmpl.Execute(&buf, pkg); err != nil {
		return pkg.Name
	}

	return buf.String()
}
