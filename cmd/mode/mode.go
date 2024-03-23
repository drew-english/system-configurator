package mode

import (
	"github.com/spf13/cobra"
)

var ModeCmd = &cobra.Command{
	Use:   "mode",
	Short: "Manage the system-configurator mode",
	Long:  "Manage the system-configurator mode",
}
