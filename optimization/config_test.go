/*
 * Copyright (c) New Cloud Technologies, Ltd. 2013-2022.
 * Author: Vitaly Isaev <vitaly.isaev@myoffice.team>
 * License: https://github.com/newcloudtechnologies/rbfopt-go/blob/master/LICENSE
 */

package optimization

import (
	"context"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRBFOptConfig(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var c *RBFOptConfig
		require.Error(t, c.validate())
	})

	t.Run("no parameters", func(t *testing.T) {
		c := &RBFOptConfig{}
		require.Error(t, c.validate())
	})

	t.Run("no cost function", func(t *testing.T) {
		c := &RBFOptConfig{
			Parameters: []*ParameterDescription{
				{
					Bound: &Bound{
						Left:  1,
						Right: 2,
					},
					Name: "crab",
				},
			},
		}

		require.Error(t, c.validate())
	})

	t.Run("invalid parameters", func(t *testing.T) {
		c := &RBFOptConfig{
			Parameters: []*ParameterDescription{
				{
					Bound: &Bound{
						Left:  1,
						Right: 2,
					},
					Name: "crab",
				},
			},
			CostFunction: func(ctx context.Context) (Cost, error) {
				return 0, nil
			},
		}

		require.Error(t, c.validate())
	})

	t.Run("invalid max evaluations", func(t *testing.T) {
		c := &RBFOptConfig{
			Parameters: []*ParameterDescription{
				{
					Bound: &Bound{
						Left:  1,
						Right: 2,
					},
					ConfigModifier: func(i int) {},
					Name:           "crab",
				},
			},
			CostFunction: func(ctx context.Context) (Cost, error) {
				return 0, nil
			},
		}

		require.Error(t, c.validate())
	})

	t.Run("invalid max iterations", func(t *testing.T) {
		c := &RBFOptConfig{
			Parameters: []*ParameterDescription{
				{
					Bound: &Bound{
						Left:  1,
						Right: 2,
					},
					ConfigModifier: func(i int) {},
					Name:           "crab",
				},
			},
			CostFunction: func(ctx context.Context) (Cost, error) {
				return 0, nil
			},
			MaxEvaluations: 100,
		}

		require.Error(t, c.validate())
	})

	t.Run("invalid parameter combination cost", func(t *testing.T) {
		c := &RBFOptConfig{
			Parameters: []*ParameterDescription{
				{
					Bound: &Bound{
						Left:  1,
						Right: 2,
					},
					ConfigModifier: func(i int) {},
					Name:           "crab",
				},
			},
			CostFunction: func(ctx context.Context) (Cost, error) {
				return 0, nil
			},
			MaxEvaluations:                  100,
			MaxIterations:                   100,
			InvalidParameterCombinationCost: math.MaxFloat64,
		}

		require.Error(t, c.validate())
	})
}

func TestPlotConfig(t *testing.T) {
	t.Run("invalid scatter plot policy", func(t *testing.T) {
		c := &PlotConfig{
			ScatterPlotPolicy: 0,
		}
		require.Error(t, c.validate())
	})

	t.Run("invalid heatmap render policy", func(t *testing.T) {
		c := &PlotConfig{
			ScatterPlotPolicy: Omit,
		}
		require.Error(t, c.validate())
	})
}

func TestConfig(t *testing.T) {
	t.Run("no root dir", func(t *testing.T) {
		c := &Config{}
		require.Error(t, c.validate())
	})

	t.Run("no endpoint", func(t *testing.T) {
		c := &Config{
			RootDir: "/tmp/crab",
		}
		require.Error(t, c.validate())
	})

	t.Run("invalid RBFOpt", func(t *testing.T) {
		c := &Config{
			RootDir:  "/tmp/crab",
			Endpoint: "localhost:8080",
		}
		require.Error(t, c.validate())
	})

	t.Run("invalid Plot", func(t *testing.T) {
		c := &Config{
			RootDir:  "/tmp/crab",
			Endpoint: "localhost:8080",
			RBFOpt: &RBFOptConfig{
				Parameters: []*ParameterDescription{
					{
						Bound: &Bound{
							Left:  1,
							Right: 2,
						},
						ConfigModifier: func(i int) {},
						Name:           "crab",
					},
				},
				CostFunction: func(ctx context.Context) (Cost, error) {
					return 0, nil
				},
				MaxEvaluations:                  100,
				MaxIterations:                   100,
				InvalidParameterCombinationCost: 100,
			},
		}
		require.Error(t, c.validate())
	})

	t.Run("valid", func(t *testing.T) {
		c := &Config{
			RootDir:  "/tmp/crab",
			Endpoint: "localhost:8080",
			RBFOpt: &RBFOptConfig{
				Parameters: []*ParameterDescription{
					{
						Bound: &Bound{
							Left:  1,
							Right: 2,
						},
						ConfigModifier: func(i int) {},
						Name:           "crab",
					},
				},
				CostFunction: func(ctx context.Context) (Cost, error) {
					return 0, nil
				},
				MaxEvaluations:                  100,
				MaxIterations:                   100,
				InvalidParameterCombinationCost: 100,
			},
			Plot: &PlotConfig{
				ScatterPlotPolicy:   Omit,
				HeatmapRenderPolicy: AssignClosestValidValue,
			},
		}
		require.NoError(t, c.validate())
	})
}
