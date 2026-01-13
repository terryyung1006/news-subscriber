#!/usr/bin/env python3
"""
Test runner for RAG Financial News Worker.

This script runs all tests in the project to ensure everything works correctly.
"""

import os
import subprocess
import sys
from pathlib import Path


def run_command(command, description):
    """Run a command and handle errors."""
    print(f"\n{'='*60}")
    print(f"Running: {description}")
    print(f"Command: {command}")
    print("=" * 60)

    try:
        result = subprocess.run(
            command, shell=True, check=True, capture_output=True, text=True
        )
        print("✅ SUCCESS")
        if result.stdout:
            print("Output:")
            print(result.stdout)
        return True
    except subprocess.CalledProcessError as e:
        print("❌ FAILED")
        print(f"Error code: {e.returncode}")
        if e.stdout:
            print("Stdout:")
            print(e.stdout)
        if e.stderr:
            print("Stderr:")
            print(e.stderr)
        return False


def main():
    """Run all tests."""
    print("🧪 RAG Financial News Worker - Test Suite")
    print("=" * 60)

    # Change to project root
    project_root = Path(__file__).parent
    os.chdir(project_root)

    # Test results
    results = []

    # 1. Check if required files exist
    print("\n📁 Checking project structure...")
    required_files = [
        "requirements.txt",
        "config/config.yaml",
        "src/main.py",
        "docker-compose.yml",
        "docker/Dockerfile",
    ]

    for file_path in required_files:
        if Path(file_path).exists():
            print(f"✅ {file_path}")
        else:
            print(f"❌ {file_path} - MISSING")
            results.append(False)

    # 2. Run linting
    results.append(
        run_command(
            "python -m flake8 src/ --max-line-length=100 --ignore=E501,W503",
            "Code linting with flake8",
        )
    )

    # 3. Run type checking
    results.append(
        run_command(
            "python -m mypy src/ --ignore-missing-imports", "Type checking with mypy"
        )
    )

    # 4. Run unit tests
    results.append(run_command("python -m pytest tests/ -v", "Unit tests with pytest"))

    # 5. Test configuration loading
    results.append(
        run_command(
            "python -c \"import yaml; yaml.safe_load(open('config/config.yaml'))\"",
            "Configuration file validation",
        )
    )

    # 6. Test Docker build (without running)
    results.append(
        run_command(
            "docker build -f docker/Dockerfile -t rag-worker-test .",
            "Docker build test",
        )
    )

    # Summary
    print(f"\n{'='*60}")
    print("📊 TEST SUMMARY")
    print("=" * 60)

    passed = sum(results)
    total = len(results)

    print(f"Tests passed: {passed}/{total}")
    print(f"Success rate: {(passed/total)*100:.1f}%")

    if passed == total:
        print("🎉 All tests passed!")
        return 0
    else:
        print("❌ Some tests failed!")
        return 1


if __name__ == "__main__":
    sys.exit(main())
