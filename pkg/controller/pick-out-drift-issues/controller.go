package issues

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sethvargo/go-githubactions"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/tfaction-go/pkg/config"
	createdriftissues "github.com/suzuki-shunsuke/tfaction-go/pkg/controller/create-drift-issues"
	issue "github.com/suzuki-shunsuke/tfaction-go/pkg/controller/get-or-create-drift-issue"
	"github.com/suzuki-shunsuke/tfaction-go/pkg/github"
)

type Controller struct {
	gh     github.Client
	fs     afero.Fs
	action issue.Action
}

func New(gh github.Client, fs afero.Fs, action issue.Action) *Controller {
	return &Controller{
		gh:     gh,
		fs:     fs,
		action: action,
	}
}

type Param struct {
	RepoOwner string
	RepoName  string
	PWD       string
	Now       time.Time
}

func (c *Controller) Run(ctx context.Context, logE *logrus.Entry, param *Param) error { //nolint:funlen,cyclop
	cfg, err := config.Read(c.fs)
	if err != nil {
		return fmt.Errorf("read tfaction-root.yaml: %w", err)
	}
	if cfg.DriftDetection == nil {
		githubactions.SetOutput("has_issues", "false")
		githubactions.SetOutput("issues", "[]")
		return nil
	}

	repoOwner := param.RepoOwner
	repoName := param.RepoName
	if cfg.DriftDetection.IssueRepoOwner != "" {
		repoOwner = cfg.DriftDetection.IssueRepoOwner
	}
	if cfg.DriftDetection.IssueRepoName != "" {
		repoName = cfg.DriftDetection.IssueRepoName
	}

	workingDirectories, err := createdriftissues.ListWorkingDirectories(c.fs, cfg, param.PWD)
	if err != nil {
		return fmt.Errorf("list working directories: %w", err)
	}

	targets := map[string]string{}
	for workingDirectoryPath, workingDirectory := range workingDirectories {
		targetGroup := createdriftissues.GetTargetGroupByWorkingDirectory(cfg.TargetGroups, workingDirectoryPath)
		if targetGroup == nil {
			continue
		}
		if !createdriftissues.CheckEnabled(cfg, targetGroup, workingDirectory) {
			continue
		}
		target := createdriftissues.GetTargetByWorkingDirectory(workingDirectoryPath, targetGroup)
		// Merge cfg and targetGroup and workingDirectory
		runsOn := "ubuntu-latest"
		for _, r := range []string{
			workingDirectory.TerraformPlanConfig.GetRunsOn(),
			workingDirectory.RunsOn,
			targetGroup.TerraformPlanConfig.GetRunsOn(),
			targetGroup.RunsOn,
			cfg.RunsOn,
		} {
			if r == "" {
				continue
			}
			runsOn = r
			break
		}
		targets[target] = runsOn
	}

	logE.WithField("num_of_working_dirs", len(workingDirectories)).Debug("search working directories")
	logE.WithField("num_of_targets", len(targets)).Debug("convert working directories to targets")

	deadline := getDeadline(param.Now, cfg.DriftDetection.MinimumDetectionInterval)
	logE.WithFields(logrus.Fields{
		"minimum_detection_interval": cfg.DriftDetection.MinimumDetectionInterval,
		"deadline":                   deadline,
	}).Info("check a deadline")

	issues, err := c.gh.ListLeastRecentlyUpdatedIssues(ctx, repoOwner, repoName, cfg.DriftDetection.NumOfIssues, deadline)
	if err != nil {
		return fmt.Errorf("list drift issues: %w", err)
	}
	logE.WithField("num_of_issues", len(issues)).Info("list least recently updated issues")

	arr := make([]*github.Issue, 0, len(issues))
	for _, issue := range issues {
		if runsOn, ok := targets[issue.Target]; ok {
			issue.RunsOn = runsOn
			arr = append(arr, issue)
			continue
		}
		logE := logE.WithFields(logrus.Fields{
			"issue_number": issue.Number,
			"target":       issue.Target,
		})
		if _, err := c.gh.ArchiveIssue(ctx, repoOwner, repoName, issue.Number, "Archived "+issue.Title); err != nil {
			logE.WithError(err).Error("archive an issue")
		}
		logE.Info("archive an issue")
	}

	return c.setOutput(arr)
}

func (c *Controller) setOutput(issues []*github.Issue) error {
	if len(issues) == 0 {
		c.action.SetOutput("has_issues", "false")
		c.action.SetOutput("issues", "[]")
		return nil
	}

	c.action.SetOutput("has_issues", "true")
	b, err := json.Marshal(issues)
	if err != nil {
		return fmt.Errorf("marshal issues as JSON: %w", err)
	}
	c.action.SetOutput("issues", string(b))
	return nil
}

func getDeadline(now time.Time, duration int) string {
	return now.In(time.FixedZone("UTC", 0)).Add(-time.Duration(duration) * time.Hour).Format("2006-01-02T15:04:05+00:00")
}
