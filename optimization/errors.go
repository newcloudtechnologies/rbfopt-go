/*
 * Copyright (c) New Cloud Technologies, Ltd. 2013-2022.
 * Author: Vitaly Isaev <vitaly.isaev@myoffice.team>
 * License: https://github.com/newcloudtechnologies/rbfopt-go/blob/master/LICENSE
 */

package optimization

import (
	"github.com/pkg/errors"
)

var (
	// ErrInvalidParameterCombination notifies optimizer about invalid combination of parameters.
	// CostFunction must return MaxCost and ErrInvalidParameterCombination if it happened.
	ErrInvalidParameterCombination = errors.New("invalid parameter combination")

	// ErrTooHighInvalidParameterCombinationCost is returned when one sets too high
	// InvalidParameterCombinationCost setting
	ErrTooHighInvalidParameterCombinationCost = errors.New(
		"too high value of InvalidParameterCombinationCost: " +
			"visit https://github.com/coin-or/rbfopt/issues/28#issuecomment-629720480 to pick a good one")

	// ErrTooLowInvalidParameterCombinationCost is returned when one sets too low
	// InvalidParameterCombinationCost setting
	ErrTooLowInvalidParameterCombinationCost = errors.New(
		"observed cost value is higher than invalid parameter combination cost",
	)
)
