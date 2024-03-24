package mode

import (
	"github.com/drew-english/system-configurator/internal/mode"
	"github.com/drew-english/system-configurator/pkg/termio"
	"github.com/spf13/cobra"
)

var PeakCmd = &cobra.Command{
	Use:   "peak",
	Short: "View the current mode",
	Long: `View the currently set mode of system-configurator.

Usage: scfg mode peak`,
	RunE: func(cmd *cobra.Command, args []string) error {
		termio.Printf("%s\n", mode.Current().String())
		return nil
	},
}

func init() {
	ModeCmd.AddCommand(PeakCmd)
}
