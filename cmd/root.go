package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile     string
	version     string
	commitSHA   string
	buildDate   string
	verboseMode bool
)

var rootCmd = &cobra.Command{
	Use:   "backstage-gen",
	Short: "Backstage catalog generator and validator",
	Long: `A CLI tool for generating and validating Backstage catalog-info.yaml files.

backstage-gen inspects your repository to automatically detect technology stack,
dependencies, and organizational patterns to generate properly formatted
Backstage catalog files.`,
	Version: version,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is .backstage-gen.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verboseMode, "verbose", "v", false, "verbose output")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName(".backstage-gen")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("$HOME")
	}
	viper.SetEnvPrefix("BACKSTAGE_GEN")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil && verboseMode {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func SetVersionInfo(v, sha, date string) {
	version = v
	commitSHA = sha
	buildDate = date
	rootCmd.Version = fmt.Sprintf("%s (commit: %s, built: %s)", v, sha, date)
}
