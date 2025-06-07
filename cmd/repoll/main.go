package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/khicago/repoll/internal/cli"
	"github.com/khicago/repoll/internal/config"
	"github.com/khicago/repoll/internal/process"
	"github.com/khicago/repoll/internal/reporter"
	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "2025-06-02"
)

// Global flags
var (
	reportFlag  bool
	dryRunFlag  bool
	verboseFlag bool
	quietFlag   bool
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "repoll",
		Short: "A powerful tool for managing multiple Git repositories",
		Long: `repoll (Repository Puller) is a lightning-fast, developer-friendly CLI tool 
that revolutionizes how you manage multiple Git repositories.`,
		Version: fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date),
	}

	// Add global flags
	rootCmd.PersistentFlags().BoolVar(&reportFlag, "report", false, "Generate detailed execution report")
	rootCmd.PersistentFlags().BoolVarP(&dryRunFlag, "dry-run", "n", false, "Show what would be done without executing")
	rootCmd.PersistentFlags().BoolVarP(&verboseFlag, "verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().BoolVarP(&quietFlag, "quiet", "q", false, "Enable quiet mode (minimal output)")

	// Main command for processing config files
	runCmd := &cobra.Command{
		Use:   "run [config-files...]",
		Short: "Process configuration files to clone/update repositories",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runProcessConfigs(args)
		},
	}

	// mkconf command
	mkconfCmd := &cobra.Command{
		Use:   "mkconf [directory]",
		Short: "Generate configuration file from existing Git repositories",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			targetDir := "."
			if len(args) > 0 {
				targetDir = args[0]
			}
			return runMakeConfig(targetDir)
		},
	}

	// Legacy support: direct config file processing (no subcommand)
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}
		
		// Check if the first argument is a config file
		if len(args) > 0 && strings.HasSuffix(args[0], ".toml") {
			err := runProcessConfigs(args)
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		
		cmd.Help()
	}

	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(mkconfCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// runProcessConfigs processes multiple configuration files
func runProcessConfigs(configPaths []string) error {
	ui := cli.NewUIManager(quietFlag, verboseFlag)
	
	if !quietFlag {
		ui.Banner()
	}
	
	startTime := time.Now()
	totalRepos := 0
	successCount := 0
	failCount := 0
	
	var report *reporter.MakeReport
	if reportFlag {
		report = &reporter.MakeReport{}
	}
	
	opts := &process.ProcessorOptions{
		UI:     ui,
		DryRun: dryRunFlag,
	}
	
	if dryRunFlag {
		ui.Warning("DRY RUN MODE: No actual changes will be made")
	}
	
	for i, configPath := range configPaths {
		if len(configPaths) > 1 {
			ui.Section(fmt.Sprintf("Processing %s (%d/%d)", configPath, i+1, len(configPaths)))
		}
		
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			ui.Error("Configuration file not found: %s", configPath)
			failCount++
			continue
		}
		
		err := process.ProcessConfig(configPath, report, opts)
		if err != nil {
			ui.Error("Failed to process %s: %v", configPath, err)
			failCount++
			continue
		}
		
		// Count repositories from this config
		if cfg, err := config.ReadFromFile(configPath); err == nil {
			for _, site := range cfg.Sites {
				totalRepos += len(site.Repos)
			}
		}
	}
	
	// Calculate success/fail counts from report if available
	if report != nil {
		successCount = 0
		failCount = 0
		for _, action := range report.Actions {
			if action.Success {
				successCount++
			} else {
				failCount++
			}
		}
		totalRepos = len(report.Actions)
	}
	
	totalDuration := time.Since(startTime)
	
	if !dryRunFlag {
		ui.Summary(totalRepos, successCount, failCount, totalDuration)
	} else {
		ui.Info("Dry run completed in %v", totalDuration)
	}
	
	// Generate report if requested
	if reportFlag && report != nil {
		ui.Info("Generating execution report...")
		fmt.Println(report.Report())
	}
	
	return nil
}

// runMakeConfig generates a configuration file from existing repositories
func runMakeConfig(targetDir string) error {
	ui := cli.NewUIManager(quietFlag, verboseFlag)
	
	if !quietFlag {
		ui.Banner()
	}
	
	startTime := time.Now()
	
	ui.Info("Scanning directory: %s", targetDir)
	
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		return fmt.Errorf("directory does not exist: %s", targetDir)
	}
	
	var report *reporter.MkconfReport
	if reportFlag {
		report = &reporter.MkconfReport{}
	}
	
	if dryRunFlag {
		ui.Warning("DRY RUN MODE: Configuration will be printed to stdout only")
	}
	
	cfg, err := config.GenerateFromDirectory(targetDir, report)
	if err != nil {
		return fmt.Errorf("failed to generate configuration: %w", err)
	}
	
	if cfg == nil || len(cfg.Sites) == 0 {
		ui.Warning("No Git repositories found in %s", targetDir)
		return nil
	}
	
	// Convert config to TOML string
	configContent, err := config.ToTOML(cfg)
	if err != nil {
		return fmt.Errorf("failed to convert configuration to TOML: %w", err)
	}
	
	if dryRunFlag {
		ui.Section("Generated Configuration (DRY RUN)")
		fmt.Println(configContent)
	} else {
		outputFile := "repos.toml"
		err = os.WriteFile(outputFile, []byte(configContent), 0644)
		if err != nil {
			return fmt.Errorf("failed to write configuration file: %w", err)
		}
		ui.Success("Configuration written to %s", outputFile)
	}
	
	duration := time.Since(startTime)
	ui.Info("Configuration generation completed in %v", duration)
	
	// Generate report if requested
	if reportFlag && report != nil {
		ui.Info("Generating discovery report...")
		fmt.Println(report.Report())
	}
	
	return nil
} 