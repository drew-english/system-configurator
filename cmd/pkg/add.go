package pkg

import (
	"fmt"

	"github.com/drew-english/system-configurator/internal/model"
	"github.com/drew-english/system-configurator/internal/store"
	"github.com/drew-english/system-configurator/pkg/termio"
	"github.com/spf13/cobra"
)

var AddCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"a"},
	Short:   "Add packages",
	Long: `Add an arbitrary number of packages.
Packages are specified in the form <package-name>[@<version>], where the version is optional.

Usage: system-configurator pkg add <package-name>@<version> <package-name> ...`,
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

		cfg, err := store.LoadConfiguration()
		if err != nil {
			return fmt.Errorf("Unable to load configuration: %w", err)
		}

		for _, pkg := range pkgsToAdd {
			if err := cfg.AddPackage(pkg); err != nil {
				termio.Warnf("Failed to add package `%s`: %v\n", pkg, err)
			}
		}

		if err := store.WriteConfiguration(cfg); err != nil {
			return fmt.Errorf("Failed to write configuration: %w", err)
		}

		termio.Printf("Successfully added %d packages\n", len(pkgsToAdd))
		return nil
	},
}

func init() {
	PkgCmd.AddCommand(AddCmd)
}
