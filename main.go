package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"Konfetti/config"
	"Konfetti/parser"
	"Konfetti/scanner"

	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v2"
)

const version = "2.0.0"

type ConfigResult struct {
	File     string                 `json:"file"`
	Format   string                 `json:"format"`
	Settings map[string]interface{} `json:"settings"`
}

func main() {
	app := &cli.App{
		Name:     "Konfetti 🎉",
		Usage:    "Wrangle your config chaos with style",
		Version:  version,
		Compiled: time.Now(),
		Commands: []*cli.Command{
			{
				Name:    "scan",
				Aliases: []string{"s"},
				Usage:   "Scan directories for config files",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "profile", Usage: "Use a named profile from ~/.konfetti.yaml"},
					&cli.StringFlag{Name: "path", Aliases: []string{"p"}, Usage: "Path to scan"},
					&cli.StringFlag{Name: "key", Usage: "Filter by key substring"},
					&cli.StringFlag{Name: "value", Usage: "Filter by value substring"},
					&cli.StringFlag{Name: "filter", Usage: "Filter filenames that contain substring"},
					&cli.StringFlag{Name: "output", Usage: "Output format: text, json, table", Value: "text"},
					&cli.BoolFlag{Name: "interactive", Aliases: []string{"i"}, Usage: "Enable interactive mode"},
					&cli.BoolFlag{Name: "no-warn", Usage: "Suppress warning output for unreadable paths"},
				},
				Action: scanCommand,
			},
			{
				Name:    "explore",
				Aliases: []string{"x"},
				Usage:   "Interactive config browser",
				Action:  exploreCommand,
			},
			{
				Name:   "explain",
				Usage:  "Explain what a config file does (experimental)",
				Action: explainCommand,
			},
			{
				Name:  "init",
				Usage: "Generate a sample config file at ~/.konfetti.yaml",
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "force", Aliases: []string{"f"}, Usage: "Overwrite existing config file"},
				},
				Action: initCommand,
			},
			{
				Name:  "tips",
				Usage: "Show a random Konfetti tip",
				Action: func(c *cli.Context) error {
					tips := []string{
						"Pro tip: Add `--interactive` to get guided scanning.",
						"Not sure where to start? Try `konfetti scan --path ~`.",
						"Keep configs tidy — chaos loves clutter.",
					}
					fmt.Println("💡 " + tips[time.Now().Unix()%int64(len(tips))])
					return nil
				},
			},
		},
		Before: func(c *cli.Context) error {
			fmt.Fprintln(os.Stderr, "🎉 Welcome to Konfetti – It doesn’t judge your configs. Much. 🎉")
			return nil
		},
		After: func(c *cli.Context) error {
			fmt.Fprintln(os.Stderr, "\n✨ All scanned. No judgment (mostly).")
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println("❌ Error:", err)
		os.Exit(1)
	}
}

// ---------------- Scan Command ----------------
func hasStdinData() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	// If not a char device and size > 0 or pipe
	if (fi.Mode() & os.ModeCharDevice) == 0 {
		return true
	}
	return false
}

