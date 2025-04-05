package issue_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	issues "github.com/suzuki-shunsuke/tfaction-go/pkg/controller/create-drift-issues"
	issue "github.com/suzuki-shunsuke/tfaction-go/pkg/controller/get-or-create-drift-issue"
	"github.com/suzuki-shunsuke/tfaction-go/pkg/github"
	"github.com/suzuki-shunsuke/tfaction-go/pkg/util"
)

func TestController_Run(t *testing.T) { //nolint:funlen
	t.Parallel()
	data := []struct {
		name      string
		isErr     bool
		files     map[string]string
		param     *issue.Param
		setGH     func(ctx context.Context, gh *github.MockClient)
		setAction func(action *issue.MockAction)
	}{
		{
			name: "create an issue",
			param: &issue.Param{
				RepoOwner:       "suzuki-shunsuke",
				RepoName:        "test-tfaction",
				Target:          "aws/bar/production",
				GitHubServerURL: "https://github.com",
			},
			files: map[string]string{
				"tfaction-root.yaml": `
drift_detection: {}
target_groups:
- working_directory: aws/
  target: aws/
`,
				"aws/bar/production/tfaction.yaml": `{}`,
			},
			setGH: func(ctx context.Context, gh *github.MockClient) {
				gh.EXPECT().CreateIssue(ctx, "suzuki-shunsuke", "test-tfaction", &github.IssueRequest{
					Title: util.StrP("Terraform Drift (aws/bar/production)"),
					Body:  util.StrP(issues.IssueBodyTemplate),
				}).Return(&github.GitHubIssue{
					Number: util.IntP(5),
				}, nil)

				gh.EXPECT().GetIssue(ctx, "suzuki-shunsuke", "test-tfaction", "Terraform Drift (aws/bar/production)").Return(nil, nil)
			},
			setAction: func(action *issue.MockAction) {
				action.EXPECT().SetEnv("TFACTION_DRIFT_ISSUE_NUMBER", "5").Return()
				action.EXPECT().SetEnv("TFACTION_DRIFT_ISSUE_STATE", "open").Return()
				action.EXPECT().Infof("https://github.com/suzuki-shunsuke/test-tfaction/pull/5").Return()
				action.EXPECT().AddStepSummary("Drift Issue: https://github.com/suzuki-shunsuke/test-tfaction/pull/5").Return()
			},
		},
		{
			name: "get an issue",
			param: &issue.Param{
				RepoOwner:       "suzuki-shunsuke",
				RepoName:        "test-tfaction",
				Target:          "aws/bar/production",
				GitHubServerURL: "https://github.com",
			},
			files: map[string]string{
				"tfaction-root.yaml": `
drift_detection: {}
target_groups:
- working_directory: aws/
  target: aws/
`,
				"aws/bar/production/tfaction.yaml": `{}`,
			},
			setGH: func(ctx context.Context, gh *github.MockClient) {
				gh.EXPECT().GetIssue(ctx, "suzuki-shunsuke", "test-tfaction", "Terraform Drift (aws/bar/production)").Return(&github.Issue{
					Number: 10,
					State:  "closed",
				}, nil)
			},
			setAction: func(action *issue.MockAction) {
				action.EXPECT().SetEnv("TFACTION_DRIFT_ISSUE_NUMBER", "10").Return()
				action.EXPECT().SetEnv("TFACTION_DRIFT_ISSUE_STATE", "closed").Return()
				action.EXPECT().Infof("https://github.com/suzuki-shunsuke/test-tfaction/pull/10").Return()
				action.EXPECT().AddStepSummary("Drift Issue: https://github.com/suzuki-shunsuke/test-tfaction/pull/10").Return()
			},
		},
	}
	logE := logrus.NewEntry(logrus.New())
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			gh := github.NewMockClient(t)
			d.setGH(ctx, gh)

			fs := afero.NewMemMapFs()
			for k, v := range d.files {
				if err := util.MkdirAll(fs, filepath.Dir(k)); err != nil {
					t.Fatal(err)
				}
				if err := util.WriteFile(fs, k, []byte(v)); err != nil {
					t.Fatal(err)
				}
			}

			action := issue.NewMockAction(t)
			d.setAction(action)

			ctrl := issue.New(gh, fs, action)
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
