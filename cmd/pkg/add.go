package pkg

import (
	"fmt"

	"github.com/drew-english/system-configurator/internal/mode"
	"github.com/drew-english/system-configurator/internal/model"
	"github.com/drew-english/system-configurator/internal/store"
	"github.com/drew-english/system-configurator/pkg/sys/pkgmanager"
	"github.com/drew-english/system-configurator/pkg/termio"
	"github.com/spf13/cobra"
)

var AddCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"a"},
	Short:   "Add packages",
	Long: `Add an arbitrary number of packages.
Packages are specified in the form <package-name>[@<version>], where the version is optional.

Usage: scfg pkg add <package-name>@<version> <package-name> ...`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pkgsToAdd := make([]*model.Package, 0, len(args))
		for _, pkgStr := range args {
			pkg, err := model.ParsePackage(pkgStr)
			if err != nil {
				return err
			}

			pkgsToAdd = append(pkgsToAdd, pkg)
		}

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

		for _, pkg := range pkgsToAdd {
			if err := addPackage(cfg, manager, pkg); err != nil {
				return fmt.Errorf("Failed to add package `%s`: %v\n", pkg, err)
			}
		}

		if cfg != nil {
			if err := store.WriteConfiguration(cfg); err != nil {
				return fmt.Errorf("Failed to write configuration: %w", err)
			}
		}

		termio.Printf("Successfully added %d packages\n", len(pkgsToAdd))
		return nil
	},
}

func init() {
	PkgCmd.AddCommand(AddCmd)
}

func addPackage(cfg *store.Configuration, manager pkgmanager.PacakgeManager, pkg *model.Package) error {
	if manager != nil {
		if err := manager.AddPackage(pkg); err != nil {
			return err
		}
	}

	if cfg != nil {
		if err := cfg.AddPackage(pkg); err != nil {
			return err
		}
	}

	return nil
}

func modifySystem() bool {
	currentMode := mode.Current()
	return currentMode == mode.ModeHybrid || currentMode == mode.ModeSystem
}

func modifyConfig() bool {
	currentMode := mode.Current()
	return currentMode == mode.ModeHybrid || currentMode == mode.ModeConfiguration
}
