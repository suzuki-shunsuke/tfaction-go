package cli

import (
	"context"
	"io"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/urfave-cli-v3-util/helpall"
	"github.com/suzuki-shunsuke/urfave-cli-v3-util/vcmd"
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

func (r *Runner) Run(ctx context.Context, args ...string) error {
	return helpall.With(&cli.Command{ //nolint:wrapcheck
		Name:    "tfaction",
		Usage:   "GitHub Actions Workflow for Terraform. https://github/com/suzuki-shunsuke/tfaction-go",
		Version: r.LDFlags.Version + " (" + r.LDFlags.Commit + ")",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "log-level",
				Usage:   "log level",
				Sources: cli.EnvVars("TFACTION_LOG_LEVEL"),
			},
		},
		EnableShellCompletion: true,
		Commands: []*cli.Command{
			r.newCreateDriftIssuesCommand(),
			r.newPickOutDriftIssuesCommand(),
			r.newGetOrCreateDriftIssueCommand(),
			vcmd.New(&vcmd.Command{
				Name:    "tfaction",
				Version: r.LDFlags.Version,
				SHA:     r.LDFlags.Commit,
			}),
		},
	}, nil).Run(ctx, args)
}
