package issues

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
	"github.com/suzuki-shunsuke/tfaction-go/pkg/github"
	"github.com/suzuki-shunsuke/tfaction-go/pkg/util"
	"gopkg.in/yaml.v3"
)

type Controller struct {
	gh github.Client
}

func New(gh github.Client) *Controller {
	return &Controller{
		gh: gh,
	}
}

type Config struct {
	BaseWorkingDirectory string        `yaml:"base_working_directory"`
	WorkingDirectoryFile string        `yaml:"working_directory_file"`
	TargetGroups         []TargetGroup `yaml:"target_groups"`
}

type TargetGroup struct {
	WorkingDirectory string `yaml:"working_directory"`
	Target           string
}

func (cfg *Config) GetWorkingDirectoryFile() string {
	if cfg.WorkingDirectoryFile != "" {
		return cfg.WorkingDirectoryFile
	}
	return "tfaction.yaml"
}

type Param struct {
	RepoOwner string
	RepoName  string
	PWD       string
}

func (ctrl *Controller) Run(ctx context.Context, logE *logrus.Entry, param *Param) error { //nolint:cyclop,funlen
	f, err := os.Open("tfaction-root.yaml")
	if err != nil {
		return fmt.Errorf("open tfaction-root.yaml: %w", err)
	}
	defer f.Close()
	cfg := &Config{}
	if err := yaml.NewDecoder(f).Decode(cfg); err != nil {
		return fmt.Errorf("read tfaction-root.yaml: %w", err)
	}
	workingDirectoryFileName := cfg.GetWorkingDirectoryFile()
	workingDirectoryPaths := []string{}
	baseWorkingDirectory := filepath.Join(param.PWD, cfg.BaseWorkingDirectory)
	if err := filepath.WalkDir(baseWorkingDirectory, func(p string, dirEntry fs.DirEntry, e error) error {
		if dirEntry.Name() != workingDirectoryFileName {
			return nil
		}
		f, err := filepath.Rel(baseWorkingDirectory, filepath.Dir(p))
		if err != nil {
			return fmt.Errorf("get a relative path of a working directory: %w", err)
		}
		workingDirectoryPaths = append(workingDirectoryPaths, f)
		return nil
	}); err != nil {
		return fmt.Errorf("search working directories: %w", err)
	}
	logE.WithField("num_of_working_dirs", len(workingDirectoryPaths)).Debug("search working directories")
	// Convert working directory to target
	targets := make([]string, 0, len(workingDirectoryPaths))
	for _, workingDirectoryPath := range workingDirectoryPaths {
		for _, targetGroup := range cfg.TargetGroups {
			if strings.HasPrefix(workingDirectoryPath, targetGroup.WorkingDirectory) {
				targets = append(targets, strings.Replace(workingDirectoryPath, targetGroup.WorkingDirectory, targetGroup.Target, 1))
				break
			}
		}
	}
	logE.WithField("num_of_targets", len(targets)).Debug("convert working directories to targets")
	// Search GitHub Issues
	issues, err := ctrl.gh.ListIssues(ctx, param.RepoOwner, param.RepoName)
	if err != nil {
		return fmt.Errorf("list issues: %w", err)
	}
	logE.WithField("num_of_issues", len(issues)).Debug("search GiHub issues")
	issueTargets := make(map[string]struct{}, len(issues))
	for _, issue := range issues {
		issueTargets[issue.Target] = struct{}{}
	}
	for _, target := range targets {
		logE := logE.WithField("target", target)
		if _, ok := issueTargets[target]; ok {
			continue
		}
		if _, err := ctrl.gh.CreateIssue(ctx, param.RepoOwner, param.RepoName, &github.IssueRequest{
			Title: util.StrP(fmt.Sprintf(`Terraform Drift (%s)`, target)),
			Body: util.StrP(`
This issus was created by [tfaction](https://suzuki-shunsuke.github.io/tfaction/docs/).

## :warning: Don't change the issue title

tfaction searches Issues by Issue title. So please don't change the issue title.
`),
		}); err != nil {
			logerr.WithError(logE, err).Error("create an issue")
		}
		logE.Info("created an issue")
	}
	return nil
}
