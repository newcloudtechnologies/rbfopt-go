#!/usr/bin/env python

#  Copyright (c) New Cloud Technologies, Ltd. 2013-2022.
#  Author: Vitaly Isaev <vitaly.isaev@myoffice.team>
#  License: https://github.com/newcloudtechnologies/rbfopt-go/blob/master/LICENSE

import pathlib

import pandas as pd

from config import Config
from plot import Renderer
from report import Report

debug_dir = pathlib.Path('/tmp/rbfopt_20220830_105502')


def main():
    config = Config.from_file(debug_dir.joinpath("config.json"))
    evaluations = pd.read_csv(debug_dir.joinpath("evaluations.csv"))
    report = Report.load_from_file(debug_dir.joinpath("report.json"))
    renderer = Renderer(config=config, df=evaluations, report=report)
    renderer.radar()


if __name__ == '__main__':
    main()
