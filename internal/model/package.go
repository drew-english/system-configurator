package model

type Package struct {
	Name       string              `json:"name"`
	Version    string              `json:"version"`
	Alternates map[string]*Package `json:"alternates"` // map of alternative package manager name to package info
}

func (p *Package) ForManager(managerName string) *Package {
	if p.Alternates == nil {
		return p
	}

	return p.Alternates[managerName]
}

func (p *Package) String() string {
	return p.Name + " " + p.Version
}
