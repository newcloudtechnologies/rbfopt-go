#  Copyright (c) New Cloud Technologies, Ltd. 2013-2022.
#  Author: Vitaly Isaev <vitaly.isaev@myoffice.team>
#  License: https://github.com/newcloudtechnologies/rbfopt-go/blob/master/LICENSE

import functools
import typing

import matplotlib.axes
import matplotlib.image
import matplotlib.pyplot as plt
import matplotlib.ticker
import numpy as np
import pandas as pd
import scipy.interpolate
import scipy.stats
from adjustText import adjust_text
from colorhash import ColorHash

from rbfoptgo import names
from rbfoptgo.config import Config, InvalidParameterCombinationRenderPolicy
from rbfoptgo.report import Report


class Renderer:
    __df: pd.DataFrame
    __report: Report
    __config: Config
    __transparent: bool

    def __init__(self, config: Config, df: pd.DataFrame, report: Report, transparent: bool = False):
        self.__df = df.loc[:, df.columns != names.Iteration]
        self.__report = report
        self.__config = config
        self.__transparent = transparent

    def __prepare_df(self, policy: InvalidParameterCombinationRenderPolicy) -> pd.DataFrame:
        match policy:
            case InvalidParameterCombinationRenderPolicy.omit:
                # Filter values corresponding to ErrInvalidParameterCombination params
                return self.__df[self.__df.invalid_parameter_combination == False]
            case InvalidParameterCombinationRenderPolicy.assign_closest_valid_value:
                # avoid white spots on the heatmap
                unique_costs = self.__df.cost.unique()
                unique_costs.sort()
                # -1 stands for ErrInvalidParameterCombinationCost, take -2 - the next closest
                closest_valid_value = unique_costs[-2]

                # make df with replaced values
                df = self.__df.copy()
                mask = df.invalid_parameter_combination == True
                df.cost.where(~mask, closest_valid_value, inplace=True)
                return df
            case _:
                raise ValueError(f"unknown policy: {policy}")

    @property
    @functools.cache
    def __parameter_column_names(self) -> typing.List[str]:
        utility_columns = (names.Cost, names.InvalidParameterCombination)
        return list(filter(lambda x: x not in utility_columns, self.__df.columns))

    def __cost_bounds(self, df: pd.DataFrame) -> typing.List[str]:
        cost = df[names.Cost]
        return [cost.min(), cost.max()]

    def scatterplots(self):
        self.__render_scatterplot_group(only_optimal_values=False)
        self.__render_scatterplot_group(only_optimal_values=True)

    def __render_scatterplot_group(self, only_optimal_values: bool):
        column_names = self.__parameter_column_names

        if len(column_names) <= 2:
            n_rows, n_columns = 1, len(column_names)
        else:
            n_columns = 2
            if len(column_names) % 2 == 0:
                n_rows = int(len(column_names) / n_columns)
            else:
                n_rows = int(len(column_names) / n_columns) + 1

        # this multiplier is purely empirical...
        figsize = (6 * n_rows, 6 * n_columns)

        fig, axes = plt.subplots(nrows=n_rows, ncols=n_columns, figsize=figsize,
                                 squeeze=False, constrained_layout=True)
        axes = axes.flat

        for i in range(len(axes)):
            if i < len(column_names):
                self.__render_scatterplot(axes[i], column_names[i], only_optimal_values=only_optimal_values)
            else:
                axes[i].axis('off')

        fig.tight_layout()

        suffix = "only_optimal_values" if only_optimal_values else "all_values"
        figure_path = self.__config.root_dir.joinpath(f"scatterplot_{suffix}.png")
        fig.savefig(figure_path, transparent=self.__transparent)

    def __render_scatterplot(self, ax: matplotlib.axes.Axes, col_name: str, only_optimal_values: bool):
        df = self.__prepare_df(policy=self.__config.plot.scatter_plot_policy)
        data = pd.DataFrame({col_name: df[col_name], names.Cost: df[names.Cost]})

        if only_optimal_values:
            # for every argument value, pick the best cost function value
            data = data.groupby(col_name)[names.Cost].agg(lambda x: x.min()).reset_index()

        color = ColorHash(col_name).hex
        ax.plot(data[col_name], data[names.Cost], linewidth=0, marker='o', color=color)

        ax.set_xlabel(col_name, fontsize=14)
        ax.set_ylabel('Cost function', fontsize=14)

        # draw point with optimum
        opt_arg = self.__report.optimum_argument(col_name)
        opt_val = self.__report.cost
        ax.scatter(opt_arg, opt_val, color='red', marker='o', s=100)
        ax.annotate("{:.2f}".format(opt_val), (opt_arg, opt_val))

        # set equal limits
        (cost_min, cost_max) = self.__cost_bounds(df)
        # ax.set_ybound(lower=cost_min, upper=cost_max)
        ax.set_ylim(bottom=cost_min, top=cost_max)

    def heatmaps(self):
        df = self.__prepare_df(policy=self.__config.plot.heatmap_render_policy)

        # TODO: one can pass a particular set of interpolation methods,
        #  but honestly I can't see any significant difference between them.
        # methods = ['none', 'nearest', 'bilinear', 'bicubic', 'spline16',
        #            'spline36', 'hanning', 'hamming', 'hermite', 'kaiser', 'quadric',
        #            'catrom', 'gaussian', 'bessel', 'mitchell', 'sinc', 'lanczos']
        methods = ['hamming']

        for method in methods:
            print(f"rendering heatmap matrix using method {method}")
            self.__pairwise_heatmaps(df=df, interpolation=method)
            self.__pairwise_heatmap_matrix(df=df, interpolation=method)

    def __pairwise_heatmaps(self, df: pd.DataFrame, interpolation: str):
        column_names = self.__parameter_column_names
        for i in range(len(column_names) - 1):
            for j in range(i + 1, len(column_names)):
                self.__pairwise_heatmap(df=df,
                                        col_name_1=column_names[i],
                                        col_name_2=column_names[j],
                                        interpolation=interpolation)

    def __pairwise_heatmap(self,
                           df: pd.DataFrame,
                           col_name_1: str,
                           col_name_2: str,
                           interpolation: str,
                           ):

        figsize = (6, 6)

        fig, ax = plt.subplots(figsize=figsize, constrained_layout=True)

        im = self.__pairwise_heatmap_interpolate(
            df=df,
            ax=ax,
            col_name_1=col_name_1,
            col_name_2=col_name_2,
            interpolation=interpolation,
            x_label=True,
            y_label=True,
        )

        fig.colorbar(im, ax=ax, shrink=0.6)

        figure_path = self.__config.root_dir.joinpath(f"heatmap_{col_name_1}_{col_name_2}_{interpolation}.png")
        fig.savefig(figure_path, transparent=self.__transparent)

    def __pairwise_heatmap_matrix(self, df: pd.DataFrame, interpolation: str):
        column_names = self.__parameter_column_names

        # approximate size that make image look well
        figsize = (3 * len(column_names), 3 * len(column_names))

        fig, axes = plt.subplots(nrows=len(column_names) - 1,
                                 ncols=len(column_names) - 1,
                                 figsize=figsize,
                                 constrained_layout=True,
                                 )
        im = None
        for i in range(len(column_names) - 1):
            for j in range(0, i):
                axes[j, i].axis('off')
            for j in range(i + 1, len(column_names)):
                col_name_1, col_name_2 = column_names[i], column_names[j]
                ax = axes[j - 1, i]
                im = self.__pairwise_heatmap_interpolate(
                    df=df,
                    ax=ax,
                    col_name_1=col_name_1,
                    col_name_2=col_name_2,
                    interpolation=interpolation,
                    x_label=j == len(column_names) - 1,
                    y_label=i == 0,
                )

        cbar = fig.colorbar(im, ax=axes, shrink=0.6)
        cbar.ax.tick_params(labelsize=24)

        figure_path = self.__config.root_dir.joinpath(f"heatmap_matrix_{interpolation}.png")
        fig.savefig(figure_path, transparent=self.__transparent, dpi=300)

    def __pairwise_heatmap_interpolate(self,
                                       df: pd.DataFrame,
                                       ax: matplotlib.axes.Axes,
                                       col_name_1: str,
                                       col_name_2: str,
                                       interpolation: str,
                                       x_label: bool,
                                       y_label: bool,
                                       ) -> matplotlib.image.AxesImage:
        data = df[[col_name_1, col_name_2, names.Cost]]

        # select the minimums
        data = data.groupby([col_name_1, col_name_2])[names.Cost].agg(lambda x: x.min()).reset_index()

        # check if data is sufficient
        if data.shape[0] < 4:
            raise ValueError('Too little data to render grid, try to increase number of iterations')

        # compute grid bounds
        xs0 = data[col_name_1]
        ys0 = data[col_name_2]
        x_min, x_max = xs0.min(), xs0.max()
        y_min, y_max = ys0.min(), ys0.max()
        extent = [x_min, x_max, y_min, y_max]
        N = 100j
        xs, ys = np.mgrid[x_min:x_max:N, y_min:y_max:N]

        zs0 = data[names.Cost]

        # interpolate data
        resampled = scipy.interpolate.griddata(
            points=(xs0, ys0),
            values=zs0,
            xi=(xs, ys),
            method='cubic',
        )

        # render interpolated grid
        (cost_min, cost_max) = self.__cost_bounds(df)
        im = ax.imshow(
            resampled.T,
            cmap='jet',
            origin='lower',
            interpolation=interpolation,
            vmin=cost_min,
            vmax=cost_max,
            extent=extent,
        )

        # draw point with optimum
        opt_x, opt_y, opt_val = self.__derive_optimum_coordinates(col_name_1, col_name_2, )

        t = ax.text(opt_x, opt_y, opt_val, ha='center', va='center', fontsize=20)
        t.set_bbox(dict(facecolor='white', alpha=0.75, edgecolor='red'))

        adjust_text(
            texts=[t],
            ax=ax,
            arrowprops=dict(arrowstyle='->', color='red'),
            expand_text=(2, 2),
            expand_align=(2, 2),
            expand_objects=(2, 2),
            expand_points=(2, 2),
        )
        ax.scatter(opt_x, opt_y, color='red', marker='o', s=100)

        # borders
        x_gap = (x_max - x_min) * 0.025
        y_gap = (y_max - y_min) * 0.025
        ax.set_xlim(left=x_min-x_gap, right=x_max+x_gap)
        ax.set_ylim(bottom=y_min-y_gap, top=y_max+y_gap)

        # assign axes labels
        ax.tick_params(axis='x', which='major', labelsize=16)
        ax.tick_params(axis='y', which='major', labelsize=16)
        if x_label:
            ax.set_xlabel(col_name_1, fontsize=18)
        if y_label:
            ax.set_ylabel(col_name_2, fontsize=18)

        return im

    def __derive_optimum_coordinates(self,
                                     col_name_1: str, col_name_2: str,
                                     ) -> (float, float, float):
        col_val_1 = self.__report.optimum_argument(col_name_1)
        col_val_2 = self.__report.optimum_argument(col_name_2)
        cost_val = self.__report.cost

        # There must be a much better way to achieve that, PRs are welcome
        cost_val_abs = np.abs(cost_val)
        if cost_val_abs > 100:
            cost_val_abs = int(cost_val)
        if cost_val_abs > 10:
            cost_val = np.around(cost_val, 1)
        if cost_val_abs > 1:
            cost_val = np.around(cost_val, 2)
        if cost_val_abs > 0.1:
            cost_val = np.around(cost_val, 3)
        if cost_val_abs > 0.01:
            cost_val = np.around(cost_val, 4)
        if cost_val_abs > 0.001:
            cost_val = np.around(cost_val, 5)

        return col_val_1, col_val_2, cost_val

    def radar(self):
        # convert optimum values to df
        df = pd.DataFrame(self.__report.optimum).T
        df.columns = df.iloc[0]
        df.drop(df.index[0], inplace=True)

        # explanation: https://www.python-graph-gallery.com/390-basic-radar-chart
        values = df.iloc[0].values.flatten().tolist()
        values += values[:1]

        N = len(df.iloc[0])
        angles = [n / float(N) * 2 * np.pi for n in range(N)]
        angles += angles[:1]

        ax = plt.subplot(111, polar=True)
        ax.set_xticks(angles[:-1], df.columns, size=14)

        ax.set_rlabel_position(0)

        ax.plot(angles, values, linewidth=1, linestyle='solid')
        ax.fill(angles, values, 'b', alpha=0.1)

        plt.savefig(self.__config.root_dir.joinpath('polar.png'), transparent=self.__transparent)
