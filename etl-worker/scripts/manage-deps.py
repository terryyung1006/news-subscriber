#!/usr/bin/env python3
"""
Modern dependency management script for RAG Financial News Worker.

This script provides utilities for managing dependencies using UV.
"""

import argparse
import json
import subprocess
import sys
from pathlib import Path
from typing import Any, Dict, List


class DependencyManager:
    """Manages Python dependencies using UV."""

    def __init__(self):
        self.project_root = Path(__file__).parent.parent
        self.pyproject_path = self.project_root / "pyproject.toml"

        if not self.pyproject_path.exists():
            raise FileNotFoundError("pyproject.toml not found")

    def run_command(
        self, command: List[str], capture_output: bool = True
    ) -> subprocess.CompletedProcess:
        """Run a command and return the result."""
        try:
            result = subprocess.run(
                command,
                capture_output=capture_output,
                text=True,
                cwd=self.project_root,
                check=True,
            )
            return result
        except subprocess.CalledProcessError as e:
            print(f"Error running command {' '.join(command)}: {e}")
            sys.exit(1)

    def install_dependencies(self, extras: List[str] = None) -> None:
        """Install dependencies with optional extras."""
        if extras:
            cmd = ["uv", "sync", "--extra"] + extras
        else:
            cmd = ["uv", "sync"]

        print(f"Installing dependencies: {' '.join(cmd)}")
        self.run_command(cmd, capture_output=False)
        print("✅ Dependencies installed successfully")

    def update_dependencies(self) -> None:
        """Update all dependencies to latest versions."""
        print("🔄 Updating dependencies...")

        # Update lock file and sync
        self.run_command(["uv", "lock", "--upgrade"])
        self.run_command(["uv", "sync", "--dev"])

        print("✅ Dependencies updated successfully")

    def build_package(self) -> None:
        """Build the package using hatch."""
        print("🔨 Building package...")
        self.run_command(["uv", "run", "hatch", "build"])
        print("✅ Package built successfully")

    def clean_build(self) -> None:
        """Clean build artifacts."""
        print("🧹 Cleaning build artifacts...")

        # Remove build directories
        for path in ["dist", "build", "*.egg-info"]:
            for item in self.project_root.glob(path):
                if item.is_dir():
                    import shutil

                    shutil.rmtree(item)
                else:
                    item.unlink()

        # Remove Python cache
        for cache_dir in self.project_root.rglob("__pycache__"):
            import shutil

            shutil.rmtree(cache_dir)

        for pyc_file in self.project_root.rglob("*.pyc"):
            pyc_file.unlink()

        print("✅ Build artifacts cleaned")

    def check_dependencies(self) -> Dict[str, Any]:
        """Check for outdated dependencies."""
        print("🔍 Checking for outdated dependencies...")

        try:
            result = self.run_command(
                ["uv", "pip", "list", "--outdated", "--format=json"]
            )
            outdated = json.loads(result.stdout)

            if outdated:
                print("📦 Outdated dependencies found:")
                for pkg in outdated:
                    print(
                        f"  {pkg['name']}: {pkg['version']} → {pkg['latest_version']}"
                    )
            else:
                print("✅ All dependencies are up to date")

            return {"outdated": outdated}
        except subprocess.CalledProcessError:
            print("⚠️  Could not check for outdated dependencies")
            return {"outdated": []}

    def security_check(self) -> None:
        """Run security checks on dependencies."""
        print("🔒 Running security checks...")

        try:
            self.run_command(["uv", "pip", "install", "safety"])
            self.run_command(["uv", "run", "safety", "check"], capture_output=False)
        except subprocess.CalledProcessError:
            print("⚠️  Some security issues found. Check the output above.")
        else:
            print("✅ No security issues found")

    def generate_requirements(self, output_file: str = "requirements.txt") -> None:
        """Generate requirements.txt from pyproject.toml."""
        print(f"📝 Generating {output_file}...")

        # Use uv to export requirements
        self.run_command(["uv", "pip", "freeze", "--output-file", output_file])

        print(f"✅ {output_file} generated successfully")

    def lock_dependencies(self) -> None:
        """Generate lock files for reproducible builds."""
        print("🔒 Locking dependencies...")

        # Generate lock file
        self.run_command(["uv", "lock"])
        print("✅ Lock file generated")

    def show_dependency_tree(self) -> None:
        """Show dependency tree."""
        print("🌳 Dependency tree:")

        try:
            self.run_command(["uv", "pip", "install", "pipdeptree"])
            result = self.run_command(["uv", "run", "pipdeptree"])
            print(result.stdout)
        except subprocess.CalledProcessError:
            print("❌ Could not generate dependency tree")

    def validate_pyproject(self) -> None:
        """Validate pyproject.toml syntax."""
        print("✅ Validating pyproject.toml...")

        try:
            import tomllib

            with open(self.pyproject_path, "rb") as f:
                tomllib.load(f)
            print("✅ pyproject.toml is valid")
        except Exception as e:
            print(f"❌ pyproject.toml is invalid: {e}")
            sys.exit(1)


def main():
    """Main entry point."""
    parser = argparse.ArgumentParser(
        description="Modern dependency management for RAG Financial News Worker"
    )

    subparsers = parser.add_subparsers(dest="command", help="Available commands")

    # Install command
    install_parser = subparsers.add_parser("install", help="Install dependencies")
    install_parser.add_argument(
        "--extras", nargs="+", help="Extra dependencies to install (dev, test, docs)"
    )

    # Update command
    subparsers.add_parser("update", help="Update all dependencies")

    # Build command
    subparsers.add_parser("build", help="Build the package")

    # Clean command
    subparsers.add_parser("clean", help="Clean build artifacts")

    # Check command
    subparsers.add_parser("check", help="Check for outdated dependencies")

    # Security command
    subparsers.add_parser("security", help="Run security checks")

    # Generate requirements command
    req_parser = subparsers.add_parser("requirements", help="Generate requirements.txt")
    req_parser.add_argument(
        "--output", default="requirements.txt", help="Output file name"
    )

    # Lock command
    subparsers.add_parser("lock", help="Generate lock files")

    # Tree command
    subparsers.add_parser("tree", help="Show dependency tree")

    # Validate command
    subparsers.add_parser("validate", help="Validate pyproject.toml")

    args = parser.parse_args()

    if not args.command:
        parser.print_help()
        return

    try:
        manager = DependencyManager()

        if args.command == "install":
            manager.install_dependencies(extras=args.extras)
        elif args.command == "update":
            manager.update_dependencies()
        elif args.command == "build":
            manager.build_package()
        elif args.command == "clean":
            manager.clean_build()
        elif args.command == "check":
            manager.check_dependencies()
        elif args.command == "security":
            manager.security_check()
        elif args.command == "requirements":
            manager.generate_requirements(args.output)
        elif args.command == "lock":
            manager.lock_dependencies()
        elif args.command == "tree":
            manager.show_dependency_tree()
        elif args.command == "validate":
            manager.validate_pyproject()

    except Exception as e:
        print(f"❌ Error: {e}")
        sys.exit(1)


if __name__ == "__main__":
    main()
