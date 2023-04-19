package github_test

import (
	"context"
	"testing"

	"github.com/suzuki-shunsuke/tfaction-go/pkg/github"
)

func TestNew(t *testing.T) {
	t.Parallel()
	data := []struct {
		name  string
		param *github.ParamNew
		isErr bool
	}{
		{
			name: "normal",
			param: &github.ParamNew{
				Token: "dummy",
			},
		},
	}
	for _, d := range data {
		d := d
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			cl, err := github.New(ctx, d.param)
			if err != nil {
				if d.isErr {
					return
				}
				t.Fatal(err)
			}
			if cl == nil {
				t.Fatal("client is nil")
			}
		})
	}
}
