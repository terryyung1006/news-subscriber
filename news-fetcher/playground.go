package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Fetchers []struct {
		Name   string `yaml:"name"`
		APIKey string `yaml:"api_key"`
		URL    string `yaml:"url"`
	} `yaml:"fetchers"`
}

func main() {
	// 1. Read config file
	data, err := os.ReadFile("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	// 2. Parse YAML
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}

	// 3. Find NewsAPI key
	var apiKey string
	for _, f := range cfg.Fetchers {
		if f.Name == "newsapi" {
			apiKey = f.APIKey
			break
		}
	}

	if apiKey == "" {
		log.Fatal("NewsAPI key not found in config")
	}

	// 4. Make request
	// Using a simple endpoint for testing
	url := fmt.Sprintf("https://newsapi.org/v2/top-headlines?country=us&apiKey=%s", apiKey)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Printf("Body: %s\n", string(body))
}
