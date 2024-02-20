package store

import (
	"github.com/drew-english/system-configurator/lib/model"
)

type Store interface {
	LoadConfiguration() (*model.Configuration, error)
	WriteConfiguration(*model.Configuration) error
}
