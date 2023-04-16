package config

import (
	"fmt"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

/*
drift_detection:
  number_of_issues: 1
  duration: 168 # 24 * 7 = 7 days
*/

type configRaw struct {
	BaseWorkingDirectory string          `yaml:"base_working_directory"`
	WorkingDirectoryFile string          `yaml:"working_directory_file"`
	TargetGroups         []TargetGroup   `yaml:"target_groups"`
	DriftDetection       *DriftDetection `yaml:"drift_detection"`
}

func (cr *configRaw) Config() *Config {
	cfg := &Config{
		BaseWorkingDirectory: cr.BaseWorkingDirectory,
		WorkingDirectoryFile: cr.WorkingDirectoryFile,
		TargetGroups:         cr.TargetGroups,
		DriftDetection:       cr.DriftDetection,
	}
	if cfg.DriftDetection != nil {
		if cfg.DriftDetection.NumOfIssues == 0 {
			cfg.DriftDetection.NumOfIssues = 1
		}
		if cfg.DriftDetection.Duration == 0 {
			cfg.DriftDetection.Duration = 168
		}
	}
	return cfg
}

type Config struct {
	BaseWorkingDirectory string
	WorkingDirectoryFile string
	TargetGroups         []TargetGroup
	DriftDetection       *DriftDetection
}

type DriftDetection struct {
	NumOfIssues int
	Duration    int
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
