package wizard

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/gautampachnanda101/backstage-gen-cli/pkg/detector"
	"github.com/gautampachnanda101/backstage-gen-cli/pkg/generator"
	"github.com/gautampachnanda101/backstage-gen-cli/pkg/llm"
)

type Wizard struct {
	reader    *bufio.Reader
	llmClient *llm.Client
	cyan      *color.Color
	green     *color.Color
	yellow    *color.Color
	magenta   *color.Color
}

type WizardConfig struct {
	Name        string
	Description string
	Type        string
	Owner       string
	Lifecycle   string
	System      string
	Tags        []string
}

func New(llmClient *llm.Client) *Wizard {
	return &Wizard{
		reader:    bufio.NewReader(os.Stdin),
		llmClient: llmClient,
		cyan:      color.New(color.FgCyan, color.Bold),
		green:     color.New(color.FgGreen, color.Bold),
		yellow:    color.New(color.FgYellow, color.Bold),
		magenta:   color.New(color.FgMagenta),
	}
}

func (w *Wizard) Run(info *detector.RepositoryInfo, config *generator.Config) (*WizardConfig, error) {
	w.cyan.Println("\n╔════════════════════════════════════════════════════════════╗")
	w.cyan.Println("║          Interactive Catalog Generation Wizard            ║")
	w.cyan.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Println()

	wizConfig := &WizardConfig{}

	// Get LLM suggestions if available
	var llmSuggestions string
	if w.llmClient != nil && w.llmClient.IsAvailable() {
		w.yellow.Println("🤖 Getting AI-powered suggestions...")
		repoInfo := w.formatRepoInfo(info)
		suggestions, err := w.llmClient.GenerateCatalogSuggestions(repoInfo, "")
		if err == nil {
			llmSuggestions = suggestions
			fmt.Println("\n💡 AI Suggestions:")
			fmt.Println(w.wrapText(suggestions, 60))
			fmt.Println()
		}
	}

	// Component name
	wizConfig.Name = w.askQuestion(
		"Component Name",
		info.Name,
		"The unique identifier for your component",
	)

	// Description
	defaultDesc := info.Description
	if llmSuggestions != "" && strings.Contains(llmSuggestions, "description") {
		defaultDesc = w.extractFromSuggestions(llmSuggestions, "description", info.Description)
	}
	wizConfig.Description = w.askQuestion(
		"Description",
		defaultDesc,
		"A brief description of what this component does",
	)

	// Component type
	wizConfig.Type = w.askChoice(
		"Component Type",
		info.Type,
		[]string{"service", "library", "website", "resource"},
		"The type of component this represents",
	)

	// Owner
	defaultOwner := config.Defaults.Owner
	if info.DefaultOwner != "" {
		defaultOwner = info.DefaultOwner
	}
	wizConfig.Owner = w.askQuestion(
		"Owner",
		defaultOwner,
		"Team or person responsible for this component",
	)

	// Lifecycle
	wizConfig.Lifecycle = w.askChoice(
		"Lifecycle",
		config.Defaults.Lifecycle,
		[]string{"production", "experimental", "deprecated"},
		"Current lifecycle stage",
	)

	// System (optional)
	wizConfig.System = w.askQuestion(
		"System (optional)",
		config.Defaults.System,
		"The system this component belongs to",
	)

	// Tags
	defaultTags := w.generateDefaultTags(info)
	if llmSuggestions != "" {
		suggestedTags := w.extractTagsFromSuggestions(llmSuggestions, defaultTags)
		defaultTags = suggestedTags
	}

	tagsInput := w.askQuestion(
		"Tags (comma-separated)",
		strings.Join(defaultTags, ", "),
		"Relevant tags for discovery and categorization",
	)
	wizConfig.Tags = w.parseTags(tagsInput)

	w.green.Println("\n✓ Configuration complete!")
	return wizConfig, nil
}

