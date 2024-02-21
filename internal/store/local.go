package store

import (
	"encoding/json"
	"errors"
	"os"
	"path"

	"github.com/drew-english/system-configurator/internal/model"
)

const (
	LocalDefaultLocation = "~/.config/system-configurator"
	LocalDefaultFileName = "config.json"
)

type (
	LocalCfg struct {
		Location string
		FileName string
	}

	localStore struct {
		cfg        *LocalCfg
		configFile *os.File
	}
)

func NewLocal(cfg *LocalCfg) (s Store, err error) {
	localStore := &localStore{cfg: cfg}
	s = localStore

	file, err := localStore.localConfigFile()
	if err != nil {
		return
	}

	localStore.configFile = file
	return
}

func (ls *localStore) LoadConfiguration() (*model.Configuration, error) {
	if ls.configFile == nil {
		return nil, errors.New("error referencing local configuration file")
	}

	configData := &model.Configuration{}
	decoder := json.NewDecoder(ls.configFile)
	err := decoder.Decode(configData)
	if err != nil {
		return nil, err
	}

	return configData, nil
}

func (ls *localStore) WriteConfiguration(configData *model.Configuration) error {
	if ls.configFile == nil {
		return errors.New("error referencing local configuration file")
	}

	if configData == nil {
		return errors.New("configuration data cannot be nil")
	}

	data, err := json.Marshal(configData)
	if err != nil {
		return err
	}

	_, err = ls.configFile.Write(data)
	return err
}

func (ls *localStore) localConfigFile() (*os.File, error) {
	file, err := os.Open(ls.cfg.filePath())
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ls.createLocalConfigFile()
		}

		return nil, err
	}

	return file, nil
}

func (ls *localStore) createLocalConfigFile() (*os.File, error) {
	err := os.MkdirAll(ls.cfg.location(), 0755)
	if err != nil {
		return nil, err
	}

	file, err := os.Create(ls.cfg.filePath())
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (lc *LocalCfg) filePath() string {
	return path.Join(lc.location(), lc.fileName())
}

func (lc *LocalCfg) location() string {
	if lc != nil && lc.Location != "" {
		return lc.Location
	}

	return LocalDefaultLocation
}

func (lc *LocalCfg) fileName() string {
	if lc != nil && lc.FileName != "" {
		return lc.FileName
	}

	return LocalDefaultFileName
}