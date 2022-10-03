/*
 * Copyright (c) New Cloud Technologies, Ltd. 2013-2022.
 * Author: Vitaly Isaev <vitaly.isaev@myoffice.team>
 * License: https://github.com/newcloudtechnologies/rbfopt-go/blob/master/LICENSE
 */

package optimization

import (
	"context"
	"fmt"
	"math"

	"github.com/pkg/errors"
)

// Cost represents the value returned by cost function.
type Cost = float64

// CostFunction (or objective function) is implemented by clients.
// Optimizer will try to find the best possible combination of your parameters on the basis of this function.
// CostFunction call is expected to be expensive, so client should check context expiration.
type CostFunction func(ctx context.Context) (Cost, error)

// InitStrategy determines the way RBFOpt selects the initial sample points:
// See:
// 1. https://github.com/coin-or/rbfopt/blob/4.2.3/src/rbfopt/rbfopt_settings.py#L119
// 2. Https://usermanual.wiki/Document/manual.890598645/view (part 2.1)
type InitStrategy int8

const (
	// LHDMaximin - the default init strategy
	LHDMaximin InitStrategy = iota
	// LHDCorr - one of possible init strategies
	LHDCorr
	// AllCorners - one of possible init strategies
	AllCorners
	// LowerCorners - one of possible init strategies
	LowerCorners
	// RandCorners - one of possible init strategies
	RandCorners
)

// ErrUnknownInitStrategy is returned when user set unknown InitStrategy
var ErrUnknownInitStrategy = errors.New("unknown strategy")

// MarshalJSON renders InitStrategy to JSON.
func (s InitStrategy) MarshalJSON() ([]byte, error) {
	switch s {
	case LHDMaximin:
		return []byte("\"lhd_maximin\""), nil
	case LHDCorr:
		return []byte("\"lhd_corr\""), nil
	case AllCorners:
		return []byte("\"all_corners\""), nil
	case LowerCorners:
		return []byte("\"lower_corners\""), nil
	case RandCorners:
		return []byte("\"rand_corners\""), nil
	default:
		return nil, errors.Wrapf(ErrUnknownInitStrategy, "%v", s)
	}
}

// RBFOptConfig - RBFOpt configuration
type RBFOptConfig struct {
	CostFunction   CostFunction            `json:"-"`               // CostFunction itself
	Parameters     []*ParameterDescription `json:"parameters"`      // Arguments of a CostFunctions
	MaxEvaluations uint                    `json:"max_evaluations"` // Evaluations limit
	MaxIterations  uint                    `json:"max_iterations"`  // Iterations limit
	InitStrategy   InitStrategy            `json:"init_strategy"`   // Strategy to select initial points
	// RBFOpt: reason: https://github.com/coin-or/rbfopt/issues/28
	InvalidParameterCombinationCost Cost `json:"invalid_parameter_combination_cost"`
}

//nolint:revive // too simple function to split
func (c *RBFOptConfig) validate() error {
	if c == nil {
		return errors.New("empty")
	}

	if len(c.Parameters) == 0 {
		return errors.New("field Parameters is empty")
	}

	if c.CostFunction == nil {
		return errors.New("field CostFunction is empty")
	}

	for _, param := range c.Parameters {
		if err := param.validate(); err != nil {
			return errors.Wrapf(err, "validate parameter '%s'", param.Name)
		}
	}

	if c.MaxEvaluations == 0 {
		return errors.New("field MaxEvaluations is empty")
	}

	if c.MaxIterations == 0 {
		return errors.New("field MaxIterations is empty")
	}

	if c.InvalidParameterCombinationCost == math.MaxFloat64 {
		return ErrTooHighInvalidParameterCombinationCost
	}

	return nil
}

func (c *RBFOptConfig) getParameterByName(name string) (*ParameterDescription, error) {
	for _, param := range c.Parameters {
		if param.Name == name {
			return param, nil
		}
	}

	return nil, errors.Errorf("param '%s' does not exist", name)
}

// InvalidParameterCombinationRenderPolicy defines how to handle points corresponding to the
// ErrInvalidParameterCombination on plots
//go:generate stringer -type=InvalidParameterCombinationRenderPolicy
type InvalidParameterCombinationRenderPolicy int8

const (
	// Omit - points corresponding to the ErrInvalidParameterCombination won't be rendered at all
	Omit InvalidParameterCombinationRenderPolicy = iota + 1
	// AssignClosestValidValue - points corresponding to the ErrInvalidParameterCombination
	// will be rendered as if they have max observed valid value
	AssignClosestValidValue
)

// ErrUnknownInvalidParameterCombinationRenderPolicy notifies about wrong InvalidParameterCombinationRenderPolicy value
var ErrUnknownInvalidParameterCombinationRenderPolicy = errors.New("unknown InvalidParameterCombinationRenderPolicy")

// MarshalJSON renders InvalidParameterCombinationRenderPolicy to JSON
func (p InvalidParameterCombinationRenderPolicy) MarshalJSON() ([]byte, error) {
	switch p {
	case Omit:
		return []byte("\"omit\""), nil
	case AssignClosestValidValue:
		return []byte("\"assign_closest_valid_value\""), nil
	}

	return nil, errors.Wrapf(ErrUnknownInvalidParameterCombinationRenderPolicy, "%v", p)
}

// PlotConfig - plot renderer configuration
type PlotConfig struct {
	ScatterPlotPolicy   InvalidParameterCombinationRenderPolicy `json:"scatter_plot_policy"`
	HeatmapRenderPolicy InvalidParameterCombinationRenderPolicy `json:"heatmap_render_policy"`
}

func (c *PlotConfig) String() string {
	return fmt.Sprintf("scatter_%v_heatmap_%v", c.ScatterPlotPolicy, c.ScatterPlotPolicy)
}

func (c *PlotConfig) validate() error {
	if c == nil {
		return errors.New("empty plot config")
	}

	if c.ScatterPlotPolicy == 0 {
		return errors.New("field ScatterPlotIPCR is empty")
	}

	if c.HeatmapRenderPolicy == 0 {
		return errors.New("field HeatmapRenderErrIPCPolicy is empty")
	}

	return nil
}

// Config is top level configuration
//nolint:govet // have no time to find better configuration
type Config struct {
	// RBFOpt - config of rbfopt library itself
	RBFOpt *RBFOptConfig `json:"rbfopt"` //nolint:tagliatelle // RBFOpt is a well-known name, no underscore is needed
	// PlotConfig - config of plots made by wrapper
	Plot *PlotConfig `json:"plot"`
	// RootDir - place to store reports and other things
	// (optimizer will create it if it doesn't exist).
	RootDir string `json:"root_dir"`
	// Endpoint for the server that will work as a middleware
	Endpoint string `json:"endpoint"`
}

func (c *Config) validate() error {
	if c == nil {
		return errors.New("empty config")
	}

	if c.RootDir == "" {
		return errors.New("field RootDir is empty")
	}

	if len(c.Endpoint) == 0 {
		c.Endpoint = "0.0.0.0:8080"
	}

	if err := c.RBFOpt.validate(); err != nil {
		return errors.Wrap(err, "validate RBFOpt")
	}

	if err := c.Plot.validate(); err != nil {
		return errors.Wrap(err, "field Plot")
	}

	return nil
}
