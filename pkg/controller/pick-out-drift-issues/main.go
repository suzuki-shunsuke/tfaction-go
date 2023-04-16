package issues

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sethvargo/go-githubactions"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/tfaction-go/pkg/config"
	"github.com/suzuki-shunsuke/tfaction-go/pkg/github"
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
}

func (ctrl *Controller) Run(ctx context.Context, logE *logrus.Entry, param *Param) error {
	cfg, err := config.Read(ctrl.fs)
	if err != nil {
		return fmt.Errorf("read tfaction-root.yaml: %w", err)
	}
	if cfg.DriftDetection == nil {
		return nil
	}

	deadline := time.Now().In(time.FixedZone("UTC", 0)).Add(-time.Duration(cfg.DriftDetection.Duration) * time.Hour).Format("2006-01-02T15:04:05+00:00")
	logE.WithFields(logrus.Fields{
		"duration": cfg.DriftDetection.Duration,
		"deadline": deadline,
	}).Info("check a deadline")

	// Search least recently updated issues
	issues, err := ctrl.gh.ListLeastRecentlyUpdatedIssues(ctx, param.RepoOwner, param.RepoName, cfg.DriftDetection.NumOfIssues, deadline)
	if err != nil {
		return fmt.Errorf("list drift issues: %w", err)
	}
	logE.WithField("num_of_issues", len(issues)).Info("list least recently updated issues")

	if len(issues) == 0 {
		githubactions.SetOutput("has_issues", "false")
		return nil
	}

	githubactions.SetOutput("has_issues", "true")
	b, err := json.Marshal(issues)
	if err != nil {
		return fmt.Errorf("marshal issues as JSON: %w", err)
	}
	githubactions.SetOutput("issues", string(b))
	return nil
}
