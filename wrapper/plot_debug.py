#  Copyright (c) New Cloud Technologies, Ltd. 2013-2022.
#  Author: Vitaly Isaev <vitaly.isaev@myoffice.team>
#  License: https://github.com/newcloudtechnologies/rbfopt-go/blob/master/LICENSE
import pathlib

import pandas as pd

from plot import Renderer
from report import Report
from settings import Settings

debug_dir = pathlib.Path('/home/isaev/trouble/optimization/crab_problem_1')


def main():
    ss = Settings.from_file(debug_dir.joinpath("settings.json"))
    evaluations = pd.read_csv(debug_dir.joinpath("evaluations.csv"))
    evaluations.drop("Unnamed: 0", axis=1, inplace=True)
    report = Report.load_from_file(debug_dir.joinpath("report.json"))
    renderer = Renderer(ss=ss, df=evaluations, report=report, root_dir=debug_dir)
    renderer.scatterplots()
    renderer.pairwise_heatmap_matrix()


if __name__ == '__main__':
    main()
