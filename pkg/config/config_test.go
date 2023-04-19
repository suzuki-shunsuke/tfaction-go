package config_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/suzuki-shunsuke/tfaction-go/pkg/config"
	"github.com/suzuki-shunsuke/tfaction-go/pkg/testutil"
)

func TestRead(t *testing.T) { //nolint:funlen
	t.Parallel()
	data := []struct {
		name  string
		isErr bool
		exp   *config.Config
		files map[string]string
	}{
		{
			name:  "tfaction-root.yaml is not found",
			isErr: true,
		},
		{
			name:  "invalid yaml",
			isErr: true,
			files: map[string]string{
				"tfaction-root.yaml": "}",
			},
		},
		{
			name: "empty yaml",
			files: map[string]string{
				"tfaction-root.yaml": "{}",
			},
			exp: &config.Config{
				WorkingDirectoryFile: "tfaction.yaml",
			},
		},
		{
			name: "empty drift detection",
			files: map[string]string{
				"tfaction-root.yaml": `
base_working_directory: tfaction
working_directory_file: tf.yaml
drift_detection: {}
target_groups:
- working_directory: aws/
  target: aws/
`,
			},
			exp: &config.Config{
				BaseWorkingDirectory: "tfaction",
				WorkingDirectoryFile: "tf.yaml",
				DriftDetection: &config.DriftDetection{
					NumOfIssues: 1,
					Duration:    168,
				},
				TargetGroups: []*config.TargetGroup{
					{
						WorkingDirectory: "aws/",
						Target:           "aws/",
					},
				},
			},
		},
		{
			name: "overwrite default values",
			files: map[string]string{
				"tfaction-root.yaml": `
base_working_directory: tfaction
working_directory_file: tf.yaml
drift_detection:
  num_of_issues: 3
  duration: 5
  issue_repo_owner: foo
  issue_repo_name: bar
target_groups:
- working_directory: aws/
  target: aws/
`,
			},
			exp: &config.Config{
				BaseWorkingDirectory: "tfaction",
				WorkingDirectoryFile: "tf.yaml",
				DriftDetection: &config.DriftDetection{
					NumOfIssues:    3,
					Duration:       5,
					IssueRepoOwner: "foo",
					IssueRepoName:  "bar",
				},
				TargetGroups: []*config.TargetGroup{
					{
						WorkingDirectory: "aws/",
						Target:           "aws/",
					},
				},
			},
		},
	}
	for _, d := range data {
		d := d
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			fs, err := testutil.NewFs(d.files)
			if err != nil {
				t.Fatal(err)
			}
			cfg, err := config.Read(fs)
			if err != nil {
				if d.isErr {
					return
				}
				t.Fatal(err)
			}
			if d.isErr {
				t.Fatal("error must be returned")
			}
			if diff := cmp.Diff(d.exp, cfg); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
