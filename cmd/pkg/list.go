package pkg

import (
	"fmt"

	"github.com/drew-english/system-configurator/internal/store"
	"github.com/drew-english/system-configurator/pkg/termio"
	"github.com/spf13/cobra"
)

var ListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List packages",
	Long: `List packages.

Usage: system-configurator pkg list`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := store.LoadConfiguration()
		if err != nil {
			return fmt.Errorf("Unable to load configuration: %w", err)
		}

		pkgs, err := cfg.ResolvedPkgs()
		if err != nil {
			termio.Warn("Unable to resolve packages for host manager, showing raw configuration\n")
			pkgs = cfg.Packages
		}

		for _, pkg := range pkgs {
			termio.Print(pkg.String() + "\n")
		}

		return nil
	},
}

func init() {
	PkgCmd.AddCommand(ListCmd)
}
