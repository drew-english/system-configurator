package pkg

import (
	"github.com/drew-english/system-configurator/cmd/pkg/alternate"
	"github.com/spf13/cobra"
)

var PkgCmd = &cobra.Command{
	Use:     "package",
	Aliases: []string{"pkg"},
	Short:   "Manage configuration and system packages",
	Long:    "Manage configuration and system packages.",
}

func init() {
	PkgCmd.AddCommand(alternate.AlternateCmd)
}
