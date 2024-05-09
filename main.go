package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: repoll [path to the TOML config file]")
		os.Exit(1)
	}

	configPath, err := filepath.Abs(os.Args[1])
	if err != nil {
		fmt.Printf("Error determining absolute path: %s\n", err)
		os.Exit(1)
	}

	config, err := readConfig(configPath)
	if err != nil {
		fmt.Printf("Error reading config file: %s\n", err)
		os.Exit(1)
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
}
