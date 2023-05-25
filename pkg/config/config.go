package config

import (
	"fmt"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

/*
drift_detection:
  number_of_issues: 1
  minimum_detection_interval: 168 # 24 * 7 = 7 days
*/

type configRaw struct {
	BaseWorkingDirectory string          `yaml:"base_working_directory"`
	WorkingDirectoryFile string          `yaml:"working_directory_file"`
	TargetGroups         []*TargetGroup  `yaml:"target_groups"`
	DriftDetection       *DriftDetection `yaml:"drift_detection"`
	RunsOn               string          `yaml:"runs_on"`
}

func (cr *configRaw) Config() *Config {
	cfg := &Config{
		BaseWorkingDirectory: cr.BaseWorkingDirectory,
		WorkingDirectoryFile: cr.WorkingDirectoryFile,
		TargetGroups:         cr.TargetGroups,
		DriftDetection:       cr.DriftDetection,
		RunsOn:               cr.RunsOn,
	}
	if cfg.DriftDetection != nil {
		if cfg.DriftDetection.NumOfIssues == 0 {
			cfg.DriftDetection.NumOfIssues = 1
		}
		if cfg.DriftDetection.MinimumDetectionInterval == 0 {
			cfg.DriftDetection.MinimumDetectionInterval = 168
		}
	}
	if cfg.WorkingDirectoryFile == "" {
		cfg.WorkingDirectoryFile = "tfaction.yaml"
	}
	return cfg
}

type Config struct {
	BaseWorkingDirectory string
	WorkingDirectoryFile string
	TargetGroups         []*TargetGroup
	DriftDetection       *DriftDetection
	RunsOn               string
}

type WorkingDirectory struct {
	RunsOn              string `yaml:"runs_on"`
	TerraformPlanConfig *Job   `yaml:"terraform_plan_config"`
}

type DriftDetection struct {
	NumOfIssues              int    `yaml:"num_of_issues"`
	MinimumDetectionInterval int    `yaml:"minimum_detection_interval"`
	IssueRepoOwner           string `yaml:"issue_repo_owner"`
	IssueRepoName            string `yaml:"issue_repo_name"`
}

type TargetGroup struct {
	WorkingDirectory    string `yaml:"working_directory"`
	Target              string
	RunsOn              string `yaml:"runs_on"`
	TerraformPlanConfig *Job   `yaml:"terraform_plan_config"`
}

type Job struct {
	RunsOn string `yaml:"runs_on"`
}

func (job *Job) GetRunsOn() string {
	if job == nil {
		return ""
	}
	return job.RunsOn
}

func Read(fs afero.Fs) (*Config, error) {
	f, err := fs.Open("tfaction-root.yaml")
	if err != nil {
		return nil, fmt.Errorf("open tfaction-root.yaml: %w", err)
	}
	defer f.Close()
	cfg := &configRaw{}
	if err := yaml.NewDecoder(f).Decode(cfg); err != nil {
		return nil, fmt.Errorf("read tfaction-root.yaml: %w", err)
	}
	return cfg.Config(), nil
}

func ReadWorkingDirectory(fs afero.Fs, p string) (*WorkingDirectory, error) {
	f, err := fs.Open(p)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", p, err)
	}
	defer f.Close()
	cfg := &WorkingDirectory{}
	if err := yaml.NewDecoder(f).Decode(cfg); err != nil {
		return nil, fmt.Errorf("read %s: %w", p, err)
	}
	return cfg, nil
}
