package cli

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/sethvargo/go-githubactions"
	"github.com/spf13/afero"
	issue "github.com/suzuki-shunsuke/tfaction-go/pkg/controller/get-or-create-drift-issue"
	"github.com/suzuki-shunsuke/tfaction-go/pkg/github"
	"github.com/suzuki-shunsuke/tfaction-go/pkg/log"
	"github.com/urfave/cli/v3"
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
		Token:              os.Getenv("GITHUB_TOKEN"),
		GHEBaseURL:         os.Getenv("GITHUB_API_URL"),
		GHEGraphQLEndpoint: os.Getenv("GITHUB_GRAPHQL_URL"),
	})
	if err != nil {
		return fmt.Errorf("set up a GitHub Client: %w", err)
	}
	fs := afero.NewOsFs()
	ctrl := issue.New(gh, fs, githubactions.New())
	log.SetLevel(c.String("log-level"), runner.LogE)
	repo := os.Getenv("GITHUB_REPOSITORY")
	repoOwner, repoName, found := strings.Cut(repo, "/")
	if !found {
		return errors.New("GITHUB_REPOSITORY is invalid")
	}
	target := os.Getenv("TFACTION_TARGET")
	if target == "" {
		return errors.New("TFACTION_TARGET is not set")
	}
	return ctrl.Run(c.Context, runner.LogE, &issue.Param{ //nolint:wrapcheck
		RepoOwner:       repoOwner,
		RepoName:        repoName,
		Target:          target,
		GitHubServerURL: os.Getenv("GITHUB_SERVER_URL"),
	})
}
