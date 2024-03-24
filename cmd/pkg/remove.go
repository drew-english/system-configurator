package pkg

import (
	"fmt"

	"github.com/drew-english/system-configurator/internal/store"
	"github.com/drew-english/system-configurator/pkg/sys/pkgmanager"
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
		var cfg *store.Configuration
		if modifyConfig() {
			var err error
			if cfg, err = store.LoadConfiguration(); err != nil {
				return fmt.Errorf("Unable to load configuration: %w", err)
			}
		}

		var manager pkgmanager.PacakgeManager
		if modifySystem() {
			var err error
			if manager, err = pkgmanager.FindPackageManager(); err != nil {
				return fmt.Errorf("Failed to resolve a package manager: %w", err)
			}
		}

		for _, pkgName := range args {
			if err := removePackage(cfg, manager, pkgName); err != nil {
				return fmt.Errorf("Failed to remove package `%s`: %s\n", pkgName, err)
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

func removePackage(cfg *store.Configuration, manager pkgmanager.PacakgeManager, pkgName string) error {
	if manager != nil {
		if err := manager.RemovePackage(pkgName); err != nil {
			return err
		}
	}

	if cfg != nil {
		if err := cfg.RemovePackage(pkgName); err != nil {
			return err
		}
	}

	return nil
}
