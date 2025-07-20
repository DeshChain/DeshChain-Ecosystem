#!/usr/bin/env python3

from setuptools import setup, find_packages

with open("README.md", "r", encoding="utf-8") as fh:
    long_description = fh.read()

with open("requirements.txt", "r", encoding="utf-8") as fh:
    requirements = [line.strip() for line in fh if line.strip() and not line.startswith("#")]

setup(
    name="deshchain-sdk",
    version="1.0.0",
    author="DeshChain Development Team",
    author_email="dev@deshchain.network",
    description="Official Python SDK for DeshChain blockchain",
    long_description=long_description,
    long_description_content_type="text/markdown",
    url="https://github.com/deshchain/deshchain-sdk-python",
    project_urls={
        "Bug Tracker": "https://github.com/deshchain/deshchain-sdk-python/issues",
        "Documentation": "https://docs.deshchain.network/sdk/python",
        "Source Code": "https://github.com/deshchain/deshchain-sdk-python",
    },
    packages=find_packages(exclude=["tests", "tests.*"]),
    classifiers=[
        "Development Status :: 5 - Production/Stable",
        "Intended Audience :: Developers",
        "License :: OSI Approved :: MIT License",
        "Operating System :: OS Independent",
        "Programming Language :: Python :: 3",
        "Programming Language :: Python :: 3.8",
        "Programming Language :: Python :: 3.9",
        "Programming Language :: Python :: 3.10",
        "Programming Language :: Python :: 3.11",
        "Programming Language :: Python :: 3.12",
        "Topic :: Software Development :: Libraries :: Python Modules",
        "Topic :: Office/Business :: Financial",
        "Topic :: System :: Distributed Computing",
    ],
    python_requires=">=3.8",
    install_requires=requirements,
    extras_require={
        "dev": [
            "pytest>=7.0.0",
            "pytest-asyncio>=0.21.0",
            "pytest-cov>=4.0.0",
            "black>=23.0.0",
            "isort>=5.12.0",
            "flake8>=6.0.0",
            "mypy>=1.0.0",
            "sphinx>=6.0.0",
            "sphinx-rtd-theme>=1.2.0",
        ],
        "test": [
            "pytest>=7.0.0",
            "pytest-asyncio>=0.21.0",
            "pytest-cov>=4.0.0",
            "responses>=0.23.0",
        ],
    },
    entry_points={
        "console_scripts": [
            "deshchain-cli=deshchain.cli:main",
        ],
    },
    keywords=[
        "deshchain",
        "blockchain",
        "cosmos",
        "cultural",
        "heritage",
        "defi",
        "lending",
        "sdk",
        "python",
        "india",
    ],
    package_data={
        "deshchain": ["py.typed"],
    },
    include_package_data=True,
    zip_safe=False,
)