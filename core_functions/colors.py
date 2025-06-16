"""
Color utilities for the CLI tool.
"""

from colorama import Fore, Back, Style, init
import sys

# Initialize colorama for cross-platform colored output
init(autoreset=True)

class Colors:
    """Color constants and utility functions."""
    
    # Basic colors
    RED = Fore.RED
    GREEN = Fore.GREEN
    YELLOW = Fore.YELLOW
    BLUE = Fore.BLUE
    MAGENTA = Fore.MAGENTA
    CYAN = Fore.CYAN
    WHITE = Fore.WHITE
    
    # Styles
    BRIGHT = Style.BRIGHT
    DIM = Style.DIM
    RESET = Style.RESET_ALL
    
    # Background colors
    BG_RED = Back.RED
    BG_GREEN = Back.GREEN
    BG_YELLOW = Back.YELLOW
    BG_BLUE = Back.BLUE

def colorize(text: str, color: str, style: str = "") -> str:
    """Apply color and style to text."""
    return f"{style}{color}{text}{Style.RESET_ALL}"

def success(text: str) -> str:
    """Format success message in green."""
    return colorize(text, Colors.GREEN, Colors.BRIGHT)

def error(text: str) -> str:
    """Format error message in red."""
    return colorize(text, Colors.RED, Colors.BRIGHT)

def warning(text: str) -> str:
    """Format warning message in yellow."""
    return colorize(text, Colors.YELLOW, Colors.BRIGHT)

def info(text: str) -> str:
    """Format info message in blue."""
    return colorize(text, Colors.BLUE, Colors.BRIGHT)

def highlight(text: str) -> str:
    """Format highlighted text in cyan."""
    return colorize(text, Colors.CYAN, Colors.BRIGHT)

def dim(text: str) -> str:
    """Format dimmed text."""
    return colorize(text, Colors.WHITE, Colors.DIM)

def print_success(text: str):
    """Print success message."""
    print(success(text))

def print_error(text: str):
    """Print error message."""
    print(error(text), file=sys.stderr)

def print_warning(text: str):
    """Print warning message."""
    print(warning(text))

def print_info(text: str):
    """Print info message."""
    print(info(text))

def print_highlight(text: str):
    """Print highlighted text."""
    print(highlight(text))

def print_header(text: str):
    """Print a header with styling."""
    separator = "─" * len(text)
    print(f"\n{colorize(separator, Colors.CYAN)}")
    print(f"{colorize(text, Colors.CYAN, Colors.BRIGHT)}")
    print(f"{colorize(separator, Colors.CYAN)}")

def print_section(text: str):
    """Print a section header."""
    print(f"\n{colorize('●', Colors.GREEN, Colors.BRIGHT)} {colorize(text, Colors.WHITE, Colors.BRIGHT)}") 