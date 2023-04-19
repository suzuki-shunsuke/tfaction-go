package testutil

import (
	"os"

	"github.com/spf13/afero"
)

const FilePermission os.FileMode = 0o644

func NewFs(files map[string]string) (afero.Fs, error) {
	fs := afero.NewMemMapFs()
	for name, body := range files {
		if err := afero.WriteFile(fs, name, []byte(body), FilePermission); err != nil {
			return nil, err //nolint:wrapcheck
		}
	}
	return fs, nil
}
