//go:build testing

/*
 * Copyright (c) New Cloud Technologies, Ltd. 2013-2022.
 * Author: Vitaly Isaev <vitaly.isaev@myoffice.team>
 * License: https://github.com/newcloudtechnologies/rbfopt-go/blob/master/LICENSE
 */

package optimization

import (
	"fmt"
	"go/build"
	"os"
	"os/exec"
	"path/filepath"
)

func makeCmd(rootDir string) *exec.Cmd {
	projectPath := filepath.Join(
		build.Default.GOPATH,
		"src/github.com/newcloudtechnologies/rbfopt-go",
	)

	rbfOptGoExecutable := filepath.Join(projectPath, "rbfoptgo/main.py")

	cmd := exec.Command(rbfOptGoExecutable, rootDir)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("PYTHONPATH=%s", projectPath))

	return cmd
}
