package issues

import (
	"context"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
	"github.com/suzuki-shunsuke/tfaction-go/pkg/config"
	"github.com/suzuki-shunsuke/tfaction-go/pkg/github"
	"github.com/suzuki-shunsuke/tfaction-go/pkg/util"
)

type Controller struct {
	gh github.Client
	fs afero.Fs
}

func New(gh github.Client, fs afero.Fs) *Controller {
	return &Controller{
		gh: gh,
		fs: fs,
	}
}

func (c *Controller) Run(ctx context.Context, logE *logrus.Entry, param *Param) error { //nolint:cyclop,funlen
	cfg, err := config.Read(c.fs)
	if err != nil {
		return fmt.Errorf("read tfaction-root.yaml: %w", err)
	}
	if cfg.DriftDetection == nil {
		logE.Info("drift_detection is null")
		return nil
	}

	repoOwner := param.RepoOwner
	repoName := param.RepoName
	if cfg.DriftDetection != nil {
		if cfg.DriftDetection.IssueRepoOwner != "" {
			repoOwner = cfg.DriftDetection.IssueRepoOwner
		}
		if cfg.DriftDetection.IssueRepoName != "" {
			repoName = cfg.DriftDetection.IssueRepoName
		}
	}

	workingDirectories, err := ListWorkingDirectories(c.fs, cfg, param.PWD)
	if err != nil {
		return err
	}

	logE.WithField("num_of_working_dirs", len(workingDirectories)).Debug("search working directories")
	targets := ListTargets(cfg, workingDirectories)
	logE.WithField("num_of_targets", len(targets)).Debug("convert working directories to targets")

	// Search GitHub Issues
	issues, err := c.gh.ListIssues(ctx, repoOwner, repoName)
	if err != nil {
		return fmt.Errorf("list issues: %w", err)
	}
	logE.WithField("num_of_issues", len(issues)).Debug("search GiHub issues")

	issueTargets := make(map[string]struct{}, len(issues))
	issueMap := make(map[string]*github.Issue, len(issues))
	for _, issue := range issues {
		issueTargets[issue.Target] = struct{}{}
		issueMap[issue.Target] = issue
	}

	for target := range targets {
		// Create and close issues
		logE := logE.WithField("target", target)
		if _, ok := issueTargets[target]; ok {
			continue
		}
		issue, err := c.gh.CreateIssue(ctx, repoOwner, repoName, &github.IssueRequest{
			Title: util.StrP(fmt.Sprintf(`Terraform Drift (%s)`, target)),
			Body:  util.StrP(IssueBodyTemplate),
		})
		if err != nil {
			logerr.WithError(logE, err).Error("create an issue")
		}
		logE.Info("created an issue")
		if _, err := c.gh.CloseIssue(ctx, repoOwner, repoName, issue.GetNumber()); err != nil {
			logerr.WithError(logE, err).Error("close an issue")
		}
		logE.Debug("closed an issue")
		issueMap[target] = &github.Issue{
			Number: issue.GetNumber(),
			Title:  issue.GetTitle(),
		}
	}

	for target, issue := range issueMap {
		// Rename issues whose targets are not found
		if _, ok := targets[target]; ok {
			continue
		}
		logE := logE.WithFields(logrus.Fields{
			"target":       target,
			"issue_number": issue.Number,
		})
		if _, err := c.gh.ArchiveIssue(ctx, repoOwner, repoName, issue.Number, "Archived "+issue.Title); err != nil {
			logE.WithError(err).Error("archive an issue")
		}
		logE.Info("archive an issue")
	}

	return nil
}

type Param struct {
	RepoOwner string
	RepoName  string
	PWD       string
}

func ListWorkingDirectories(aferoFs afero.Fs, cfg *config.Config, pwd string) (map[string]*config.WorkingDirectory, error) {
	workingDirectoryPaths := map[string]*config.WorkingDirectory{}
	baseWorkingDirectory := filepath.Join(pwd, cfg.BaseWorkingDirectory)
	if err := fs.WalkDir(afero.NewIOFS(aferoFs), baseWorkingDirectory, func(p string, dirEntry fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if dirEntry.Name() != cfg.WorkingDirectoryFile {
			return nil
		}
		f, err := filepath.Rel(pwd, filepath.Dir(p))
		if err != nil {
			return fmt.Errorf("get a relative path of a working directory: %w", err)
		}

		wdCfg, err := config.ReadWorkingDirectory(aferoFs, p)
		if err != nil {
			return fmt.Errorf("read %s: %w", p, err)
		}
		workingDirectoryPaths[f] = wdCfg

		return nil
	}); err != nil {
		return nil, fmt.Errorf("search working directories: %w", err)
	}
	return workingDirectoryPaths, nil
}

func GetTargetGroupByWorkingDirectory(targetGroups []*config.TargetGroup, workingDirectoryPath string) *config.TargetGroup {
	for _, targetGroup := range targetGroups {
		if strings.HasPrefix(workingDirectoryPath, targetGroup.WorkingDirectory) {
			return targetGroup
		}
	}
	return nil
}

func GetTargetByWorkingDirectory(workingDirectoryPath string, targetGroup *config.TargetGroup) string {
	return strings.Replace(workingDirectoryPath, targetGroup.WorkingDirectory, targetGroup.Target, 1)
}

func ListTargets(cfg *config.Config, workingDirectories map[string]*config.WorkingDirectory) map[string]struct{} {
	targets := make(map[string]struct{}, len(workingDirectories))
	for workingDirectoryPath, wdCfg := range workingDirectories {
		if targetGroup := GetTargetGroupByWorkingDirectory(cfg.TargetGroups, workingDirectoryPath); targetGroup != nil {
			if CheckEnabled(cfg, targetGroup, wdCfg) {
				targets[GetTargetByWorkingDirectory(workingDirectoryPath, targetGroup)] = struct{}{}
			}
		}
	}
	return targets
}

func CheckEnabled(cfg *config.Config, targetGroup *config.TargetGroup, wdCfg *config.WorkingDirectory) bool {
	if wdCfg.DriftDetection != nil {
		if wdCfg.DriftDetection.Enabled != nil {
			return *wdCfg.DriftDetection.Enabled
		}
		return true
	}
	if targetGroup.DriftDetection != nil {
		if targetGroup.DriftDetection.Enabled != nil {
			return *targetGroup.DriftDetection.Enabled
		}
		return true
	}
	return cfg.DriftDetection.Enabled == nil || *cfg.DriftDetection.Enabled
}

const IssueBodyTemplate = `
This issus was created by [tfaction](https://suzuki-shunsuke.github.io/tfaction/docs/).

About this issue, please see [the document](https://suzuki-shunsuke.github.io/tfaction/docs/feature/drift-detection).
`
