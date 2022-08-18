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

type rbfOptSettings struct {
	Endpoint                               string                  `json:"endpoint"`
	Parameters                             []*ParameterDescription `json:"parameters"`
	MaxEvaluations                         uint                    `json:"max_evaluations"`
	MaxIterations                          uint                    `json:"max_iterations"`
	SkipInvalidParameterCombinationOnPlots bool                    `json:"skip_invalid_parameter_combination_on_plots"`
	InitStrategy                           InitStrategy            `json:"init_strategy"`
}

type rbfOptWrapper struct {
	ctx      context.Context
	settings *Settings
	endpoint string
}

const rbfOptGoExecutable = "rbfopt-go-wrapper"

func (r *rbfOptWrapper) run() error {
	path := filepath.Join(r.settings.RootDir, "settings.json")

	if err := r.dumpConfig(path); err != nil {
		return errors.Wrap(err, "dump config")
	}

	//nolint:gosec
	cmd := exec.Command(rbfOptGoExecutable, r.settings.RootDir)
	if err := r.executeCommand(r.ctx, cmd); err != nil {
		return errors.Wrap(err, "execute command")
	}

	return nil
}

func (r *rbfOptWrapper) dumpConfig(path string) error {
	cfg := &rbfOptSettings{
		Endpoint:                               r.endpoint,
		Parameters:                             r.settings.Parameters,
		MaxEvaluations:                         r.settings.MaxEvaluations,
		MaxIterations:                          r.settings.MaxIterations,
		SkipInvalidParameterCombinationOnPlots: r.settings.SkipInvalidParameterCombinationOnPlots,
	}

	data, err := json.Marshal(cfg)
	if err != nil {
		return errors.Wrap(err, "json marshal")
	}

	if err := ioutil.WriteFile(path, data, 0644); err != nil {
		return errors.Wrap(err, "write file")
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

func runRbfOpt(ctx context.Context, settings *Settings, endpoint string) error {
	wrapper := &rbfOptWrapper{
		ctx:      ctx,
		settings: settings,
		endpoint: endpoint,
	}

	if err := wrapper.run(); err != nil {
		return errors.Wrap(err, "run RbfOpt wrapper")
	}

	return nil
}
