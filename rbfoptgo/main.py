#!/usr/bin/env python

#  Copyright (c) New Cloud Technologies, Ltd. 2013-2022.
#  Author: Vitaly Isaev <vitaly.isaev@myoffice.team>
#  License: https://github.com/newcloudtechnologies/rbfopt-go/blob/master/LICENSE

import pathlib
import sys

import rbfopt

from rbfoptgo.client import Client
from rbfoptgo.config import Config
from rbfoptgo.evaluator import Evaluator
from rbfoptgo.plot import Renderer


def main():
    # prepare configuration
    root_dir = pathlib.Path(sys.argv[1])
    config_path = root_dir.joinpath("config.json")
    config = Config.from_file(config_path)
    print(f"config: {config}")

    client = Client(config.endpoint)
    evaluator = Evaluator(config=config, client=client, parameter_names=config.rbfopt.var_names, root_dir=root_dir)

    bb = rbfopt.RbfoptUserBlackBox(obj_funct=evaluator.estimate_cost, **config.rbfopt.user_black_box)

    # perform optimization
    rbfopt_settings = rbfopt.RbfoptSettings(**config.rbfopt.settings)

    alg = rbfopt.RbfoptAlgorithm(rbfopt_settings, bb)

    # post report to server
    evaluator.register_report(*alg.optimize())
    evaluations, report = evaluator.dump()

    # render plots
    renderer = Renderer(config, evaluations, report)
    renderer.scatterplots()
    renderer.heatmaps()
    renderer.radar()


if __name__ == "__main__":
    main()
