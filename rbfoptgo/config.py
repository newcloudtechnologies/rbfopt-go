#  Copyright (c) New Cloud Technologies, Ltd. 2013-2022.
#  Author: Vitaly Isaev <vitaly.isaev@myoffice.team>
#  License: https://github.com/newcloudtechnologies/rbfopt-go/blob/master/LICENSE
import json
import os
import pathlib
from dataclasses import dataclass
from datetime import datetime
from time import mktime
from enum import Enum
from typing import List, Any, Dict

import numpy as np


@dataclass
class Bound:
    left: int
    right: int

    @staticmethod
    def from_dict(obj: Any) -> 'Bound':
        _left = int(obj.get("left"))
        _right = int(obj.get("right"))
        return Bound(_left, _right)


@dataclass
class Parameter:
    bound: Bound
    name: str

    @staticmethod
    def from_dict(obj: Any) -> 'Parameter':
        _bound = Bound.from_dict(obj.get("bound"))
        _name = str(obj.get("name"))
        return Parameter(_bound, _name)


class InvalidParameterCombinationRenderPolicy(Enum):
    omit = 1
    assign_closest_valid_value = 2


@dataclass
class PlotConfig:
    scatter_plot_policy: InvalidParameterCombinationRenderPolicy
    heatmap_render_policy: InvalidParameterCombinationRenderPolicy

    @staticmethod
    def from_dict(obj: Any) -> 'PlotConfig':
        _scatter_plot_policy = InvalidParameterCombinationRenderPolicy[str(obj.get("scatter_plot_policy"))]
        _heatmap_render_policy = InvalidParameterCombinationRenderPolicy[str(obj.get("heatmap_render_policy"))]
        return PlotConfig(_scatter_plot_policy, _heatmap_render_policy)


@dataclass
class RBFOptConfig:
    parameters: List[Parameter]
    max_evaluations: int
    max_iterations: int
    init_strategy: str
    invalid_parameter_combination_cost: int

    @staticmethod
    def from_dict(obj: Any) -> 'RBFOptConfig':
        _parameters = [Parameter.from_dict(y) for y in obj.get("parameters")]
        _max_evaluations = int(obj.get("max_evaluations"))
        _max_iterations = int(obj.get("max_iterations"))
        _init_strategy = str(obj.get("init_strategy"))
        _invalid_parameter_combination_cost = int(obj.get("invalid_parameter_combination_cost"))
        return RBFOptConfig(_parameters, _max_evaluations, _max_iterations, _init_strategy,
                            _invalid_parameter_combination_cost)

    @property
    def var_names(self) -> List[str]:
        var_names = []
        for i, param in enumerate(self.parameters):
            var_names.append(param.name)

        return var_names

    @property
    def user_black_box(self) -> Dict:
        dimensions = len(self.parameters)
        var_lower = np.zeros(shape=(dimensions,))
        var_upper = np.zeros(shape=(dimensions,))
        var_types = np.chararray(shape=(dimensions,))

        for i, param in enumerate(self.parameters):
            var_lower[i] = param.bound.left
            var_upper[i] = param.bound.right
            var_types[i] = 'I'

        return dict(
            dimension=dimensions,
            var_lower=var_lower,
            var_upper=var_upper,
            var_type=var_types,
        )

    @property
    def settings(self) -> Dict:
        return dict(
            max_evaluations=self.max_evaluations,
            max_iterations=self.max_iterations,
            rand_seed=int(mktime(datetime.now().timetuple())),
            init_strategy=self.init_strategy,
        )


@dataclass
class Config:
    root_dir: os.PathLike
    endpoint: str
    rbfopt: RBFOptConfig
    plot: PlotConfig

    @staticmethod
    def from_dict(obj: Any) -> 'Config':
        _root_dir = pathlib.Path(str(obj.get("root_dir")))
        _endpoint = str(obj.get("endpoint"))
        _rbfopt = RBFOptConfig.from_dict(obj.get("rbfopt"))
        _plot = PlotConfig.from_dict(obj.get("plot"))
        return Config(_root_dir, _endpoint, _rbfopt, _plot)

    @staticmethod
    def from_file(path: os.PathLike) -> 'Config':
        with open(path, "r") as f:
            data = json.load(f)
            return Config.from_dict(data)

# Example Usage
# jsonstring = json.loads(myjsonstring)
# root = Root.from_dict(jsonstring)
