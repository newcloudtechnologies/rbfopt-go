#!/usr/bin/env python
#  Copyright (c) New Cloud Technologies, Ltd. 2013-2022.
#  Author: Vitaly Isaev <vitaly.isaev@myoffice.team>
#  License: https://github.com/newcloudtechnologies/rbfopt-go/blob/master/LICENSE

import pathlib

import pandas as pd
from matplotlib import rc

from rbfoptgo.config import Config
from rbfoptgo.plot import Renderer
from rbfoptgo.report import Report

rc('font', **{'family': 'serif', 'serif': ['Roboto']})


def run(debug_dir: pathlib.Path):
    """
    An entrypoint to plotting debug script. Use only for testing purposes.
    :return:
    """
    config = Config.from_file(debug_dir.joinpath("config.json"))
    evaluations = pd.read_csv(debug_dir.joinpath("evaluations.csv"))
    report = Report.load_from_file(debug_dir.joinpath("report.json"))
    renderer = Renderer(config=config, df=evaluations, report=report, transparent=False)
    renderer.heatmaps()


def main():
    """
    Change debug_dir if you want to re-render plots using the data left from
    the one of the previous optimization sessions.
    :return:
    """
    debug_dir = pathlib.Path('/home/isaev/trouble/optimization/20220908_2330_HDD_SSD_Volume_IO_CPU_262144')
    run(debug_dir)


if __name__ == '__main__':
    main()
