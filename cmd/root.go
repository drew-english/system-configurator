package cmd

import (
	"fmt"
	"os"

	"github.com/drew-english/system-configurator/cmd/pkg"
	"github.com/drew-english/system-configurator/internal/mode"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "scfg",
	Short: "System Configurator CLI",
	Long:  `A CLI tool for managing system packages and configuration scripts.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringP("mode", "m", "configuration", "Set the mode of system-configurator. Valid options are: conf[iguration], sys[tem], and hyb[rid].")
	viper.BindPFlag("mode", rootCmd.PersistentFlags().Lookup("mode"))
	viper.SetDefault("mode", "configuration")

	rootCmd.AddCommand(pkg.PkgCmd)
}

func initConfig() {
	viper.SetEnvPrefix("SCFG")
	viper.AutomaticEnv()

	if mode.Parse(viper.GetString("mode")) == -1 {
		cobra.CheckErr(fmt.Sprintf("mode `%s` is invalid\n", viper.GetString("mode")))
	}
}
