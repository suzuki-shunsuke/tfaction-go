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
	"github.com/suzuki-shunsuke/tfaction-go/pkg/github"
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

type Param struct {
	RepoOwner string
	RepoName  string
	PWD       string
}

func (ctrl *Controller) Run(ctx context.Context, logE *logrus.Entry, param *Param) error { //nolint:funlen,cyclop
	cfg, err := config.Read(ctrl.fs)
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

	workingDirectories, err := createdriftissues.ListWorkingDirectories(ctrl.fs, cfg, param.PWD)
	if err != nil {
		return fmt.Errorf("list working directories: %w", err)
	}

	targets := map[string]string{}
	for workingDirectoryPath, workingDirectory := range workingDirectories {
		targetGroup := createdriftissues.GetTargetGroupByWorkingDirectory(cfg.TargetGroups, workingDirectoryPath)
		if targetGroup == nil {
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

	deadline := time.Now().In(time.FixedZone("UTC", 0)).Add(-time.Duration(cfg.DriftDetection.Duration) * time.Hour).Format("2006-01-02T15:04:05+00:00")
	logE.WithFields(logrus.Fields{
		"duration": cfg.DriftDetection.Duration,
		"deadline": deadline,
	}).Info("check a deadline")

	issues, err := ctrl.gh.ListLeastRecentlyUpdatedIssues(ctx, param.RepoOwner, param.RepoName, cfg.DriftDetection.NumOfIssues, deadline)
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
		if _, err := ctrl.gh.ArchiveIssue(ctx, repoOwner, repoName, issue.Number, fmt.Sprintf(`Archived %s`, issue.Title)); err != nil {
			logE.WithError(err).Error("archive an issue")
		}
		logE.Info("archive an issue")
	}

	if len(arr) == 0 {
		githubactions.SetOutput("has_issues", "false")
		githubactions.SetOutput("issues", "[]")
		return nil
	}

	githubactions.SetOutput("has_issues", "true")
	b, err := json.Marshal(arr)
	if err != nil {
		return fmt.Errorf("marshal issues as JSON: %w", err)
	}
	githubactions.SetOutput("issues", string(b))
	return nil
}
