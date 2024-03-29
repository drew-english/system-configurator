package alternate

import "github.com/spf13/cobra"

var AlternateCmd = &cobra.Command{
	Use:     "alternate",
	Aliases: []string{"alt"},
	Short:   "Manage package alternates",
	Long: `Manage package alternates.
Alternates provide the ability to specify a different package name and version for a given package manager.`,
}
