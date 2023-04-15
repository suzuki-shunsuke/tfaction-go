package cli

import (
	"context"
	"io"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
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
	compiledDate, err := time.Parse(time.RFC3339, runner.LDFlags.Date)
	if err != nil {
		compiledDate = time.Now()
	}
	app := cli.App{
		Name:     "tfaction",
		Usage:    "GitHub Actions Workflow for Terraform. https://github/com/suzuki-shunsuke/tfaction-go",
		Version:  runner.LDFlags.Version + " (" + runner.LDFlags.Commit + ")",
		Compiled: compiledDate,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "log-level",
				Usage:   "log level",
				EnvVars: []string{"TFACTION_LOG_LEVEL"},
			},
		},
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			runner.newVersionCommand(),
			runner.newCreateDriftIssuesCommand(),
		},
	}

	return app.RunContext(ctx, args) //nolint:wrapcheck
}
