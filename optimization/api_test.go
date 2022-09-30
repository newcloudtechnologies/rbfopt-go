/*
 * Copyright (c) New Cloud Technologies, Ltd. 2013-2022.
 * Author: Vitaly Isaev <vitaly.isaev@myoffice.team>
 * License: https://github.com/newcloudtechnologies/rbfopt-go/blob/master/LICENSE
 */

package optimization

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEstimateCostRequest(t *testing.T) {
	cfg := &RBFOptConfig{
		Parameters: []*ParameterDescription{
			{
				Name: "Param1",
			},
		},
	}

	ecr := &estimateCostRequest{
		ParameterValues: []*ParameterValue{
			{
				Name:  "Param2",
				Value: 1,
			},
		},
	}

	require.Error(t, ecr.applyValues(cfg))
}
