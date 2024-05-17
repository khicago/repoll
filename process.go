package main

import (
	"path/filepath"
	"time"

	"github.com/bagaking/goulp/wlog"
)

func processConfig(configPath string, outReport *MakeReport) (err error) {
	configPath, err = filepath.Abs(configPath)
	if err != nil {
		wlog.Common().Errorf("Error determining absolute path: %s\n", err)
		return err
	}

	config, err := readConfig(configPath)
	if err != nil {
		wlog.Common().Errorf("Error reading config file: %s\n", err)
		return err
	}

	for _, site := range config.Sites {
		for _, repo := range site.Repos {
			ma := &MakeAction{
				Time:       time.Now(),
				Repository: site.RemotePrefix + "/" + repo.Repo,
				Success:    true,
				Memo:       repo.Memo,
			}
			outReport.Actions = append(outReport.Actions, ma)

			if err = gitCloneOrUpdate(repo, site); err != nil {
				wlog.Common().Errorf("Processing repository %s Failed, err= %s .\n", repo.Repo, err)
				ma.Success = false
				ma.Error = err.Error()
				continue
			}
			wlog.Common().Infof("Processing repository %s success.\n", repo.Repo)

			if repo.WarmUp || site.WarmUpAll {
				if err = warmUpRepo(repo.FullPath(site)); err != nil {
					wlog.Common().Errorf("- warm-up operations for repo %s failed, err= %s\n", repo.Repo, err)
				} else {
					wlog.Common().Infof("- warm-up operations for repo %s success.\n", repo.Repo)
				}
			}

			ma.Duration = time.Since(ma.Time)
		}
	}

	wlog.Common().Info("Origin processing complete.")
	return nil
}
