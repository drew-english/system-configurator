package store

import (
	"errors"

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

var WriteConfiguration = func(cfg *Configuration) error {
	s, err := NewLocal(nil)
	if err != nil {
		return err
	}

	return s.WriteConfiguration(cfg)
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

func (c *Configuration) AddPackage(pkg *model.Package) error {
	for i, p := range c.Packages {
		if p.Name == pkg.Name {
			return errors.New("package already exists in configuration")
		}

		if p.Name > pkg.Name {
			c.Packages = append(c.Packages[:i], append([]*model.Package{pkg}, c.Packages[i:]...)...)
			return nil
		}
	}

	c.Packages = append(c.Packages, pkg)
	return nil
}

func (c *Configuration) RemovePackage(name string) error {
	_, i := c.FindPackage(name)
	if i == -1 {
		return errors.New("package does not exist in configuration")
	}

	c.Packages = append(c.Packages[:i], c.Packages[i+1:]...)
	return nil
}

func (c *Configuration) FindPackage(name string) (*model.Package, int) {
	for i, p := range c.Packages {
		if p.Name == name {
			return p, i
		}
	}

	return nil, -1
}
