package pkg

import (
	"github.com/spf13/cobra"
)

var PkgCmd = &cobra.Command{
	Use:     "package",
	Aliases: []string{"pkg"},
	Short:   "Manage configuration and system packages",
	Long:    "Manage configuration and system packages.",
}
