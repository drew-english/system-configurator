package mode

import (
	"fmt"

	"github.com/drew-english/system-configurator/internal/mode"
	"github.com/spf13/cobra"
)

var SetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set the mode",
	Long: `Set the mode of system-configurator.
Valid options are: conf[iguration], sys[tem], and hyb[rid].
conf[iguration]: Sets operations to only apply to the system-configurator configuration.
sys[tem]: Sets the operations to only apply to the system.
hyb[rid]: Sets the operations to apply to both the configuration and the system.

Usage: scfg mode set hybrid`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		modeVal := mode.Parse(args[0])
		if modeVal == -1 {
			return fmt.Errorf("Invalid mode `%s`\n", args[0])
		}

		mode.Set(modeVal)
		return nil
	},
}

func init() {
	ModeCmd.AddCommand(SetCmd)
}
