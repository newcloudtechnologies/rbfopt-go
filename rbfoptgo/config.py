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
    """
    Describes the area of variation for a cost function parameter
    """
    left: int
    right: int

    @staticmethod
    def from_dict(obj: Any) -> 'Bound':
        """
        Constructs object from an arbitrary dictionary
        :param obj: Dictionary with parameter values
        :return: an object of desired type
        """
        _left = int(obj.get("left"))
        _right = int(obj.get("right"))
        return Bound(_left, _right)


@dataclass
class Parameter:
    """
    An arbitrary parameter of a cost function
    """
    bound: Bound
    name: str

    @staticmethod
    def from_dict(obj: Any) -> 'Parameter':
        """
        Constructs object from an arbitrary dictionary
        :param obj: Dictionary with parameter values
        :return: an object of desired type
        """
        _bound = Bound.from_dict(obj.get("bound"))
        _name = str(obj.get("name"))
        return Parameter(_bound, _name)


class InvalidParameterCombinationRenderPolicy(Enum):
    """
    Describes the behavior of the plot renderer when the computed cost function value
    corresponds to ErrInvalidParameterCombinationCost.
    """
    omit = 1
    assign_closest_valid_value = 2


@dataclass
class PlotConfig:
    """
    Configuration of plotting config.
    """
    scatter_plot_policy: InvalidParameterCombinationRenderPolicy
    heatmap_render_policy: InvalidParameterCombinationRenderPolicy

    @staticmethod
    def from_dict(obj: Any) -> 'PlotConfig':
        """
        Constructs object from an arbitrary dictionary
        :param obj: Dictionary with parameter values
        :return: an object of desired type
        """
        _scatter_plot_policy = InvalidParameterCombinationRenderPolicy[str(obj.get("scatter_plot_policy"))]
        _heatmap_render_policy = InvalidParameterCombinationRenderPolicy[str(obj.get("heatmap_render_policy"))]
        return PlotConfig(_scatter_plot_policy, _heatmap_render_policy)


@dataclass
class RBFOptConfig:
    """
    Configuration of RBFOpt library.
    """
    parameters: List[Parameter]
    max_evaluations: int
    max_iterations: int
    init_strategy: str
    invalid_parameter_combination_cost: int

    @staticmethod
    def from_dict(obj: Any) -> 'RBFOptConfig':
        """
        Constructs object from an arbitrary dictionary
        :param obj: Dictionary with parameter values
        :return: an object of desired type
        """
        _parameters = [Parameter.from_dict(y) for y in obj.get("parameters")]
        _max_evaluations = int(obj.get("max_evaluations"))
        _max_iterations = int(obj.get("max_iterations"))
        _init_strategy = str(obj.get("init_strategy"))
        _invalid_parameter_combination_cost = int(obj.get("invalid_parameter_combination_cost"))
        return RBFOptConfig(_parameters, _max_evaluations, _max_iterations, _init_strategy,
                            _invalid_parameter_combination_cost)

    @property
    def var_names(self) -> List[str]:
        """
        Returns parameter names suitable for RBFOpt
        :return:  list of parameter namse
        """
        var_names = []
        for param in self.parameters:
            var_names.append(param.name)

        return var_names

    @property
    def user_black_box(self) -> Dict:
        """
        Returns the area of the function scope where to perform optimization.
        :return:  list of parameter namse
        """
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
        """
        Returns dictionary with RBFOpt settings.
        :return: Dict with settings.
        """
        return dict(
            max_evaluations=self.max_evaluations,
            max_iterations=self.max_iterations,
            rand_seed=int(mktime(datetime.now().timetuple())),
            init_strategy=self.init_strategy,
        )


@dataclass
class Config:
    """
    A configuration for the Python part of rbfopt-go
    """
    root_dir: os.PathLike
    endpoint: str
    rbfopt: RBFOptConfig
    plot: PlotConfig

    @staticmethod
    def from_dict(obj: Any) -> 'Config':
        """
        Constructs object from an arbitrary dictionary
        :param obj: Dictionary with parameter values
        :return: an object of desired type
        """
        _root_dir = pathlib.Path(str(obj.get("root_dir")))
        _endpoint = str(obj.get("endpoint"))
        _rbfopt = RBFOptConfig.from_dict(obj.get("rbfopt"))
        _plot = PlotConfig.from_dict(obj.get("plot"))
        return Config(_root_dir, _endpoint, _rbfopt, _plot)

    @staticmethod
    def from_file(path: os.PathLike) -> 'Config':
        """
        Loads configuration from file
        :param path: full path to config file
        :return:
        """
        with open(path, "r") as f:
            data = json.load(f)
            return Config.from_dict(data)
