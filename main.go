package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/bagaking/goulp/wlog"
	"github.com/spf13/cobra"
)

// 主命令和根命令
var rootCmd = &cobra.Command{
	Use: "repoll",
}

// make 子命令
var cmdMake = &cobra.Command{
	Use:   "make [paths to the TOML config file]",
	Short: "Repoll clones or updates repositories based on the TOML configuration.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Usage: repoll [path to the TOML config file]")
		}

		report := MakeReport{Actions: make([]*MakeAction, 0)}
		for _, path := range args {
			configPath, err := filepath.Abs(path)
			if err != nil {
				wlog.Common().Errorf("Error determining absolute path: %s\n", err)
			}

			// 假设已经有了一个 processConfig 函数
			if err := processConfig(configPath, &report); err != nil {
				wlog.Common().Errorf("Error processing config file: %s\n", err)
			}
		}

		if reportFlag, _ := cmd.Flags().GetBool("report"); reportFlag {
			reportFileName := time.Now().Format("20060102-150405") + "_make_report.log"
			if err := os.WriteFile(reportFileName, []byte(report.Report()), os.ModePerm); err != nil {
				wlog.Common().Errorf("Write report %s failed, err= %s\n", reportFileName, err)
			}
		}
	},
}

// mkconf 子命令
var cmdMakeConf = &cobra.Command{
	Use:   "mkconf [directory]",
	Short: "Mkconf scans the given directory for git repositories and creates a TOML config for repoll.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Usage: mkconf [directory]")
		}
		report := MkconfReport{Actions: make([]*MkconfAction, 0)}
		if err := makeConfig(args[0], &report); err != nil {
			wlog.Common().Errorf("Error making config for %s: %s\n", args[0], err)
		}
		if reportFlag, _ := cmd.Flags().GetBool("report"); reportFlag {
			reportFileName := time.Now().Format("20060102-150405") + "_mkconf_report.log"
			if err := os.WriteFile(reportFileName, []byte(report.Report()), os.ModePerm); err != nil {
				wlog.Common().Errorf("Write report %s failed, err= %s\n", reportFileName, err)
			}
		}
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

func init() {
	cmdMake.Flags().Bool("report", false, "Generate a detailed report after command execution.")
	cmdMakeConf.Flags().Bool("report", false, "Generate a detailed report after command execution.")

	rootCmd.AddCommand(cmdMake)
	rootCmd.AddCommand(cmdMakeConf)
}
