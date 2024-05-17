package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/khicago/irr"
)

func getRepoNameFromURL(repoURL string) (ret string) {
	return strings.TrimPrefix(strings.TrimSuffix(repoURL, ".git"), getRemotePrefix(repoURL))
}

func getRemotePrefix(repoURL string) string {
	if strings.HasPrefix(repoURL, "http") || strings.HasPrefix(repoURL, "ssh") {
		// 解析HTTPS URL
		parsedURL, err := url.Parse(repoURL)
		if err != nil {
			fmt.Printf("error parsing URL: %s\n", err)
			return ""
		}
		// 重建URL，仅包含到域名部分
		parsedURL.Path = ""
		parsedURL.RawQuery = ""
		parsedURL.Fragment = ""
		return parsedURL.String() + "/"
	} else if strings.HasPrefix(repoURL, "git@") {
		// SSH格式URL，提取用户名和主机名
		colonIndex := strings.Index(repoURL, ":")
		if colonIndex != -1 {
			return repoURL[:colonIndex+1]
		}
	}
	return repoURL
}

func mkStatusMemo(str string, need bool, memo string) string {
	if !need {
		return str
	}
	if str != "" {
		str += " & "
	}
	str = memo
	return str
}

func makeConfig(dir string, outReport *MkconfReport) error {
	dirAbs, err := filepath.Abs(dir)
	if err != nil {
		return irr.Wrap(err, "wrong dir")
	}
	// The new configuration struct
	config := Config{Sites: []SiteConfig{}}

	remoteTable := make(map[string]SiteConfig)

	// Walk the directory to find git repositories
	err = filepath.Walk(dirAbs, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			return nil
		}

		repoPath, err := filepath.Abs(path)
		if err != nil {
			return err
		}

		// Check if this is a .git directory
		if _, err = os.Stat(filepath.Join(repoPath, ".git")); os.IsNotExist(err) {
			return nil
		}

		ma := &MkconfAction{
			Time:   time.Now(),
			Origin: path,
		}
		outReport.Actions = append(outReport.Actions, ma)

		origin, err := getGitRemoteOrigin(repoPath)
		if err != nil {
			fmt.Println(irr.Wrap(err, "get git remote origin failed, path= %s", repoPath))
		} else {
			ma.HasOrigin = true
			ma.Origin = origin
		}
		uncommitted, err := hasUncommittedChanges(repoPath)
		if err != nil {
			fmt.Println(irr.Wrap(err, "get git uncommitted changes failed, path= %s", repoPath))
			ma.Uncommitted = true
		}
		unmerged, err := hasUnmergedCommits(repoPath)
		if err != nil {
			fmt.Println(irr.Wrap(err, "get git unmerged changes failed, path= %s", repoPath))
			ma.Unmerged = true
		}

		statusMemo := mkStatusMemo("", uncommitted, "uncommitted")
		statusMemo = mkStatusMemo("", unmerged, "unmerged")

		remotePrefix := getRemotePrefix(origin)
		repo := Repo{Repo: getRepoNameFromURL(origin), Memo: statusMemo}
		relatedLocalPath := filepath.Join(dir, strings.TrimPrefix(
			strings.TrimSuffix(
				strings.TrimSuffix(repoPath, repo.Repo),
				filepath.Base(repo.Repo),
			),
			dirAbs,
		))
		fmt.Println("origin", origin)
		uniqueKey := remotePrefix + " @@ " + relatedLocalPath
		if r, ok := remoteTable[uniqueKey]; ok {
			r.Repos = append(r.Repos, repo)
			remoteTable[uniqueKey] = r
		} else {
			siteConfig := SiteConfig{
				// Assume remote prefix is the same as the origin but without the last part
				RemotePrefix: remotePrefix,
				Dir:          relatedLocalPath,
				Repos:        []Repo{repo},
			}
			remoteTable[uniqueKey] = siteConfig
		}
		// Skip checking subdirectories since we found a .git directory
		return filepath.SkipDir
	})

	for remote, siteConfig := range remoteTable {
		config.Sites = append(config.Sites, siteConfig)
		fmt.Printf("create site config for remote `%s`.\n", remote)
	}

	if err != nil {
		return err
	}

	// 创建输出文件
	file, err := os.Create(time.Now().Format("20060102-150405") + "_conf.toml")
	if err != nil {
		return err
	}
	defer file.Close()

	// 使用 toml.Encode 将 config 写入文件
	encoder := toml.NewEncoder(file)
	if err = encoder.Encode(config); err != nil {
		return err
	}

	return nil
}

func getGitRemoteOrigin(repoPath string) (string, error) {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	cmd.Dir = repoPath
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	str := strings.TrimSpace(string(out))
	if str == "" {
		return "", errors.New("git remote get-url origin not found")
	}
	return strings.TrimSpace(str), nil
}

func hasUncommittedChanges(repoPath string) (bool, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = repoPath
	out, err := cmd.CombinedOutput()
	if err != nil {
		return false, err
	}
	return bytes.TrimSpace(out) != nil, nil
}

func hasUnmergedCommits(repoPath string) (bool, error) {
	cmd := exec.Command("git", "cherry", "-v")
	cmd.Dir = repoPath
	out, err := cmd.CombinedOutput()
	if err != nil {
		return false, err
	}
	return bytes.TrimSpace(out) != nil, nil
}
