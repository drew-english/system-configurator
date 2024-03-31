package model

import (
	"fmt"
	"regexp"

	"gopkg.in/yaml.v3"
)

var packageRegex = regexp.MustCompile(`^([^@\s]+)(?:$|@(\S+$))`)

type (
	Package struct {
		Name       string
		Version    string
		Alternates map[string]*Package // map of alternative package manager name to package info

		yamlStoredString string
	}

	yamlPkg struct {
		Name       string              `yaml:"name"`
		Version    string              `yaml:"version,omitempty"`
		Alternates map[string]*Package `yaml:"alternates,omitempty"`
	}
)

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

func (p *Package) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind != yaml.ScalarNode {
		decodedValue := new(yamlPkg)
		if err := value.Decode(decodedValue); err != nil {
			return err
		}

		p.Name = decodedValue.Name
		p.Version = decodedValue.Version
		p.Alternates = decodedValue.Alternates
		return nil
	}

	pkg, err := ParsePackage(value.Value)
	if err != nil {
		return err
	}

	*p = *pkg
	*&p.yamlStoredString = value.Value
	return nil
}

func (p *Package) MarshalYAML() (interface{}, error) {
	if p.yamlStoredString != "" {
		return p.yamlStoredString, nil
	}

	return &yamlPkg{
		Name:       p.Name,
		Version:    p.Version,
		Alternates: p.Alternates,
	}, nil
}
