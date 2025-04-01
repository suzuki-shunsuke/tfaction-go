package cli

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/afero"
	issues "github.com/suzuki-shunsuke/tfaction-go/pkg/controller/create-drift-issues"
	"github.com/suzuki-shunsuke/tfaction-go/pkg/github"
	"github.com/suzuki-shunsuke/tfaction-go/pkg/log"
	"github.com/urfave/cli/v3"
)

func (runner *Runner) newCreateDriftIssuesCommand() *cli.Command {
	return &cli.Command{
		Name:   "create-drift-issues",
		Usage:  "Create GitHub Issues for Terraform drift detection",
		Action: runner.createDriftIssuesAction,
	}
}

func (runner *Runner) createDriftIssuesAction(c *cli.Context) error {
	gh, err := github.New(c.Context, &github.ParamNew{
		Token:              os.Getenv("GITHUB_TOKEN"),
		GHEBaseURL:         os.Getenv("GITHUB_API_URL"),
		GHEGraphQLEndpoint: os.Getenv("GITHUB_GRAPHQL_URL"),
	})
	if err != nil {
		return fmt.Errorf("set up a GitHub Client: %w", err)
	}
	fs := afero.NewOsFs()
	ctrl := issues.New(gh, fs)
	log.SetLevel(c.String("log-level"), runner.LogE)
	repo := os.Getenv("GITHUB_REPOSITORY")
	repoOwner, repoName, found := strings.Cut(repo, "/")
	if !found {
		return errors.New("GITHUB_REPOSITORY is invalid")
	}
	pwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get a current directory path: %w", err)
	}
	return ctrl.Run(c.Context, runner.LogE, &issues.Param{ //nolint:wrapcheck
		RepoOwner: repoOwner,
		RepoName:  repoName,
		PWD:       pwd,
	})
}
