package store

type Store interface {
	LoadConfiguration() (*Configuration, error)
	WriteConfiguration(*Configuration) error
}

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
