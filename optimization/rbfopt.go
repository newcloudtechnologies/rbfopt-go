/*
 * Copyright (c) New Cloud Technologies, Ltd. 2013-2022.
 * Author: Vitaly Isaev <vitaly.isaev@myoffice.team>
 * License: https://github.com/newcloudtechnologies/rbfopt-go/blob/master/LICENSE
 */

package optimization

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
)

type rbfOptWrapper struct {
	ctx    context.Context
	config *Config
}

const rbfOptGoExecutable = "/home/isaev/go/src/github.com/newcloudtechnologies/rbfopt-go/wrapper/main.py"

func (r *rbfOptWrapper) run() error {
	// render config to JSON because it will be used by Python part
	path := filepath.Join(r.config.RootDir, "config.json")

	data, err := json.Marshal(r.config)
	if err != nil {
		return errors.Wrap(err, "marshal json")
	}

	if err := ioutil.WriteFile(path, data, 0600); err != nil {
		return errors.Wrap(err, "write file")
	}

	//nolint:gosec
	cmd := exec.Command(rbfOptGoExecutable, r.config.RootDir)
	if err := r.executeCommand(r.ctx, cmd); err != nil {
		return errors.Wrap(err, "execute command")
	}

	return nil
}

func (r *rbfOptWrapper) executeCommand(ctx context.Context, cmd *exec.Cmd) error {
	logger, err := logr.FromContext(ctx)
	if err != nil {
		return errors.Wrap(err, "logr from context")
	}

	logger.Info("executing command", "cmd", cmd)

	var (
		stdoutBuf = &bytes.Buffer{}
		stderrBuf = &bytes.Buffer{}
	)

	cmd.Stdout = stdoutBuf
	cmd.Stderr = stderrBuf

	err = cmd.Run()

	// print

	if stdoutBuf.Len() > 0 {
		logger.V(0).Info("subprocess stdout")

		_, err := io.Copy(os.Stdout, stdoutBuf)
		if err != nil {
			return errors.Wrap(err, "copy")
		}
	}

	if stderrBuf.Len() > 0 {
		logger.Info("subprocess stderr")

		_, err := io.Copy(os.Stdout, stderrBuf)
		if err != nil {
			return errors.Wrap(err, "copy")
		}
	}

	if err != nil {
		return errors.Wrap(err, "cmd run")
	}

	return nil
}

func runRbfOpt(ctx context.Context, config *Config) error {
	wrapper := &rbfOptWrapper{
		ctx:    ctx,
		config: config,
	}

	if err := wrapper.run(); err != nil {
		return errors.Wrap(err, "run RbfOpt wrapper")
	}

	return nil
}
