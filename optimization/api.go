/*
 * Copyright (c) New Cloud Technologies, Ltd. 2013-2022.
 * Author: Vitaly Isaev <vitaly.isaev@myoffice.team>
 * License: https://github.com/newcloudtechnologies/rbfopt-go/blob/master/LICENSE
 */

package optimization

type estimateCostRequest struct {
	ParameterValues []*ParameterValue `json:"parameter_values"`
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
