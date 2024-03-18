package store

import (
	"github.com/drew-english/system-configurator/internal/model"
	"github.com/drew-english/system-configurator/pkg/sys/pkgmanager"
)

type (
	Store interface {
		LoadConfiguration() (*Configuration, error)
		WriteConfiguration(*Configuration) error
	}

	Configuration struct {
		Packages []*model.Package `json:"packages"`
	}
)

var LoadConfiguration = func() (*Configuration, error) {
	s, err := NewLocal(nil)
	if err != nil {
		return nil, err
	}

	return s.LoadConfiguration()
}

func (c *Configuration) ResolvedPkgs() ([]*model.Package, error) {
	manager, err := pkgmanager.FindPackageManager()
	if err != nil {
		return nil, err
	}

	resolvedPackages := make([]*model.Package, 0, len(c.Packages))
	for _, pkg := range c.Packages {
		resolvedPackages = append(resolvedPackages, pkg.ForManager(manager.Name()))
	}

	return resolvedPackages, nil
}
