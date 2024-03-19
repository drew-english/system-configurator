package store

import (
	"errors"

	"github.com/drew-english/system-configurator/internal/store"
)

type Configuration store.Configuration

func StubConfiguration(cfg *Configuration) {
	wrapLoadConfiguration((*store.Configuration)(cfg), nil)
}

func StubConfigurationError() {
	wrapLoadConfiguration(nil, errors.New("error loading configuration"))
}

func wrapLoadConfiguration(cfg *store.Configuration, err error) {
	origianlLoadConfiguration := store.LoadConfiguration

	store.LoadConfiguration = func() (*store.Configuration, error) {
		defer func() {
			store.LoadConfiguration = origianlLoadConfiguration
		}()

		return cfg, err
	}
}
