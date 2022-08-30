#  Copyright (c) New Cloud Technologies, Ltd. 2013-2022.
#  Author: Vitaly Isaev <vitaly.isaev@myoffice.team>
#  License: https://github.com/newcloudtechnologies/rbfopt-go/blob/master/LICENSE

from setuptools import setup, find_packages
import pathlib

HERE = pathlib.Path(__file__).parent

setup(name='rbfopt-go',
      # version_config={
      #     "dev_template": "{tag}",
      # },
      setuptools_git_versioning={
          "enabled": True,
          "template": "{tag}",
          "dirty_template": "{tag}",
      },
      description='Find better configuration of your Go service with derivative-free optimization algorithms',
      author='Vitaly Isaev',
      author_email='vitalyisaev2@gmail.com',
      url='https://github.com/newcloudtechnologies/rbfopt-go',
      packages=find_packages(),
      setup_requires=["setuptools-git-versioning"],
      install_requires=(
          "jsons",
          "numpy",
          "Pyomo==6.1.2",
          "rbfopt==4.2.2",
          "requests",
          "urllib3",
          "pandas",
          "matplotlib",
          "scipy",
          "colorhash==1.0.4"
      ),
      license_file="LICENSE",
      package_dir={
          '': '.'
      },
      classifiers=[
          "Programming Language :: Python :: 3",
          "License :: OSI Approved :: MIT License",
          "Operating System :: OS Independent",
          "Topic :: Scientific/Engineering :: Mathematics",
      ],
      entry_points={
          'console_scripts': [
              'rbfopt-go-wrapper = rbfoptgo.main:main',
          ]
      },
      zip_safe=False,
      python_requires=">=3.7",
      )
