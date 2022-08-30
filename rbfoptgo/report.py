#  Copyright (c) New Cloud Technologies, Ltd. 2013-2022.
#  Author: Vitaly Isaev <vitaly.isaev@myoffice.team>
#  License: https://github.com/newcloudtechnologies/rbfopt-go/blob/master/LICENSE

import json
import os
from dataclasses import dataclass
from typing import List

import jsons

from rbfoptgo.common import Cost, ParameterValue
from rbfoptgo.config import Parameter


@dataclass
class Report:
    bounds: List[Parameter]
    optimum: List[ParameterValue]
    cost: Cost
    iterations: int
    evaluations: int
    fast_evaluations: int

    def optimum_argument(self, name: str) -> int:
        for pv in self.optimum:
            if pv.name == name:
                return pv.value

        raise ValueError(f"unexpected name {name}")

    def save_to_file(self, file_path: os.PathLike):
        with open(file_path, "w") as f:
            obj = jsons.dump(self)
            json.dump(obj, f)

    @classmethod
    def load_from_file(cls, file_path: os.PathLike) -> 'Report':
        with open(file_path, "r") as f:
            data = json.load(f)
            return jsons.load(data, cls)
