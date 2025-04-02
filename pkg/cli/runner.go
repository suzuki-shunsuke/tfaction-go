package cli

import (
	"context"
	"io"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"
)

type Runner struct {
	Stdin   io.Reader
	Stdout  io.Writer
	Stderr  io.Writer
	LDFlags *LDFlags
	LogE    *logrus.Entry
}

type LDFlags struct {
	Version string
	Commit  string
	Date    string
}

func (runner *Runner) Run(ctx context.Context, args ...string) error {
	app := cli.Command{
		Name:    "tfaction",
		Usage:   "GitHub Actions Workflow for Terraform. https://github/com/suzuki-shunsuke/tfaction-go",
		Version: runner.LDFlags.Version + " (" + runner.LDFlags.Commit + ")",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "log-level",
				Usage:   "log level",
				Sources: cli.EnvVars("TFACTION_LOG_LEVEL"),
			},
		},
		EnableShellCompletion: true,
		Commands: []*cli.Command{
			runner.newVersionCommand(),
			runner.newCreateDriftIssuesCommand(),
			runner.newPickOutDriftIssuesCommand(),
			runner.newGetOrCreateDriftIssueCommand(),
		},
	}

	return app.Run(ctx, args) //nolint:wrapcheck
}
