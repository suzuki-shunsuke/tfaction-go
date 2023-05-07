package issues_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	issues "github.com/suzuki-shunsuke/tfaction-go/pkg/controller/create-drift-issues"
	"github.com/suzuki-shunsuke/tfaction-go/pkg/github"
	"github.com/suzuki-shunsuke/tfaction-go/pkg/util"
)

func TestController_Run(t *testing.T) { //nolint:funlen
	t.Parallel()
	data := []struct {
		name  string
		param *issues.Param
		files map[string]string
		dirs  []string
		isErr bool
	}{
		{
			name: "normal",
			param: &issues.Param{
				RepoOwner: "suzuki-shunsuke",
				RepoName:  "test-tfaction",
				PWD:       "/home/foo/workspace/test-tfaction",
			},
			files: map[string]string{
				"tfaction-root.yaml": `
drift_detection: {}
target_groups:
- working_directory: aws/
  target: aws/
`,
				"/home/foo/workspace/test-tfaction/aws/foo/development/tfaction.yaml": "{}",
				"/home/foo/workspace/test-tfaction/aws/foo/production/tfaction.yaml":  "{}",
				"/home/foo/workspace/test-tfaction/aws/bar/development/tfaction.yaml": "{}",
				"/home/foo/workspace/test-tfaction/aws/bar/production/tfaction.yaml":  "{}",
			},
		},
	}
	ctx := context.Background()
	logE := logrus.NewEntry(logrus.New())
	for _, d := range data {
		d := d
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			gh := github.NewMockClient(t)
			gh.EXPECT().ListIssues(ctx, "suzuki-shunsuke", "test-tfaction").Return([]*github.Issue{
				{
					Number: 1,
					Title:  "Terraform Drift (aws/foo/development)",
					Target: "aws/foo/development",
					State:  "open",
				},
				{
					Number: 2,
					Title:  "Terraform Drift (aws/zoo/development)",
					Target: "aws/zoo/development",
					State:  "open",
				},
			}, nil)

			gh.EXPECT().CreateIssue(ctx, "suzuki-shunsuke", "test-tfaction", &github.IssueRequest{
				Title: util.StrP("Terraform Drift (aws/foo/production)"),
				Body:  util.StrP(issues.IssueBodyTemplate),
			}).Return(&github.GitHubIssue{
				Number: util.IntP(3),
			}, nil)

			gh.EXPECT().CreateIssue(ctx, "suzuki-shunsuke", "test-tfaction", &github.IssueRequest{
				Title: util.StrP("Terraform Drift (aws/bar/development)"),
				Body:  util.StrP(issues.IssueBodyTemplate),
			}).Return(&github.GitHubIssue{
				Number: util.IntP(4),
			}, nil)

			gh.EXPECT().CreateIssue(ctx, "suzuki-shunsuke", "test-tfaction", &github.IssueRequest{
				Title: util.StrP("Terraform Drift (aws/bar/production)"),
				Body:  util.StrP(issues.IssueBodyTemplate),
			}).Return(&github.GitHubIssue{
				Number: util.IntP(5),
			}, nil)

			gh.EXPECT().CloseIssue(ctx, "suzuki-shunsuke", "test-tfaction", 3).Return(nil, nil)
			gh.EXPECT().CloseIssue(ctx, "suzuki-shunsuke", "test-tfaction", 4).Return(nil, nil)
			gh.EXPECT().CloseIssue(ctx, "suzuki-shunsuke", "test-tfaction", 5).Return(nil, nil)

			gh.EXPECT().ArchiveIssue(ctx, "suzuki-shunsuke", "test-tfaction", 2, "Archived Terraform Drift (aws/zoo/development)").Return(nil, nil)

			fs := afero.NewMemMapFs()
			for _, k := range d.files {
				if err := util.MkdirAll(fs, k); err != nil {
					t.Fatal(err)
				}
			}
			for k, v := range d.files {
				if err := util.MkdirAll(fs, filepath.Dir(k)); err != nil {
					t.Fatal(err)
				}
				if err := util.WriteFile(fs, k, []byte(v)); err != nil {
					t.Fatal(err)
				}
			}

			ctrl := issues.New(gh, fs)
			if err := ctrl.Run(ctx, logE, d.param); err != nil {
				if d.isErr {
					return
				}
				t.Fatal(err)
			}
			if d.isErr {
				t.Fatal("error must be returned")
			}
		})
	}
}