func scanCommand(c *cli.Context) error {
	// Load config file
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("[WARN] Could not load config file: %v\n", err)
		cfg = &config.Config{
			Defaults: config.ScanDefaults{Output: "text"},
			Profiles: make(map[string]config.ScanProfile),
		}
	}

	// Start with defaults from config
	path := cfg.Defaults.Path
	filterKey := cfg.Defaults.Key
	filterValue := cfg.Defaults.Value
	filterName := cfg.Defaults.Filter
	outputFormat := cfg.Defaults.Output
	noWarn := cfg.Defaults.NoWarn

	// Apply profile if specified
	if profileName := c.String("profile"); profileName != "" {
		if profile, ok := cfg.Profiles[profileName]; ok {
			if profile.Path != "" {
				path = profile.Path
			}
			if profile.Key != "" {
				filterKey = profile.Key
			}
			if profile.Value != "" {
				filterValue = profile.Value
			}
			if profile.Filter != "" {
				filterName = profile.Filter
			}
			if profile.Output != "" {
				outputFormat = profile.Output
			}
			noWarn = profile.NoWarn
		} else {
			return fmt.Errorf("profile '%s' not found in ~/.konfetti.yaml", profileName)
		}
	}

	// CLI flags override everything (only if explicitly set)
	if c.IsSet("path") {
		path = c.String("path")
	}
	if c.IsSet("key") {
		filterKey = c.String("key")
	}
	if c.IsSet("value") {
		filterValue = c.String("value")
	}
	if c.IsSet("filter") {
		filterName = c.String("filter")
	}
	if c.IsSet("output") {
		outputFormat = c.String("output")
	}
	if c.IsSet("no-warn") {
		noWarn = c.Bool("no-warn")
	}

	interactive := c.Bool("interactive")

	// STDIN mode: no path provided but data is piped in
	if path == "" && hasStdinData() {
		data, err := os.ReadFile("/dev/stdin")
		if err != nil {
			return err
		}
		parsed, format := parser.ParseBytes(data)
		// Apply filters
		if filterKey != "" || filterValue != "" {
			filtered := make(map[string]interface{})
			for k, v := range parsed {
				keyMatch := filterKey == "" || strings.Contains(strings.ToLower(k), strings.ToLower(filterKey))
				valMatch := filterValue == "" || strings.Contains(strings.ToLower(fmt.Sprintf("%v", v)), strings.ToLower(filterValue))
				if keyMatch && valMatch {
					filtered[k] = v
				}
			}
			parsed = filtered
		}
		result := []ConfigResult{{File: "stdin", Format: format, Settings: parsed}}
		switch outputFormat {
		case "json":
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			enc.Encode(result)
		case "table":
			printTable(result)
		default:
			fmt.Printf("File: stdin [%s]\n", format)
			for k, v := range parsed {
				fmt.Printf("  %s = %v\n", k, v)
			}
		}
		return nil
	}

	if path == "" && interactive {
		prompt := promptui.Prompt{
			Label:     "No path specified, scan current directory?",
			IsConfirm: true,
		}
		_, err := prompt.Run()
		if err != nil {
			fmt.Println("Aborted.")
			return nil
		}
		path = "."
	}

	if path == "" {
		cwd, err := os.Getwd()
		if err == nil {
			path = cwd
		} else {
			fmt.Println("Using default scan paths")
			// Will handle in getDefaultScanPaths()
			paths := getDefaultScanPaths()
			return runScan(paths, filterName, filterKey, filterValue, outputFormat, noWarn)
		}
	}

	return runScan([]string{path}, filterName, filterKey, filterValue, outputFormat, noWarn)
}

func runScan(paths []string, filterName, filterKey, filterValue, outputFormat string, suppressWarn bool) error {
	exts := []string{".json", ".yaml", ".yml", ".xml", ".conf", ".config", ".txt", ".ini", ".properties"}
	results, scanErrors := ScanAndFilter(paths, exts, filterName, filterKey, filterValue)

	fmt.Printf("Matched %d config files:\n", len(results))
	if len(scanErrors) > 0 && !suppressWarn {
		fmt.Println("Encountered the following errors while scanning:")
		for _, err := range scanErrors {
			fmt.Printf("  [WARN] %s\n", err)
		}
	}
	if len(results) == 0 {
		fmt.Println("No matches found.")
		return nil
	}

	switch outputFormat {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(results)
	case "table":
		printTable(results)
	default:
		for _, r := range results {
			fmt.Printf("File: %s [%s]\n", r.File, r.Format)
			for k, v := range r.Settings {
				fmt.Printf("  %s = %v\n", k, v)
			}
			fmt.Println("---")
		}
	}

	return nil
}

// ---------------- Explore Command ----------------
func exploreCommand(c *cli.Context) error {
	fmt.Println("🧭 Interactive explorer coming soon (think Bubble Tea TUI)")
	return nil
}

