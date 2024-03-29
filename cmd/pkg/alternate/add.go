package alternate

import (
	"fmt"
	"strings"

	"github.com/drew-english/system-configurator/internal/model"
	"github.com/drew-english/system-configurator/internal/store"
	"github.com/drew-english/system-configurator/pkg/sys/pkgmanager"
	"github.com/drew-english/system-configurator/pkg/termio"
	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"
)

var AddCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"a"},
	Short:   "Add an alternate",
	Long: `Add package alternates to a base package.
Alternate packages are specified in the form <package-name>[@<version>], where the version is optional.
Only modifies configuration, so modes have no effect.

Usage: scfg pkg alt add <base-package-name> <package-name>[@<version>] <manager-name>`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		basePkgName, altPkgName, mgrName := args[0], args[1], args[2]

		alternateToAdd, err := model.ParsePackage(altPkgName)
		if err != nil {
			return err
		}

		if _, ok := pkgmanager.Managers[mgrName]; !ok {
			return fmt.Errorf(
				"Invalid manager `%s`, valid managers are:\n%s\n",
				mgrName,
				strings.Join(maps.Keys(pkgmanager.Managers), "\n"),
			)
		}

		cfg, err := store.LoadConfiguration()
		if err != nil {
			return fmt.Errorf("Unable to load configuration: %w", err)
		}

		basePkg, _ := cfg.FindPackage(basePkgName)
		if basePkg == nil {
			return fmt.Errorf("Unable to find base package `%s`", basePkgName)
		}

		if err := basePkg.AddAlternate(mgrName, alternateToAdd); err != nil {
			return fmt.Errorf("Failed to add alternate `%s`: %w", altPkgName, err)
		}

		if err := store.WriteConfiguration(cfg); err != nil {
			return fmt.Errorf("Failed to write configuration: %w", err)
		}

		termio.Printf("Successfully added alternate `%s` to `%s`\n", altPkgName, basePkgName)
		return nil
	},
}

func init() {
	AlternateCmd.AddCommand(AddCmd)
}
