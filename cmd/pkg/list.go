package pkg

import (
	"cmp"
	"fmt"
	"slices"

	"github.com/drew-english/system-configurator/internal/mode"
	"github.com/drew-english/system-configurator/internal/model"
	"github.com/drew-english/system-configurator/internal/store"
	"github.com/drew-english/system-configurator/pkg/sys/pkgmanager"
	"github.com/drew-english/system-configurator/pkg/termio"
	"github.com/spf13/cobra"
)

var ListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List packages",
	Long: `List packages.

Usage: scfg pkg list`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var configPackages []*model.Package
		if modifyConfig() {
			cfg, err := store.LoadConfiguration()
			if err != nil {
				return fmt.Errorf("Unable to load configuration: %w", err)
			}

			configPackages, err = cfg.ResolvedPkgs()
			if err != nil {
				termio.Warn("Unable to resolve packages for host manager, showing base configuration\n")
				configPackages = cfg.Packages
			}
		}

		var sysPackages []*model.Package
		if modifySystem() {
			manager, err := pkgmanager.FindPackageManager()
			if err != nil {
				return fmt.Errorf("Unable to load configuration: %w", err)
			}

			sysPackages, err = manager.ListPackages()
			if err != nil {
				return fmt.Errorf("Unable to read system packages: %w", err)
			}
		}

		slices.SortFunc(configPackages, func(x *model.Package, y *model.Package) int {
			return cmp.Compare(x.String(), y.String())
		})

		slices.SortFunc(sysPackages, func(x *model.Package, y *model.Package) int {
			return cmp.Compare(x.String(), y.String())
		})

		cfgIdx := 0
		sysIdx := 0
		for cfgIdx < len(configPackages) || sysIdx < len(sysPackages) {
			var cfgPkg, sysPkg *model.Package

			if cfgIdx < len(configPackages) {
				cfgPkg = configPackages[cfgIdx]
			}

			if sysIdx < len(sysPackages) {
				sysPkg = sysPackages[sysIdx]
			}

			cmp := comparePackage(cfgPkg, sysPkg)
			var pkgName, diffSign string

			switch cmp {
			case 0:
				pkgName = cfgPkg.String()
				diffSign = " "
				cfgIdx++
				sysIdx++
			case -1:
				pkgName = cfgPkg.String()
				diffSign = "+"
				cfgIdx++
			case 1:
				pkgName = sysPkg.String()
				diffSign = "-"
				sysIdx++
			}

			sfmt := "%s\n"
			printValues := []any{pkgName}
			if mode.Current() == mode.ModeHybrid {
				sfmt = "%s %s\n"
				printValues = []any{diffSign, pkgName}
			}

			termio.Printf(sfmt, printValues...)
		}

		return nil
	},
}

func init() {
	PkgCmd.AddCommand(ListCmd)
}

func comparePackage(a, b *model.Package) int {
	if a == nil && b == nil {
		return 0
	} else if a == nil {
		return 1
	} else if b == nil {
		return -1
	}

	return cmp.Compare(a.String(), b.String())
}
