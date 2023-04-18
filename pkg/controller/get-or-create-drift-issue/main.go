package issue

import (
	"context"
	"fmt"
	"strconv"

	"github.com/sethvargo/go-githubactions"
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

type Param struct {
	RepoOwner string
	RepoName  string
	Target    string
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
		// Create a drift issue
		issue, err := ctrl.gh.CreateIssue(ctx, repoOwner, repoName, &github.IssueRequest{
			Title: util.StrP(fmt.Sprintf(`Terraform Drift (%s)`, param.Target)),
			Body: util.StrP(`
This issus was created by [tfaction](https://suzuki-shunsuke.github.io/tfaction/docs/).

## :warning: Don't change the issue title

tfaction searches Issues by Issue title. So please don't change the issue title.
`),
		})
		if err != nil {
			logerr.WithError(logE, err).Error("create an issue")
		}
		logE.Info("created an issue")

		githubactions.SetOutput("issue_number", strconv.Itoa(issue.GetNumber()))
		githubactions.SetOutput("issue_state", "open")

		githubactions.SetEnv("TFACTION_DRIFT_ISSUE_NUMBER", strconv.Itoa(issue.GetNumber()))
		githubactions.SetEnv("TFACTION_DRIFT_ISSUE_STATE", "open")

		return nil
	}

	githubactions.SetOutput("issue_number", strconv.Itoa(issue.Number))
	githubactions.SetOutput("issue_state", issue.State)

	githubactions.SetEnv("TFACTION_DRIFT_ISSUE_NUMBER", strconv.Itoa(issue.Number))
	githubactions.SetEnv("TFACTION_DRIFT_ISSUE_STATE", issue.State)

	return nil
}
