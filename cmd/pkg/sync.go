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

var SyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync packages between configuration and system",
	Long: `Sync packages between configuration and system. Has different behavior based on the current mode:
- Configuration: Add packages to the configuration that are present only on the system.
- System: Add packages to the system that are present only in the configuration.
- Hybrid: Two-way sync packages between the configuration and system, only adding packages that are present in one but not the other.

Usage: scfg pkg sync`,
	RunE: func(cmd *cobra.Command, args []string) error {
		configPackages := make(map[string]*model.Package)
		cfg, err := store.LoadConfiguration()
		if err != nil {
			return fmt.Errorf("Unable to load configuration: %w", err)
		}

		cfgPkgList, err := cfg.ResolvedPkgs()
		if err != nil {
			termio.Warn("Unable to resolve packages for host manager, showing base configuration\n")
			cfgPkgList = cfg.Packages
		}

		for _, pkg := range cfgPkgList {
			configPackages[pkg.Name] = pkg
		}

		sysPackages := make(map[string]*model.Package)
		manager, err := pkgmanager.FindPackageManager()
		if err != nil {
			return fmt.Errorf("Failed to find the package manager: %w", err)
		}

		sysPkgList, err := manager.ListPackages()
		if err != nil {
			return fmt.Errorf("Unable to read system packages: %w", err)
		}

		for _, pkg := range sysPkgList {
			sysPackages[pkg.Name] = pkg
		}

		if mode.ManageConfig() {
			for name, pkg := range sysPackages {
				if _, ok := configPackages[name]; !ok {
					termio.Printf("[Configuration] Adding package `%s`\n", pkg)
					if err := cfg.AddPackage(pkg); err != nil {
						termio.Warnf("[Configuration] Failed to add package `%s`: %v\n", pkg, err)
					}
				}
			}

			if err := store.WriteConfiguration(cfg); err != nil {
				return fmt.Errorf("Failed to write configuration: %w", err)
			}
		}

		if mode.ManageSystem() {
			for name, pkg := range configPackages {
				if _, ok := sysPackages[name]; !ok {
					managerPackageName := manager.FmtPackageVersion(pkg)
					termio.Printf("[System] Adding package `%s`\n", managerPackageName)
					if err := manager.AddPackage(pkg); err != nil {
						termio.Warnf("[System] Failed to add package `%s`: %v\n", managerPackageName, err)
					}
				}
			}
		}

		return nil
	},
}

func init() {
	PkgCmd.AddCommand(SyncCmd)
}