// ---------------- Init Command ----------------
func initCommand(c *cli.Context) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not find home directory: %w", err)
	}

	configPath := filepath.Join(homeDir, ".konfetti.yaml")

	// Check if config already exists
	if _, err := os.Stat(configPath); err == nil && !c.Bool("force") {
		return fmt.Errorf("config file already exists at %s. Use --force to overwrite", configPath)
	}

	// Write sample config
	sampleConfig := config.GetSampleConfig()
	if err := os.WriteFile(configPath, []byte(sampleConfig), 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	fmt.Printf("✅ Created sample config file at: %s\n", configPath)
	fmt.Println("\n📝 Edit this file to customize your default settings and profiles.")
	fmt.Println("💡 Use: konfetti scan --profile <name> to use a profile.")
	return nil
}

// ---------------- Explain Command ----------------
func explainCommand(c *cli.Context) error {
	if c.Args().Len() == 0 {
		if hasStdinData() {
			data, err := os.ReadFile("/dev/stdin")
			if err != nil {
				return err
			}
			parsed, format := parser.ParseBytes(data)
			fmt.Println("🔎 Analyzing stdin config ...")
			keys := make([]string, 0, len(parsed))
			for k := range parsed {
				keys = append(keys, k)
			}
			fmt.Printf("Format: %s | Keys detected: %d\n", format, len(keys))
			// Mock explanation summary
			sample := "Looks like a configuration file with typical settings. Ensure sensitive values are not exposed."
			if _, ok := parsed["debug"]; ok {
				sample += " Debug flag present; disable in production."
			}
			fmt.Println("✨ " + sample)
			return nil
		}
		fmt.Println("Usage: konfetti explain <file> OR cat file | konfetti explain")
		return nil
	}
	file := c.Args().First()
	fmt.Printf("🔎 Explaining %s ... (mock)\n", file)
	time.Sleep(1 * time.Second)
	fmt.Println("✨ This config seems to define a JSON app layout. Debug enabled. (mock data)")
	return nil
}

// ---------------- Helper Functions ----------------
func getDefaultScanPaths() []string {
	homeDir, _ := os.UserHomeDir()

	if runtime.GOOS == "windows" {
		return []string{
			os.Getenv("ProgramData"),
			os.Getenv("APPDATA"),
			os.Getenv("LOCALAPPDATA"),
		}
	}

	paths := []string{"/etc"}
	if runtime.GOOS == "darwin" {
		paths = append(paths, filepath.Join(homeDir, "Library/Application Support"))
	}
	paths = append(paths, filepath.Join(homeDir, ".config"))

	return paths
}

func printTable(results []ConfigResult) {
	fmt.Printf("%-40s | %-20s | %-30s | %-8s\n", "File", "Setting", "Value", "Format")
	fmt.Println(strings.Repeat("-", 110))
	for _, r := range results {
		for k, v := range r.Settings {
			fmt.Printf("%-40s | %-20s | %-30s | %-8s\n", r.File, k, fmt.Sprintf("%v", v), r.Format)
		}
	}
}

// ---------------- Scan & Filter ----------------
func ScanAndFilter(paths []string, exts []string, filterName, filterKey, filterValue string) ([]ConfigResult, []string) {
	files, errors := scanner.ScanDirs(paths, exts)
	var results []ConfigResult

	for _, f := range files {
		result, format := parser.ParseFile(f)
		if filterKey != "" || filterValue != "" {
			filteredResult := make(map[string]interface{})
			for k, v := range result {
				keyMatch := filterKey == "" || strings.Contains(strings.ToLower(k), strings.ToLower(filterKey))
				valMatch := filterValue == "" || strings.Contains(strings.ToLower(fmt.Sprintf("%v", v)), strings.ToLower(filterValue))
				if keyMatch && valMatch {
					filteredResult[k] = v
				}
			}
			result = filteredResult
		}

		if len(result) > 0 && (filterName == "" || strings.Contains(f, filterName)) {
			results = append(results, ConfigResult{
				File:     f,
				Format:   format,
				Settings: result,
			})
		}
	}
	return results, errors
}
