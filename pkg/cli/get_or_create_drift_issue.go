package cli //nolint:dupl

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/afero"
	issue "github.com/suzuki-shunsuke/tfaction-go/pkg/controller/get-or-create-drift-issue"
	"github.com/suzuki-shunsuke/tfaction-go/pkg/github"
	"github.com/suzuki-shunsuke/tfaction-go/pkg/log"
	"github.com/urfave/cli/v2"
)

func (runner *Runner) newGetOrCreateDriftIssueCommand() *cli.Command {
	return &cli.Command{
		Name:   "get-or-create-drift-issue",
		Usage:  "Get or Create a GitHub Issue for Terraform drift detection",
		Action: runner.getOrCreateDriftIssueAction,
	}
}

func (runner *Runner) getOrCreateDriftIssueAction(c *cli.Context) error {
	gh, err := github.New(c.Context, &github.ParamNew{
		Token: os.Getenv("GITHUB_TOKEN"),
	})
	if err != nil {
		return fmt.Errorf("set up a GitHub Client: %w", err)
	}
	fs := afero.NewOsFs()
	ctrl := issue.New(gh, fs)
	log.SetLevel(c.String("log-level"), runner.LogE)
	repo := os.Getenv("GITHUB_REPOSITORY")
	repoOwner, repoName, found := strings.Cut(repo, "/")
	if !found {
		return errors.New("GITHUB_REPOSITORY is invalid")
	}
	return ctrl.Run(c.Context, runner.LogE, &issue.Param{ //nolint:wrapcheck
		RepoOwner: repoOwner,
		RepoName:  repoName,
	})
}
