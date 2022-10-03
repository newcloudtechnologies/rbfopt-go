#  Copyright (c) New Cloud Technologies, Ltd. 2013-2022.
#  Author: Vitaly Isaev <vitaly.isaev@myoffice.team>
#  License: https://github.com/newcloudtechnologies/rbfopt-go/blob/master/LICENSE

from http import HTTPStatus
from typing import List
from urllib.parse import urljoin

import jsons
import requests

from rbfoptgo.common import Cost, ParameterValue
from rbfoptgo.report import Report
from rbfoptgo import names


class Client:
    """
    HTTP client to the Go part of library
    """
    url_head: str
    session: requests.Session

    def __init__(self, endpoint: str):
        self.url_head = f'http://{endpoint}'
        self.session = requests.Session()

    def estimate_cost(self, parameter_values: List[ParameterValue]) -> (Cost, bool):
        """
        Requests a particular cost function value for a given parameters
        :param parameter_values: a vector of parameters
        :return: 1. The value of a cost function
               2. Sign that optimizer gave a vector of parameters that was considered as non-optimal
        """
        print(f"request '{parameter_values}'")

        payload = dict(parameter_values=parameter_values)
        response = self.session.get(
            urljoin(self.url_head, 'estimate_cost'),
            json=jsons.dump(payload),
        )

        print(f"response code={response.status_code} body={response.json()}")

        if response.status_code != HTTPStatus.OK:
            raise ValueError(f'invalid status code {response.status_code}')

        return response.json()[names.Cost], response.json()[names.InvalidParameterCombination]

    def register_report(self, report: Report):
        """
        Posts the report to the Go sides
        :param report: optimizer report itself
        :return:
        """
        print(f"request '{report}'")

        payload = dict(report=report)
        response = self.session.post(
            urljoin(self.url_head, 'register_report'),
            json=jsons.dump(payload),
        )

        print(f"response code={response.status_code} response={response}")

        if response.status_code != HTTPStatus.OK:
            raise ValueError(f'invalid status code {response.status_code}')
