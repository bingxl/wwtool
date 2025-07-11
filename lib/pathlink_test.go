package lib

import (
	"os"
	"path/filepath"
	"testing"
)

var (
	dir   = filepath.Join(os.TempDir(), "pathlink_test")
	src   = filepath.Join(dir, "src")
	path2 = filepath.Join(dir, "path2")
)

func TestCreateSymlink(t *testing.T) {

}
