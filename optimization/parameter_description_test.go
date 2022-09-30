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

func TestParameterDescription(t *testing.T) {
	t.Run("invalid name pattern", func(t *testing.T) {
		pd := &ParameterDescription{
			Name: "🦃",
		}

		require.Error(t, pd.validate())
	})

	t.Run("empty config modifier", func(t *testing.T) {
		pd := &ParameterDescription{
			Name: "good_name",
		}

		require.Error(t, pd.validate())
	})
}
