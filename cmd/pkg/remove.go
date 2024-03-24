package pkg

import (
	"fmt"

	"github.com/drew-english/system-configurator/internal/store"
	"github.com/drew-english/system-configurator/pkg/termio"
	"github.com/spf13/cobra"
)

var RemoveCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm"},
	Long: `Remove an arbitrary number of packages.
	Packages are specified in the form <package-name>.

	Usage: scfg pkg rm <package-name>...`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := store.LoadConfiguration()
		if err != nil {
			return fmt.Errorf("Unable to load configuration: %w", err)
		}

		for _, pkgName := range args {
			if err := cfg.RemovePackage(pkgName); err != nil {
				termio.Warnf("Failed to remove package `%s`: %s\n", pkgName, err)
			}
		}

		if err := store.WriteConfiguration(cfg); err != nil {
			return fmt.Errorf("Failed to write configuration: %w", err)
		}

		termio.Printf("Successfully removed %d packages\n", len(args))
		return nil
	},
}

func init() {
	PkgCmd.AddCommand(RemoveCmd)
}
