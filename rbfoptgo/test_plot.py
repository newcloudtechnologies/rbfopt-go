#  Copyright (c) New Cloud Technologies, Ltd. 2013-2022.
#  Author: Vitaly Isaev <vitaly.isaev@myoffice.team>
#  License: https://github.com/newcloudtechnologies/rbfopt-go/blob/master/LICENSE

import glob
import pathlib

from rbfoptgo.plot_debug import run


def test_plot():
    """
    Integration test to estimate approximate Python coverage
    """
    # Directory with evaluations left after Go tests
    debug_dir = glob.glob("/tmp/rbfopt*")[0]
    run(pathlib.Path(debug_dir))
