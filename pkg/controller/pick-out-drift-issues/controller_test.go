package issues_test

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	issue "github.com/suzuki-shunsuke/tfaction-go/pkg/controller/get-or-create-drift-issue"
	issues "github.com/suzuki-shunsuke/tfaction-go/pkg/controller/pick-out-drift-issues"
	"github.com/suzuki-shunsuke/tfaction-go/pkg/github"
	"github.com/suzuki-shunsuke/tfaction-go/pkg/util"
)

func TestController_Run(t *testing.T) { //nolint:funlen
	t.Parallel()
	now := time.Date(2022, time.November, 10, 23, 0, 0, 0, time.UTC)
	deadline := "2022-11-03T23:00:00+00:00"
	data := []struct {
		name      string
		isErr     bool
		files     map[string]string
		param     *issues.Param
		setGH     func(ctx context.Context, gh *github.MockClient)
		setAction func(action *issue.MockAction)
	}{
		{
			name: "has issues",
			param: &issues.Param{
				RepoOwner: "suzuki-shunsuke",
				RepoName:  "test-tfaction",
				PWD:       "/home/foo/workspace/test-tfaction",
				Now:       now,
			},
			files: map[string]string{
				"tfaction-root.yaml": `
drift_detection:
  num_of_issues: 2
target_groups:
- working_directory: aws/
  target: aws/
`,
				"/home/foo/workspace/test-tfaction/aws/foo/development/tfaction.yaml": "{}",
				"/home/foo/workspace/test-tfaction/aws/foo/production/tfaction.yaml":  "{}",
				"/home/foo/workspace/test-tfaction/aws/bar/development/tfaction.yaml": "{}",
				"/home/foo/workspace/test-tfaction/aws/bar/production/tfaction.yaml":  "{}",
			},
			setGH: func(ctx context.Context, gh *github.MockClient) {
				gh.EXPECT().ListLeastRecentlyUpdatedIssues(ctx, "suzuki-shunsuke", "test-tfaction", 2, deadline).Return([]*github.Issue{
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

				gh.EXPECT().ArchiveIssue(ctx, "suzuki-shunsuke", "test-tfaction", 2, "Archived Terraform Drift (aws/zoo/development)").Return(nil, nil)
			},
			setAction: func(action *issue.MockAction) {
				action.EXPECT().SetOutput("has_issues", "true").Return()
				action.EXPECT().SetOutput("issues", `[{"number":1,"title":"Terraform Drift (aws/foo/development)","target":"aws/foo/development","state":"open","runs_on":"ubuntu-latest"}]`).Return()
			},
		},
	}
	logE := logrus.NewEntry(logrus.New())
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()

			gh := github.NewMockClient(t)
			d.setGH(t.Context(), gh)

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

			ctrl := issues.New(gh, fs, action)
			if err := ctrl.Run(t.Context(), logE, d.param); err != nil {
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
