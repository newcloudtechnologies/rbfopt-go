#  Copyright (c) New Cloud Technologies, Ltd. 2013-2022.
#  Author: Vitaly Isaev <vitaly.isaev@myoffice.team>
#  License: https://github.com/newcloudtechnologies/rbfopt-go/blob/master/LICENSE

import pathlib
from typing import List

import numpy as np
import pandas as pd

from rbfoptgo.common import Cost, ParameterValue
from rbfoptgo.client import Client
from rbfoptgo.config import Config
from rbfoptgo.report import Report
from rbfoptgo import names


class Evaluator:
    """
    Evaluator is responsible for cost function estimation. It performs HTTP calls to the Golang part of a library.
    """
    __config: Config
    __client: Client
    __parameter_names: List[str]
    __evaluations: []
    __root_dir: pathlib.Path
    __report: Report
    __iterations: int

    def __init__(self, config: Config, client: Client, parameter_names: List[str], root_dir: pathlib.Path):
        self.__config = config
        self.__client = client
        self.__parameter_names = parameter_names
        self.__evaluations = []
        self.__root_dir = root_dir
        self.__iterations = 0

    def __np_array_to_parameter_values(self, raw_values: np.ndarray) -> List[ParameterValue]:
        parameter_values = []
        for i, raw_value in enumerate(raw_values):
            parameter_values.append(
                ParameterValue(name=self.__parameter_names[i], value=int(raw_value)),
            )
        return parameter_values

    def estimate_cost(self, raw_values: np.ndarray) -> Cost:
        """
        :param raw_values: vector of cost function arguments
        :return: cost function particular value
        """
        self.__iterations += 1

        parameter_values = self.__np_array_to_parameter_values(raw_values)
        cost, invalid_parameter_combination = self.__client.estimate_cost(parameter_values)

        # store evaluation result for the future use
        entry = dict(zip(self.__parameter_names, raw_values))
        entry[names.Iteration] = self.__iterations
        entry[names.Cost] = cost
        entry[names.InvalidParameterCombination] = invalid_parameter_combination

        self.__evaluations.append(entry)

        return cost

    def register_report(
            self,
            cost: Cost,
            optimum: np.ndarray,
            iterations: int,
            evaluations: int,
            fast_evaluations: int,
    ):
        """
        Registers final report.
        :param cost: minimal value of a cost function (the optimum)
        :param optimum: the values of cost function arguments corresponding to the optimum
        :param iterations: number of iterations performed
        :param evaluations: number of evaluations performed
        :param fast_evaluations: number of fast_evaluations
        :return: None
        """
        report = Report(
            bounds=self.__config.rbfopt.parameters,
            cost=cost,
            optimum=self.__np_array_to_parameter_values(optimum),
            iterations=iterations,
            evaluations=evaluations,
            fast_evaluations=fast_evaluations,
        )

        self.__client.register_report(report)
        self.__report = report

    def dump(self) -> (pd.DataFrame, Report):
        """
        Returns the list of performed cost function evaluations.
        :return: DataFrame + json-serializable Report compatible with Golang library
        """
        # dump history of evaluations for future usage
        evaluations = pd.DataFrame(self.__evaluations)
        file_path = self.__root_dir.joinpath("evaluations.csv")
        evaluations.to_csv(file_path, header=True, index=False)

        # dump report for future usage
        file_path = self.__root_dir.joinpath("report.json")
        self.__report.save_to_file(file_path)

        return evaluations, self.__report
