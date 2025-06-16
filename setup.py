#!/usr/bin/env python3
"""
Setup script for the audio metadata editor CLI tool.
"""

from setuptools import setup, find_packages
import os

# Read the README file for long description
def read_readme():
    readme_path = os.path.join(os.path.dirname(__file__), 'README.md')
    if os.path.exists(readme_path):
        with open(readme_path, 'r', encoding='utf-8') as f:
            return f.read()
    return ""

# Read version from package
def get_version():
    version_file = os.path.join('eyed3_cli_editor', '__init__.py')
    with open(version_file, 'r', encoding='utf-8') as f:
        for line in f:
            if line.startswith('__version__'):
                return line.split('=')[1].strip().strip('"\'')
    return "0.1.0"

setup(
    name="audio-metadata-editor",
    version=get_version(),
    author="Stefano Frigerio",
    author_email="stefri3@gmail.com",
    description="Interactive CLI tool to edit audio file metadata",
    long_description=read_readme(),
    long_description_content_type="text/markdown",
    url="https://github.com/stfrigerio/audio-metadata-editor",
    packages=find_packages(),
    classifiers=[
        "Development Status :: 4 - Beta",
        "Intended Audience :: End Users/Desktop",
        "License :: OSI Approved :: MIT License",
        "Operating System :: OS Independent",
        "Programming Language :: Python :: 3",
        "Programming Language :: Python :: 3.7",
        "Programming Language :: Python :: 3.8",
        "Programming Language :: Python :: 3.9",
        "Programming Language :: Python :: 3.10",
        "Programming Language :: Python :: 3.11",
        "Topic :: Multimedia :: Sound/Audio",
        "Topic :: Utilities",
        "Environment :: Console",
    ],
    python_requires=">=3.7",
    install_requires=[
        "eyed3>=0.9.0",
        "colorama>=0.4.0",
    ],
    entry_points={
        "console_scripts": [
            "audio-meta-edit=eyed3_cli_editor.cli:cli_entry_point",
            "ame=eyed3_cli_editor.cli:cli_entry_point",  # Short alias
        ],
    },
    keywords="audio metadata mp3 id3 music editor cli",
    project_urls={
        "Bug Reports": "https://github.com/yourusername/audio-metadata-editor/issues",
        "Source": "https://github.com/yourusername/audio-metadata-editor",
    },
) 