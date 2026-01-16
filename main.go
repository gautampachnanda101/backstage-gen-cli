package main

import (
	"os"

	"github.com/yourusername/backstage-gen/cmd"
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
