#!/usr/bin/env python3
"""
Configuration management for the ETL Worker application.

This module loads config.yaml as application state instead of environment variables.
"""

import os
from pathlib import Path
from typing import Any, Dict, Optional

import yaml


class Config:
    """
    Centralized configuration management for the ETL Worker application.
    
    This class loads configuration from YAML files and provides type-safe access
    to configuration values throughout the application.
    """
    
    _instance: Optional['Config'] = None
    _config_data: Dict[str, Any] = {}
    
    def __new__(cls) -> 'Config':
        """Singleton pattern to ensure single configuration instance."""
        if cls._instance is None:
            cls._instance = super().__new__(cls)
        return cls._instance
    
    def __init__(self) -> None:
        """Initialize configuration if not already loaded."""
        if not self._config_data:
            self._load_config()
    
    def _load_config(self) -> None:
        """Load configuration from YAML file."""
        # Find config.yaml relative to this file
        config_path = Path(__file__).parent.parent / "config" / "config.yaml"
        
        if not config_path.exists():
            print(f"Config file not found: {config_path}")
            raise FileNotFoundError(f"Configuration file not found: {config_path}")
        
        try:
            with open(config_path, "r") as f:
                self._config_data = yaml.safe_load(f) or {}
            
            print(f"Loaded configuration from {config_path}")
            print(f"# Loaded {len(self._config_data)} configuration values")
            
        except Exception as e:
            print(f"Failed to load configuration: {e}")
            raise
    
    def get(self, key: str, default: Any = None) -> Any:
        """
        Get a configuration value by key.
        
        Args:
            key: Configuration key
            default: Default value if key not found
            
        Returns:
            Configuration value or default
        """
        return self._config_data.get(key, default)
    
    def get_string(self, key: str, default: str = "") -> str:
        """Get a configuration value as string."""
        value = self.get(key, default)
        return str(value) if value is not None else default
    
    def get_int(self, key: str, default: int = 0) -> int:
        """Get a configuration value as integer."""
        value = self.get(key, default)
        try:
            return int(value) if value is not None else default
        except (ValueError, TypeError):
            print(f"Warning: Invalid integer value for {key}: {value}, using default: {default}")
            return default
    
    def get_bool(self, key: str, default: bool = False) -> bool:
        """Get a configuration value as boolean."""
        value = self.get(key, default)
        if isinstance(value, bool):
            return value
        if isinstance(value, str):
            return value.lower() in ('true', '1', 'yes', 'on')
        return bool(value) if value is not None else default
    
    def get_list(self, key: str, default: Optional[list] = None) -> list:
        """Get a configuration value as list."""
        value = self.get(key, default)
        if isinstance(value, list):
            return value
        if value is None:
            return default or []
        return [value]  # Convert single value to list
    
    def reload(self) -> None:
        """Reload configuration from file."""
        self._config_data.clear()
        self._load_config()
        print("Configuration reloaded")
    
    def to_dict(self) -> Dict[str, Any]:
        """Get configuration as dictionary."""
        return self._config_data.copy()


# Global configuration instance
_config_instance: Optional[Config] = None


def load_config() -> Config:
    """
    Load configuration and return the global configuration instance.
    
    This function maintains backward compatibility with the original load_config() function
    but now returns a Config object instead of setting environment variables.
    """
    global _config_instance
    if _config_instance is None:
        _config_instance = Config()
    return _config_instance


def get_config() -> Config:
    """Get the global configuration instance."""
    global _config_instance
    if _config_instance is None:
        _config_instance = Config()
    return _config_instance


if __name__ == "__main__":
    # For backward compatibility, still output shell export commands when run directly
    config = load_config()
    
    # Output shell export commands for backward compatibility
    for key, value in config._config_data.items():
        # Escape quotes and special characters for shell
        escaped_value = str(value).replace('"', '\\"').replace("'", "\\'")
        print(f'export {key}="{escaped_value}"')
