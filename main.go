package main

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
)

// 主命令和根命令
var rootCmd = &cobra.Command{
	Use:   "repoll [path to the TOML config file]",
	Short: "Repoll clones or updates repositories based on the TOML configuration.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Usage: repoll [path to the TOML config file]")
		}

		configPath, err := filepath.Abs(args[0])
		if err != nil {
			fmt.Printf("Error determining absolute path: %s\n", err)
		}

		// 假设已经有了一个 processConfig 函数
		if err := processConfig(configPath); err != nil {
			fmt.Printf("Error processing config file: %s\n", err)
		}
	},
}

// mkconf 子命令
var mkconfCmd = &cobra.Command{
	Use:   "mkconf [directory]",
	Short: "Mkconf scans the given directory for git repositories and creates a TOML config for repoll.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Usage: mkconf [directory]")
		}
		if err := makeConfig(args[0]); err != nil {
			fmt.Printf("Error making config for %s: %s\n", args[0], err)
		}
	},
}

func processConfig(configPath string) (err error) {
	configPath, err = filepath.Abs(configPath)
	if err != nil {
		fmt.Printf("Error determining absolute path: %s\n", err)
		return err
	}

	config, err := readConfig(configPath)
	if err != nil {
		fmt.Printf("Error reading config file: %s\n", err)
		return err
	}

	for _, site := range config.Sites {
		for _, repo := range site.Repos {
			if err = gitCloneOrUpdate(repo, site); err != nil {
				fmt.Printf("Processing repository %s Failed, err= %s .\n", repo.Repo, err)
				continue
			} else {
				fmt.Printf("Processing repository %s success.\n", repo.Repo)
			}

			if repo.WarmUp || site.WarmUpAll {
				if err = warmUpRepo(repo.FullPath(site)); err != nil {
					fmt.Printf("- warm-up operations for repo %s failed, err= %s\n", repo.Repo, err)
				} else {
					fmt.Printf("- warm-up operations for repo %s success.\n", repo.Repo)
				}
			}
		}
	}

	fmt.Println("Repository processing complete.")
	return nil
}

func init() {
	rootCmd.AddCommand(mkconfCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
