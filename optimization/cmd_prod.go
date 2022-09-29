//go:build !testing

package optimization

import (
	"os/exec"
)

const rbfOptGoExecutable = "rbfopt-go-wrapper"

func makeCmd(rootDir string) *exec.Cmd {
	cmd := exec.Command(rbfOptGoExecutable, rootDir)

	return cmd
}