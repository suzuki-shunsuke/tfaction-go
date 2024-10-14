package issue

import (
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/tfaction-go/pkg/github"
)

type Controller struct {
	gh     github.Client
	fs     afero.Fs
	action Action
}

func New(gh github.Client, fs afero.Fs, action Action) *Controller {
	return &Controller{
		gh:     gh,
		fs:     fs,
		action: action,
	}
}
