/*
 * Copyright (c) New Cloud Technologies, Ltd. 2013-2022.
 * Author: Vitaly Isaev <vitaly.isaev@myoffice.team>
 * License: https://github.com/newcloudtechnologies/rbfopt-go/blob/master/LICENSE
 */

package optimization_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/newcloudtechnologies/rbfopt-go/optimization"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func newLogger() logr.Logger {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	zapLog, err := config.Build()
	if err != nil {
		panic(err)
	}

	return zapr.NewLogger(zapLog)
}

type serviceConfig struct {
	paramX int
	paramY int
	paramZ int
}

func (cfg *serviceConfig) setParamX(value int) { cfg.paramX = value }
func (cfg *serviceConfig) setParamY(value int) { cfg.paramY = value }
func (cfg *serviceConfig) setParamZ(value int) { cfg.paramZ = value }

func (cfg *serviceConfig) costFunction(_ context.Context) (optimization.Cost, error) {
	// using quite a simple polinomial function with minimum that can be easily discovered:
	// it corresponds to the upper bound of every variable
	x, y, z := cfg.paramX, cfg.paramY, cfg.paramZ

	return optimization.Cost(-1 * (x*y + z)), nil
}

func TestRbfoptGo(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		cfg := &serviceConfig{paramX: 0, paramY: 0, paramZ: 0}

		// set bounds to parameters
		config := &optimization.Config{
			RootDir: makeRootDirPath(),
			RBFOpt: &optimization.RBFOptConfig{
				CostFunction: cfg.costFunction,
				Parameters: []*optimization.ParameterDescription{
					{
						Name:           "x",
						Bound:          &optimization.Bound{Left: 0, Right: 10},
						ConfigModifier: cfg.setParamX,
					},
					{
						Name:           "y",
						Bound:          &optimization.Bound{Left: 0, Right: 10},
						ConfigModifier: cfg.setParamY,
					},
					{
						Name:           "z",
						Bound:          &optimization.Bound{Left: 0, Right: 10},
						ConfigModifier: cfg.setParamZ,
					},
				},
				MaxEvaluations:                  25,
				MaxIterations:                   25,
				InvalidParameterCombinationCost: 10,
			},
			Plot: &optimization.PlotConfig{
				ScatterPlotPolicy:   optimization.Omit,
				HeatmapRenderPolicy: optimization.Omit,
			},
		}

		logger := newLogger()
		ctx := logr.NewContext(context.Background(), logger)

		// perform optimization
		report, err := optimization.Optimize(ctx, config)
		require.NoError(t, err)
		require.NotNil(t, report)

		// validate that the optimum was reached
		expectedOptimumCfg := &serviceConfig{
			paramX: config.RBFOpt.Parameters[0].Bound.Right,
			paramY: config.RBFOpt.Parameters[1].Bound.Right,
			paramZ: config.RBFOpt.Parameters[2].Bound.Right,
		}

		expectedOptimumCost, err := expectedOptimumCfg.costFunction(context.Background())
		require.NoError(t, err)

		require.Equal(t, expectedOptimumCost, report.Cost)
		require.Len(t, report.Optimum, len(config.RBFOpt.Parameters))

		for i := 0; i < len(config.RBFOpt.Parameters); i++ {
			require.Equal(t, report.Optimum[i].Name, config.RBFOpt.Parameters[i].Name)
			require.Equal(t, report.Optimum[i].Value, config.RBFOpt.Parameters[i].Bound.Right)
		}
	})

	t.Run("invalid parameters combination", func(t *testing.T) {
		testCases := []*optimization.PlotConfig{
			{
				ScatterPlotPolicy:   optimization.Omit,
				HeatmapRenderPolicy: optimization.Omit,
			},
			{
				ScatterPlotPolicy:   optimization.Omit,
				HeatmapRenderPolicy: optimization.AssignClosestValidValue,
			},
		}

		for _, tc := range testCases {
			tc := tc
			t.Run(tc.String(), func(t *testing.T) {
				cfg := &serviceConfig{paramX: 0, paramY: 0, paramZ: 0}

				config := &optimization.Config{
					RootDir: makeRootDirPath(),
					RBFOpt: &optimization.RBFOptConfig{
						Parameters: []*optimization.ParameterDescription{
							{
								Name:           "x",
								Bound:          &optimization.Bound{Left: 0, Right: 10},
								ConfigModifier: cfg.setParamX,
							},
							{
								Name:           "y",
								Bound:          &optimization.Bound{Left: 0, Right: 10},
								ConfigModifier: cfg.setParamY,
							},
							{
								Name:           "z",
								Bound:          &optimization.Bound{Left: 0, Right: 10},
								ConfigModifier: cfg.setParamZ,
							},
						},
						// the cost function is almost the same as in previous test case, but now we impose some conditions on variables
						CostFunction: func(ctx context.Context) (optimization.Cost, error) {
							// explicitly validate parameters, throw errors if they're invalid
							if cfg.paramX < cfg.paramY {
								return 0, optimization.ErrInvalidParameterCombination
							}

							return cfg.costFunction(ctx)
						},
						MaxEvaluations:                  25,
						MaxIterations:                   25,
						InvalidParameterCombinationCost: 10,
					},
					Plot: tc,
				}

				logger := newLogger()
				ctx := logr.NewContext(context.Background(), logger)

				report, err := optimization.Optimize(ctx, config)
				require.NoError(t, err)
				require.NotNil(t, report)

				expectedOptimumCfg := &serviceConfig{
					paramX: config.RBFOpt.Parameters[0].Bound.Right,
					paramY: config.RBFOpt.Parameters[1].Bound.Right,
					paramZ: config.RBFOpt.Parameters[2].Bound.Right,
				}

				expectedOptimumCost, err := expectedOptimumCfg.costFunction(context.Background())
				require.NoError(t, err)

				require.Equal(t, expectedOptimumCost, report.Cost)
				require.Len(t, report.Optimum, len(config.RBFOpt.Parameters))

				for i := 0; i < len(config.RBFOpt.Parameters); i++ {
					require.Equal(t, report.Optimum[i].Name, config.RBFOpt.Parameters[i].Name)
					require.Equal(t, report.Optimum[i].Value, config.RBFOpt.Parameters[i].Bound.Right)
				}
			})
		}
	})
}

func makeRootDirPath() string {
	return fmt.Sprintf("/tmp/rbfopt_%s", time.Now().Format("20060102_150405"))
}
