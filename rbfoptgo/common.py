#  Copyright (c) New Cloud Technologies, Ltd. 2013-2022.
#  Author: Vitaly Isaev <vitaly.isaev@myoffice.team>
#  License: https://github.com/newcloudtechnologies/rbfopt-go/blob/master/LICENSE

from dataclasses import dataclass

Cost = float


@dataclass
class ParameterValue:
    """
    Represents value of a particular parameter
    """
    name: str
    value: int
