package process

import (
	"fmt"
	"os"
	"time"

	"github.com/khicago/repoll/internal/cli"
	"github.com/khicago/repoll/internal/config"
	"github.com/khicago/repoll/internal/git"
	"github.com/khicago/repoll/internal/reporter"
	"github.com/khicago/repoll/internal/warmup"
)

// ProcessorOptions contains options for the processor
type ProcessorOptions struct {
	UI     *cli.UIManager
	DryRun bool
}

// ProcessConfig processes a configuration file and manages repositories
func ProcessConfig(configPath string, report *reporter.MakeReport, opts *ProcessorOptions) error {
	if opts == nil {
		opts = &ProcessorOptions{
			UI: cli.NewUIManager(false, false),
		}
	}
	
	opts.UI.Info("Loading configuration from %s", configPath)
	
	cfg, err := config.ReadFromFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}
	
	totalRepos := 0
	for _, site := range cfg.Sites {
		totalRepos += len(site.Repos)
	}
	
	opts.UI.Info("Found %d site(s) with %d repositories", len(cfg.Sites), totalRepos)
	
	var progressBar *cli.ProgressBar
	if totalRepos > 1 {
		progressBar = opts.UI.NewProgressBar(totalRepos, "Processing")
	}
	
	processedCount := 0
	for _, site := range cfg.Sites {
		opts.UI.Verbose("Processing site: %s", site.RemotePrefix)
		
		err := processSite(site, report, opts, progressBar, &processedCount)
		if err != nil {
			opts.UI.Error("Error processing site %s: %v", site.RemotePrefix, err)
			continue
		}
	}
	
	if progressBar != nil {
		progressBar.Finish()
	}

	return nil
}

// processSite processes all repositories in a site configuration
func processSite(site config.SiteConfig, report *reporter.MakeReport, opts *ProcessorOptions, progressBar *cli.ProgressBar, processedCount *int) error {
	for _, repo := range site.Repos {
		startTime := time.Now()
		*processedCount++
		
		if progressBar != nil {
			progressBar.Update(*processedCount)
		}
		
		action := &reporter.MakeAction{
			Time:       startTime,
			Repository: repo.DisplayName(),
			Memo:       repo.Memo,
		}

		// Determine the action to perform
		targetPath := repo.FullPath(site)
		actionName := "Cloning"
		
		if _, err := os.Stat(targetPath); err == nil {
			actionName = "Updating"
		}
		
		if opts.DryRun {
			opts.UI.DryRun("Would %s %s -> %s", actionName, repo.DisplayName(), targetPath)
			if shouldWarmUp(repo, site) {
				opts.UI.DryRun("Would warm up %s", targetPath)
			}
			continue
		}
		
		opts.UI.ProcessingRepo(repo.DisplayName(), actionName)
		
		err := processRepository(repo, site, opts)
		duration := time.Since(startTime)
		
		action.Duration = duration
		if err != nil {
			action.Success = false
			action.Error = err.Error()
			opts.UI.RepoResult(repo.DisplayName(), false, duration, err)
		} else {
			action.Success = true
			opts.UI.RepoResult(repo.DisplayName(), true, duration, nil)
		}

		if report != nil {
			report.Actions = append(report.Actions, action)
		}
	}

	return nil
}

// processRepository processes a single repository
func processRepository(repo config.Repo, site config.SiteConfig, opts *ProcessorOptions) error {
	targetPath := repo.FullPath(site)
	repoURL := repo.RepoUrl(site)

	// Check if repository already exists
	if _, err := os.Stat(targetPath); err == nil {
		// Repository exists, try to update it
		opts.UI.Verbose("Repository exists at %s, updating...", targetPath)
		err := git.Update(targetPath)
		if err != nil {
			return fmt.Errorf("failed to update repository: %w", err)
		}
	} else {
		// Repository doesn't exist, clone it
		opts.UI.Verbose("Cloning %s to %s", repoURL, targetPath)
		err := git.Clone(repoURL, targetPath)
		if err != nil {
			return fmt.Errorf("failed to clone repository: %w", err)
		}
	}

	// Perform warm-up if needed
	if shouldWarmUp(repo, site) {
		opts.UI.Verbose("Starting warm-up for %s", targetPath)
		err := warmup.Perform(targetPath)
		if err != nil {
			opts.UI.Warning("Warm-up failed for %s: %v", targetPath, err)
			// Don't return error for warm-up failures as they're not critical
		} else {
			opts.UI.Verbose("Warm-up completed for %s", targetPath)
		}
	}

	return nil
}

// shouldWarmUp determines if warm-up should be performed for a repository
func shouldWarmUp(repo config.Repo, site config.SiteConfig) bool {
	return repo.WarmUp || site.WarmUpAll
} 