.PHONY: help install install-dev test clean build upload install-local build-exe

help:
	@echo "Available commands:"
	@echo "  install       - Install the package"
	@echo "  install-dev   - Install in development mode"
	@echo "  install-local - Install locally with pip"
	@echo "  test          - Run tests"
	@echo "  build-exe     - Build standalone executable"
	@echo "  clean         - Clean build artifacts"
	@echo "  build         - Build distribution packages"
	@echo "  upload        - Upload to PyPI (requires credentials)"

install:
	pip install .

install-dev:
	pip install -e .

install-local:
	pip install -e .
	@echo "✓ Installed! You can now use:"
	@echo "  audio-meta-edit /path/to/audio/file"
	@echo "  ame /path/to/audio/file"

test:
	python -m pytest tests/ -v

clean:
	rm -rf build/
	rm -rf dist/
	rm -rf *.egg-info/
	find . -type d -name __pycache__ -delete
	find . -type f -name "*.pyc" -delete

build: clean
	python setup.py sdist bdist_wheel

upload: build
	twine upload dist/*

build-exe:
	@echo "Building standalone executable..."
	python3 -m pip install pyinstaller eyed3 colorama
	pyinstaller audio-meta-edit.spec
	@echo "✓ Executable created in dist/audio-meta-edit"
	@echo "  You can now run: ./dist/audio-meta-edit /path/to/audio/file"

# Quick development setup
dev-setup:
	python -m venv venv
	@echo "Virtual environment created. Activate with:"
	@echo "  source venv/bin/activate  # Linux/Mac"
	@echo "  venv\\Scripts\\activate     # Windows"
	@echo "Then run: make install-dev" 