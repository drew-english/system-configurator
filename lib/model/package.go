package model

type Package struct {
	Name       string              `json:"name"`
	Version    string              `json:"version"`
	Alternates map[string]*Package `json:"alternates"` // map of alternative package manager name to package info
}
