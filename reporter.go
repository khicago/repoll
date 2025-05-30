// report.go
package main

import (
	"fmt"
	"time"
)

type MakeAction struct {
	Time       time.Time
	Repository string
	Memo       string
	Duration   time.Duration
	Success    bool
	Error      string
}

type MkconfAction struct {
	Time        time.Time
	Path        string
	Origin      string
	HasOrigin   bool
	Uncommitted bool
	Unmerged    bool
}

type MakeReport struct {
	Actions []*MakeAction
}

type MkconfReport struct {
	Actions []*MkconfAction
}

func (mr *MakeReport) Report() string {
	ret := ""
	for _, action := range mr.Actions {
		ret += fmt.Sprintf("%v | Repo: %s | Duration: %v | Success: %v | Error: %s | Memo: %s\n",
			action.Time, action.Repository, action.Duration, action.Success, action.Error, action.Memo)
	}
	return ret
}

func (mr *MkconfReport) Report() string {
	ret := ""
	for _, action := range mr.Actions {
		ret += fmt.Sprintf("%v | Repo: %s | HasOrigin: %v | Uncommitted: %v | Unmerged: %v\n",
			action.Time, action.Origin, action.HasOrigin, action.Uncommitted, action.Unmerged)
	}
	return ret
}

func (ma *MakeAction) String() string {
	return fmt.Sprintf("%v | Repo: %s | Duration: %v | Success: %v | Error: %s | Memo: %s",
		ma.Time, ma.Repository, ma.Duration, ma.Success, ma.Error, ma.Memo)
}

func (ma *MkconfAction) String() string {
	return fmt.Sprintf("%v | Path: %s | Origin: %s | HasOrigin: %v | Uncommitted: %v | Unmerged: %v",
		ma.Time, ma.Path, ma.Origin, ma.HasOrigin, ma.Uncommitted, ma.Unmerged)
}
