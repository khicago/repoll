// Package main implements repoll, a Git repository management tool that allows
// users to clone, update, and warm-up multiple repositories based on TOML configuration.
//
// repoll supports two main operations:
// 1. make - Clone/update repositories from configuration
// 2. mkconf - Generate configuration from existing Git repositories
//
// The tool provides concurrent operations, intelligent warm-up for different
// project types (Go, Node.js, etc.), and detailed reporting capabilities.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/bagaking/goulp/wlog"
	"github.com/spf13/cobra"
)

// rootCmd is the main command for repoll CLI application.
// It serves as the parent for all subcommands.
var rootCmd = &cobra.Command{
	Use:   "repoll",
	Short: "The Ultimate Git Repository Management Tool",
	Long: `repoll (Repository Puller) is a lightning-fast, developer-friendly CLI tool 
that revolutionizes how you manage multiple Git repositories. Whether you're working 
with microservices, managing open-source contributions, or handling complex multi-repo 
projects, repoll makes it effortless.

Features:
- Lightning fast concurrent operations
- Smart warm-up for Go, Node.js projects
- Simple TOML configuration
- Flexible workflows and rich reporting`,
}

// cmdMake is the subcommand for processing repository configurations.
// It clones or updates repositories based on TOML configuration files.
var cmdMake = &cobra.Command{
	Use:   "make [paths to the TOML config file]",
	Short: "Clone or update repositories based on TOML configuration",
	Long: `The make command processes one or more TOML configuration files and performs
the following operations for each repository:

1. Clone the repository if it doesn't exist locally
2. Update the repository if it already exists (git pull)
3. Run warm-up operations (go mod download, npm install, etc.) if enabled

The command supports concurrent processing and provides detailed progress feedback.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize report to track all operations
		report := MakeReport{Actions: make([]*MakeAction, 0)}

		// Process each configuration file provided as argument
		for _, path := range args {
			// Convert to absolute path for consistent handling
			configPath, err := filepath.Abs(path)
			if err != nil {
				wlog.Common().Errorf("Error determining absolute path for %s: %s\n", path, err)
				continue
			}

			// Process the configuration file and update report
			if err := processConfig(configPath, &report); err != nil {
				wlog.Common().Errorf("Error processing config file %s: %s\n", configPath, err)
				continue
			}
		}

		// Generate detailed report if requested
		if reportFlag, _ := cmd.Flags().GetBool("report"); reportFlag {
			reportFileName := time.Now().Format("20060102-150405") + "_make_report.log"
			if err := os.WriteFile(reportFileName, []byte(report.Report()), os.ModePerm); err != nil {
				wlog.Common().Errorf("Failed to write report %s: %s\n", reportFileName, err)
			} else {
				wlog.Common().Infof("Report saved to %s\n", reportFileName)
			}
		}
	},
}

// cmdMakeConf is the subcommand for generating configuration files.
// It scans directories for Git repositories and creates TOML configuration.
var cmdMakeConf = &cobra.Command{
	Use:   "mkconf [directory]",
	Short: "Generate TOML configuration from existing Git repositories",
	Long: `The mkconf command scans the specified directory for Git repositories and 
generates a TOML configuration file that can be used with the 'make' command.

The generated configuration includes:
- Repository URLs and paths
- Status information (uncommitted changes, unmerged commits)
- Organized site groupings based on remote origins

This is useful for:
- Migrating existing repository collections to repoll
- Creating backup configurations of current setups
- Sharing repository collections with team members`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize report to track discovery operations
		report := MkconfReport{Actions: make([]*MkconfAction, 0)}

		// Scan directory and generate configuration
		if err := makeConfig(args[0], &report); err != nil {
			wlog.Common().Errorf("Error generating config for directory %s: %s\n", args[0], err)
			return
		}

		// Generate detailed report if requested
		if reportFlag, _ := cmd.Flags().GetBool("report"); reportFlag {
			reportFileName := time.Now().Format("20060102-150405") + "_mkconf_report.log"
			if err := os.WriteFile(reportFileName, []byte(report.Report()), os.ModePerm); err != nil {
				wlog.Common().Errorf("Failed to write report %s: %s\n", reportFileName, err)
			} else {
				wlog.Common().Infof("Report saved to %s\n", reportFileName)
			}
		}
	},
}

// main is the entry point of the repoll application.
// It executes the root command and handles any execution errors.
func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Error executing command: %s\n", err)
		os.Exit(1)
	}
}

// init initializes the CLI application by setting up flags and subcommands.
// This function is called automatically before main().
func init() {
	// Add flags for detailed reporting
	cmdMake.Flags().Bool("report", false, "Generate a detailed report after command execution")
	cmdMakeConf.Flags().Bool("report", false, "Generate a detailed report after command execution")

	// Register subcommands with the root command
	rootCmd.AddCommand(cmdMake)
	rootCmd.AddCommand(cmdMakeConf)
}
