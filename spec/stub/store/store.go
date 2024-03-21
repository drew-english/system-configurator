package store

import (
	"errors"

	"github.com/drew-english/system-configurator/internal/store"
)

type Configuration store.Configuration

func StubLoadConfiguration(cfg *Configuration) {
	wrapLoadConfiguration((*store.Configuration)(cfg), nil)
}

func StubLoadConfigurationError() {
	wrapLoadConfiguration(nil, errors.New("error loading configuration"))
}

func StubWriteConfiguration() {
	wrapWriteConfiguration(nil)
}

func StubWriteConfigurationError() {
	wrapWriteConfiguration(errors.New("error writing configuration"))
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

func wrapWriteConfiguration(err error) {
	origianlWriteConfiguration := store.WriteConfiguration

	store.WriteConfiguration = func(*store.Configuration) error {
		defer func() {
			store.WriteConfiguration = origianlWriteConfiguration
		}()

		return err
	}
}
