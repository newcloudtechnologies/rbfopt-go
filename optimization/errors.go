/*
 * Copyright (c) New Cloud Technologies, Ltd. 2013-2022.
 * Author: Vitaly Isaev <vitaly.isaev@myoffice.team>
 * License: https://github.com/newcloudtechnologies/rbfopt-go/blob/master/LICENSE
 */

package optimization

import (
	"github.com/pkg/errors"
)

// ErrInvalidParameterCombination notifies optimizer about invalid combination of parameters.
// CostFunction must return MaxCost and ErrInvalidParameterCombination if it happened.
var ErrInvalidParameterCombination = errors.New("invalid parameter combination")

var ErrTooHighInvalidParameterCombinationCost = errors.New(
	"too high value of InvalidParameterCombinationCost: " +
		"visit https://github.com/coin-or/rbfopt/issues/28#issuecomment-629720480 to pick a good one")

var ErrTooHighObservedCost = errors.New(
	"observed cost value is higher than invalid parameter combination cost",
)
