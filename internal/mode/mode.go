package mode

import (
	"os"

	"github.com/drew-english/system-configurator/pkg/termio"
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

func Set[M ~int](mode M) {
	os.Setenv(envVar, modeToS[Mode(mode)])
}

func Current() Mode {
	if modeStr := os.Getenv(envVar); modeStr != "" {
		mode := Parse(modeStr)
		if mode == -1 {
			termio.Errorf("Current SCFG_MODE `%s` is invalid, using default of `configuration`\n", modeStr)
			return defaultMode
		}

		return Mode(mode)
	}

	return defaultMode
}

func (m Mode) String() string {
	return modeToS[m]
}