func (w *Wizard) askQuestion(label, defaultValue, help string) string {
	w.yellow.Printf("📌 %s\n", label)
	if help != "" {
		fmt.Printf("   %s\n", w.magenta.Sprint(help))
	}
	if defaultValue != "" {
		fmt.Printf("   Default: %s\n", w.green.Sprint(defaultValue))
	}
	fmt.Print("   > ")

	input, _ := w.reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		return defaultValue
	}
	return input
}

func (w *Wizard) askChoice(label, defaultValue string, options []string, help string) string {
	w.yellow.Printf("📌 %s\n", label)
	if help != "" {
		fmt.Printf("   %s\n", w.magenta.Sprint(help))
	}
	fmt.Println("   Options:")
	for i, opt := range options {
		prefix := "  "
		if opt == defaultValue {
			prefix = w.green.Sprint("→ ")
		}
		fmt.Printf("   %s %d) %s\n", prefix, i+1, opt)
	}
	fmt.Print("   Choice (number or name): ")

	input, _ := w.reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		return defaultValue
	}

	// Try to parse as number
	for i, opt := range options {
		if input == fmt.Sprintf("%d", i+1) {
			return opt
		}
	}

	// Try to match name
	for _, opt := range options {
		if strings.EqualFold(input, opt) {
			return opt
		}
	}

	return defaultValue
}

func (w *Wizard) formatRepoInfo(info *detector.RepositoryInfo) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Name: %s\n", info.Name))
	sb.WriteString(fmt.Sprintf("Type: %s\n", info.Type))
	sb.WriteString(fmt.Sprintf("Languages: %s\n", strings.Join(info.Languages, ", ")))
	sb.WriteString(fmt.Sprintf("Frameworks: %s\n", strings.Join(info.Frameworks, ", ")))
	sb.WriteString(fmt.Sprintf("Build Tools: %s\n", strings.Join(info.BuildTools, ", ")))
	if info.HasDocker {
		sb.WriteString("Has Docker: Yes\n")
	}
	if info.HasKubernetes {
		sb.WriteString("Has Kubernetes: Yes\n")
	}
	return sb.String()
}

func (w *Wizard) generateDefaultTags(info *detector.RepositoryInfo) []string {
	tags := []string{}
	for _, lang := range info.Languages {
		tags = append(tags, strings.ToLower(lang))
	}
	for _, fw := range info.Frameworks {
		tags = append(tags, strings.ToLower(fw))
	}
	if info.HasDocker {
		tags = append(tags, "docker")
	}
	if info.HasKubernetes {
		tags = append(tags, "kubernetes")
	}
	return tags
}

func (w *Wizard) parseTags(input string) []string {
	if input == "" {
		return []string{}
	}
	parts := strings.Split(input, ",")
	tags := []string{}
	for _, part := range parts {
		tag := strings.TrimSpace(part)
		if tag != "" {
			tags = append(tags, tag)
		}
	}
	return tags
}

func (w *Wizard) wrapText(text string, width int) string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return text
	}

	var result strings.Builder
	lineLen := 0

	for _, word := range words {
		wordLen := len(word)
		if lineLen+wordLen+1 > width {
			result.WriteString("\n")
			lineLen = 0
		}
		if lineLen > 0 {
			result.WriteString(" ")
			lineLen++
		}
		result.WriteString(word)
		lineLen += wordLen
	}

	return result.String()
}

func (w *Wizard) extractFromSuggestions(suggestions, key, defaultValue string) string {
	lines := strings.Split(suggestions, "\n")
	for _, line := range lines {
		if strings.Contains(strings.ToLower(line), strings.ToLower(key)) {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				return strings.TrimSpace(parts[1])
			}
		}
	}
	return defaultValue
}

func (w *Wizard) extractTagsFromSuggestions(suggestions string, defaultTags []string) []string {
	lines := strings.Split(suggestions, "\n")
	for _, line := range lines {
		if strings.Contains(strings.ToLower(line), "tag") {
			if strings.Contains(line, ":") {
				parts := strings.Split(line, ":")
				if len(parts) > 1 {
					return w.parseTags(parts[1])
				}
			}
		}
	}
	return defaultTags
}
