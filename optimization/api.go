/*
 * Copyright (c) New Cloud Technologies, Ltd. 2013-2022.
 * Author: Vitaly Isaev <vitaly.isaev@myoffice.team>
 * License: https://github.com/newcloudtechnologies/rbfopt-go/blob/master/LICENSE
 */

package optimization

import (
	"github.com/pkg/errors"
)

type estimateCostRequest struct {
	ParameterValues []*ParameterValue `json:"parameter_values"`
}

func (r *estimateCostRequest) applyValues(config *RBFOptConfig) error {
	// apply all values to config first
	for _, pv := range r.ParameterValues {
		parameterDesc, err := config.getParameterByName(pv.Name)
		if err != nil {
			return errors.Wrapf(err, "get parameter by name: %s", pv.Name)
		}

		parameterDesc.ConfigModifier(pv.Value)
	}

	return nil
}

type estimateCostResponse struct {
	Cost                        float64 `json:"cost"`
	InvalidParameterCombination bool    `json:"invalid_parameter_combination"`
}

type registerReportRequest struct {
	Report *Report `json:"report"`
}

type registerReportResponse struct {
}
