//go:build !testing

/*
 * Copyright (c) New Cloud Technologies, Ltd. 2013-2022.
 * Author: Vitaly Isaev <vitaly.isaev@myoffice.team>
 * License: https://github.com/newcloudtechnologies/rbfopt-go/blob/master/LICENSE
 */

package optimization

import (
	"os/exec"
)

const rbfOptGoExecutable = "rbfopt-go-wrapper"

func makeCmd(rootDir string) *exec.Cmd {
	cmd := exec.Command(rbfOptGoExecutable, rootDir)

	return cmd
}
