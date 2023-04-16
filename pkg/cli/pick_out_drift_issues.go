package cli

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/afero"
	issues "github.com/suzuki-shunsuke/tfaction-go/pkg/controller/pick-out-drift-issues"
	"github.com/suzuki-shunsuke/tfaction-go/pkg/github"
	"github.com/suzuki-shunsuke/tfaction-go/pkg/log"
	"github.com/urfave/cli/v2"
)

func (runner *Runner) newPickOutDriftIssuesCommand() *cli.Command {
	return &cli.Command{
		Name:   "pick-out-drift-issues",
		Usage:  "Pick out GitHub Issues for Terraform drift detection",
		Action: runner.pickOutDriftIssuesAction,
	}
}

func (runner *Runner) pickOutDriftIssuesAction(c *cli.Context) error {
	gh, err := github.New(c.Context, &github.ParamNew{
		Token: os.Getenv("GITHUB_TOKEN"),
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
	return ctrl.Run(c.Context, runner.LogE, &issues.Param{ //nolint:wrapcheck
		RepoOwner: repoOwner,
		RepoName:  repoName,
	})
}
