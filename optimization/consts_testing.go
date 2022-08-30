//go:build testing

package optimization

import (
	"go/build"
	"path/filepath"
)

var rbfOptGoExecutable = filepath.Join(build.Default.GOPATH, "src/github.com/newcloudtechnologies/rbfopt-go/wrapper/main.py")
