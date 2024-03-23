package mode

import (
	"os"

	"github.com/drew-english/system-configurator/pkg/termio"
)

const (
	ModeConfiguration = iota
	ModeSystem
	ModeHybrid

	defaultMode = ModeConfiguration
	envVar      = "SCFG_MODE"
)

var sToMode = map[string]int{
	"conf":          ModeConfiguration,
	"configuration": ModeConfiguration,
	"system":        ModeSystem,
	"sys":           ModeSystem,
	"hybrid":        ModeHybrid,
	"hyb":           ModeHybrid,
}

var modeToS = map[int]string{
	ModeConfiguration: "configuration",
	ModeSystem:        "system",
	ModeHybrid:        "hybrid",
}

func Parse(mode string) int {
	if m, ok := sToMode[mode]; ok {
		return m
	}

	return -1
}

func Set(mode int) {
	os.Setenv(envVar, modeToS[mode])
}

func Current() int {
	if modeStr := os.Getenv(envVar); modeStr != "" {
		mode := Parse(modeStr)
		if mode == -1 {
			termio.Errorf("Current SCFG_MODE `%s` is invalid, using default of `configuration`\n", modeStr)
			return defaultMode
		}

		return mode
	}

	return defaultMode
}
