package main

import (
	"os"

	"github.com/gautampachnanda101/backstage-gen-cli/cmd"
)

var (
	// Version information (set during build)
	Version   = "dev"
	CommitSHA = "unknown"
	BuildDate = "unknown"
)

func main() {
	// Set version info for cobra commands
	cmd.SetVersionInfo(Version, CommitSHA, BuildDate)

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
