package issue

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
	"github.com/suzuki-shunsuke/tfaction-go/pkg/config"
	createdriftissues "github.com/suzuki-shunsuke/tfaction-go/pkg/controller/create-drift-issues"
	"github.com/suzuki-shunsuke/tfaction-go/pkg/github"
	"github.com/suzuki-shunsuke/tfaction-go/pkg/util"
)

type Param struct {
	RepoOwner       string
	RepoName        string
	Target          string
	GitHubServerURL string
}

func (c *Controller) Run(ctx context.Context, logE *logrus.Entry, param *Param) error { //nolint:cyclop
	// Get or create a drift issue
	cfg, err := config.Read(c.fs)
	if err != nil {
		return fmt.Errorf("read tfaction-root.yaml: %w", err)
	}
	if cfg.DriftDetection == nil {
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

	var wgCfg *config.WorkingDirectory
	var targetGroup *config.TargetGroup
	for _, t := range cfg.TargetGroups {
		if !strings.HasPrefix(param.Target, t.Target) {
			continue
		}
		targetGroup = t
		wd := strings.Replace(param.Target, targetGroup.Target, targetGroup.WorkingDirectory, 1)
		p := filepath.Join(wd, cfg.WorkingDirectoryFile)
		w, err := config.ReadWorkingDirectory(c.fs, p)
		if err != nil {
			return fmt.Errorf("read %s: %w", p, err)
		}
		wgCfg = w
		break
	}
	if wgCfg == nil {
		return nil
	}

	if !createdriftissues.CheckEnabled(cfg, targetGroup, wgCfg) {
		logE.Info("drifit detection is disabled")
		return nil
	}

	// Find a drift issue from target
	issue, err := c.gh.GetIssue(ctx, repoOwner, repoName, fmt.Sprintf(`Terraform Drift (%s)`, param.Target))
	if err != nil {
		return fmt.Errorf("get a drift issue: %w", err)
	}
	if issue == nil {
		return c.createIssue(ctx, logE, repoOwner, repoName, param)
	}

	c.action.SetEnv("TFACTION_DRIFT_ISSUE_NUMBER", strconv.Itoa(issue.Number))
	c.action.SetEnv("TFACTION_DRIFT_ISSUE_STATE", issue.State)

	issueURL := fmt.Sprintf("%s/%s/%s/pull/%v", param.GitHubServerURL, repoOwner, repoName, issue.Number)
	c.action.Infof(issueURL)
	c.action.AddStepSummary("Drift Issue: " + issueURL)

	return nil
}

//go:generate mockery --name Action --testonly=false
type Action interface {
	AddStepSummary(markdown string)
	Infof(msg string, args ...interface{})
	SetEnv(k, v string)
	SetOutput(k, v string)
}

func (c *Controller) createIssue(ctx context.Context, logE *logrus.Entry, repoOwner, repoName string, param *Param) error {
	// Create a drift issue
	issue, err := c.gh.CreateIssue(ctx, repoOwner, repoName, &github.IssueRequest{
		Title: util.StrP(fmt.Sprintf(`Terraform Drift (%s)`, param.Target)),
		Body:  util.StrP(createdriftissues.IssueBodyTemplate),
	})
	if err != nil {
		logerr.WithError(logE, err).Error("create an issue")
	}
	logE.Info("created an issue")

	c.action.SetEnv("TFACTION_DRIFT_ISSUE_NUMBER", strconv.Itoa(issue.GetNumber()))
	c.action.SetEnv("TFACTION_DRIFT_ISSUE_STATE", "open")

	issueURL := fmt.Sprintf("%s/%s/%s/pull/%v", param.GitHubServerURL, repoOwner, repoName, issue.GetNumber())
	c.action.Infof(issueURL)
	c.action.AddStepSummary("Drift Issue: " + issueURL)

	return nil
}
