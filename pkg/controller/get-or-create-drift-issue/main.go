package issue

import (
	"context"
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
	"github.com/suzuki-shunsuke/tfaction-go/pkg/config"
	createdriftissues "github.com/suzuki-shunsuke/tfaction-go/pkg/controller/create-drift-issues"
	"github.com/suzuki-shunsuke/tfaction-go/pkg/github"
	"github.com/suzuki-shunsuke/tfaction-go/pkg/util"
)

type Controller struct {
	gh     github.Client
	fs     afero.Fs
	action Action
}

func New(gh github.Client, fs afero.Fs, action Action) *Controller {
	return &Controller{
		gh:     gh,
		fs:     fs,
		action: action,
	}
}

type Param struct {
	RepoOwner       string
	RepoName        string
	Target          string
	GitHubServerURL string
}

func (ctrl *Controller) Run(ctx context.Context, logE *logrus.Entry, param *Param) error {
	// Get or create a drift issue
	cfg, err := config.Read(ctrl.fs)
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

	// Find a drift issue from target
	issue, err := ctrl.gh.GetIssue(ctx, repoOwner, repoName, fmt.Sprintf(`Terraform Drift (%s)`, param.Target))
	if err != nil {
		return fmt.Errorf("get a drift issue: %w", err)
	}
	if issue == nil {
		return ctrl.createIssue(ctx, logE, repoOwner, repoName, param)
	}

	ctrl.action.SetEnv("TFACTION_DRIFT_ISSUE_NUMBER", strconv.Itoa(issue.Number))
	ctrl.action.SetEnv("TFACTION_DRIFT_ISSUE_STATE", issue.State)

	issueURL := fmt.Sprintf("%s/%s/%s/pull/%v", param.GitHubServerURL, repoOwner, repoName, issue.Number)
	ctrl.action.Infof(issueURL)
	ctrl.action.AddStepSummary(fmt.Sprintf("Drift Issue: %s", issueURL))

	return nil
}

//go:generate mockery --name Action --testonly=false
type Action interface {
	AddStepSummary(markdown string)
	Infof(msg string, args ...interface{})
	SetEnv(k, v string)
	SetOutput(k, v string)
}

func (ctrl *Controller) createIssue(ctx context.Context, logE *logrus.Entry, repoOwner, repoName string, param *Param) error {
	// Create a drift issue
	issue, err := ctrl.gh.CreateIssue(ctx, repoOwner, repoName, &github.IssueRequest{
		Title: util.StrP(fmt.Sprintf(`Terraform Drift (%s)`, param.Target)),
		Body:  util.StrP(createdriftissues.IssueBodyTemplate),
	})
	if err != nil {
		logerr.WithError(logE, err).Error("create an issue")
	}
	logE.Info("created an issue")

	ctrl.action.SetEnv("TFACTION_DRIFT_ISSUE_NUMBER", strconv.Itoa(issue.GetNumber()))
	ctrl.action.SetEnv("TFACTION_DRIFT_ISSUE_STATE", "open")

	issueURL := fmt.Sprintf("%s/%s/%s/pull/%v", param.GitHubServerURL, repoOwner, repoName, *issue.Number)
	ctrl.action.Infof(issueURL)
	ctrl.action.AddStepSummary(fmt.Sprintf("Drift Issue: %s", issueURL))

	return nil
}
