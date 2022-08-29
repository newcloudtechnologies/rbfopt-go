/*
 * Copyright (c) New Cloud Technologies, Ltd. 2013-2022.
 * Author: Vitaly Isaev <vitaly.isaev@myoffice.team>
 * License: https://github.com/newcloudtechnologies/rbfopt-go/blob/master/LICENSE
 */

package optimization

import (
	"context"
	"os"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
)

// ParameterValue describes some value of a CostFunction argument.
type ParameterValue struct {
	Name  string
	Value int
}

// Report contains information about the finished optimization process
//nolint:govet
type Report struct {
	Cost            Cost              `json:"cost"`    // Discovered optimal value of a CostFunction
	Optimum         []*ParameterValue `json:"optimum"` // Parameter values matching the optimum point
	Iterations      int               `json:"iterations"`
	Evaluations     int               `json:"evaluations"`
	FastEvaluations int               `json:"fast_evaluations"`
}

// Optimize is an entry point for the optimization routines.
// One may want to pass logger within context to have detailed logs.
func Optimize(ctx context.Context, config *Config) (*Report, error) {
	logger := logr.FromContextOrDiscard(ctx)

	// check config
	if err := config.validate(); err != nil {
		return nil, errors.Wrap(err, "validate config")
	}

	// run HTTP server that will redirect requests from Python optimizer to your Go service
	estimator := newCostEstimator(config)

	srv := newServer(logger, config.Endpoint, estimator)
	defer srv.quit()

	// create root directory for configs and artifacts if necessary
	if _, err := os.Stat(config.RootDir); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(config.RootDir, 0755); err != nil {
				return nil, errors.Wrap(err, "make directory")
			}
		} else {
			return nil, errors.Wrap(err, "stat directory")
		}
	}

	// run Python optimizer
	ctx = logr.NewContext(ctx, logger)
	if err := runRbfOpt(ctx, config); err != nil {
		if srv.lastError != nil {
			return nil, errors.Wrap(srv.lastError, "run rbfopt")
		}

		return nil, errors.Wrap(err, "run rbfopt")
	}

	// obtain final report
	report := estimator.finalReport
	if report == nil {
		return nil, errors.New("protocol error: report is nil")
	}

	return report, nil
}
