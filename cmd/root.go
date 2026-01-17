package cmd

import (
	"fmt"

	"github.com/gautampachnanda101/backstage-gen-cli/pkg/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile      string
	version      string
	commitSHA    string
	buildDate    string
	verboseLevel int
	quietMode    bool
)

var rootCmd = &cobra.Command{
	Use:   "backstage-gen-cli",
	Short: "🚀 Backstage catalog generator and validator",
	Long: `╔════════════════════════════════════════════════════════════╗
║  Backstage Catalog Generator                              ║
╚════════════════════════════════════════════════════════════╝

A powerful CLI tool for generating and validating Backstage 
catalog-info.yaml files.

✨ Features:
  • Auto-detect technology stack from your repository
  • Generate catalog files with best practices
  • Validate against Backstage schema
  • Manage git hooks for automated validation
  • Cross-platform support (macOS, Linux, Windows)

📚 Quick Start:
  backstage-gen-cli inspect          # Inspect your repository
  backstage-gen-cli generate         # Generate catalog file
  backstage-gen-cli lint             # Validate catalog file
  backstage-gen-cli hooks install    # Set up git hooks

📖 Documentation:
  https://github.com/gautampachnanda101/backstage-gen-cli`,
	Version: version,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is .backstage-gen.yaml)")
	rootCmd.PersistentFlags().CountVarP(&verboseLevel, "verbose", "v", "increase verbosity (use -v, -vv, or -vvv)")
	rootCmd.PersistentFlags().BoolVarP(&quietMode, "quiet", "q", false, "suppress non-essential output")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("quiet", rootCmd.PersistentFlags().Lookup("quiet"))
}

func initConfig() {
	// Set verbosity level based on flags
	if quietMode {
		output.SetVerbosity(output.LevelQuiet)
	} else {
		// verboseLevel is 0 by default, 1 for -v, 2 for -vv, 3 for -vvv
		output.SetVerbosity(output.LevelNormal + verboseLevel)
	}

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
	if err := viper.ReadInConfig(); err == nil && output.IsVerbose() {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func SetVersionInfo(v, sha, date string) {
	version = v
	commitSHA = sha
	buildDate = date
	rootCmd.Version = fmt.Sprintf("%s (commit: %s, built: %s)", v, sha, date)
}
