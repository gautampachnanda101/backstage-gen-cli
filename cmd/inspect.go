package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/yourusername/backstage-gen/pkg/detector"
)

var jsonOutput bool

var inspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "Inspect repository and show detected information",
	RunE:  runInspect,
}

func init() {
	rootCmd.AddCommand(inspectCmd)
	inspectCmd.Flags().BoolVar(&jsonOutput, "json", false, "Output as JSON")
}

func runInspect(cmd *cobra.Command, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	det := detector.New(cwd)
	info, err := det.Detect()
	if err != nil {
		return fmt.Errorf("failed to detect: %w", err)
	}

	if jsonOutput {
		data, _ := json.MarshalIndent(info, "", "  ")
		fmt.Println(string(data))
		return nil
	}

	fmt.Println(color.CyanString("\n=== Repository Inspection ===\n"))
	fmt.Println(color.YellowString("Basic Information:"))
	fmt.Printf("  Name: %s\n", info.Name)
	fmt.Printf("  Type: %s\n", info.Type)
	fmt.Printf("  Path: %s\n", info.Path)
	
	if info.Description != "" {
		fmt.Printf("  Description: %s\n", info.Description)
	}

	if len(info.Languages) > 0 {
		fmt.Println(color.YellowString("\nTechnology:"))
		fmt.Printf("  Languages: %s\n", strings.Join(info.Languages, ", "))
	}

	if len(info.Frameworks) > 0 {
		fmt.Printf("  Frameworks: %s\n", strings.Join(info.Frameworks, ", "))
	}

	infraItems := []string{}
	if info.HasDocker {
		infraItems = append(infraItems, "Docker")
	}
	if info.HasKubernetes {
		infraItems = append(infraItems, "Kubernetes")
	}
	if info.HasHelm {
		infraItems = append(infraItems, "Helm")
	}

	if len(infraItems) > 0 {
		fmt.Println(color.YellowString("\nInfrastructure:"))
		fmt.Printf("  %s\n", strings.Join(infraItems, ", "))
	}

	fmt.Println(color.CyanString("\n=== End Inspection ===\n"))
	return nil
}
