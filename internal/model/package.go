package model

import (
	"fmt"
	"regexp"
)

var packageRegex = regexp.MustCompile(`^([^@\s]+)(?:$|@(\S+$))`)

type Package struct {
	Name       string              `json:"name"`
	Version    string              `json:"version,omitempty"`
	Alternates map[string]*Package `json:"alternates,omitempty"` // map of alternative package manager name to package info
}

func ParsePackage(pkgStr string) (*Package, error) {
	matches := packageRegex.FindStringSubmatch(pkgStr)
	if matches == nil {
		return nil, fmt.Errorf("failed to parse package string: %s", pkgStr)
	}

	return &Package{
		Name:    matches[1],
		Version: matches[2],
	}, nil
}

func (p *Package) ForManager(managerName string) *Package {
	if p.Alternates == nil {
		return p
	}

	if alternate := p.Alternates[managerName]; alternate != nil {
		return alternate
	}

	return p
}

func (p *Package) String() string {
	s := p.Name
	if p.Version != "" {
		s += "@" + p.Version
	}

	return s
}

func (p *Package) AddAlternate(managerName string, pkg *Package) error {
	if p.Alternates == nil {
		p.Alternates = make(map[string]*Package)
	}

	if _, ok := p.Alternates[managerName]; ok {
		return fmt.Errorf("alternate already exists for `%s`", managerName)
	}

	p.Alternates[managerName] = pkg
	return nil
}
