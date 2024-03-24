package mode

import (
	"github.com/drew-english/system-configurator/pkg/termio"
	"github.com/spf13/viper"
)

const (
	ModeConfiguration = Mode(iota)
	ModeSystem
	ModeHybrid

	defaultMode = ModeConfiguration
	envVar      = "SCFG_MODE"
)

var sToMode = map[string]Mode{
	"conf":          ModeConfiguration,
	"configuration": ModeConfiguration,
	"system":        ModeSystem,
	"sys":           ModeSystem,
	"hybrid":        ModeHybrid,
	"hyb":           ModeHybrid,
}

var modeToS = map[Mode]string{
	ModeConfiguration: "configuration",
	ModeSystem:        "system",
	ModeHybrid:        "hybrid",
}

type Mode int

func Parse(mode string) Mode {
	if m, ok := sToMode[mode]; ok {
		return m
	}

	return -1
}

func Current() Mode {
	if modeStr := viper.GetString("mode"); modeStr != "" {
		mode := Parse(modeStr)
		if mode == -1 {
			termio.Errorf("Current mode `%s` is invalid, using default of `configuration`\n", modeStr)
			return defaultMode
		}

		return Mode(mode)
	}

	return defaultMode
}

func (m Mode) String() string {
	return modeToS[m]
}
