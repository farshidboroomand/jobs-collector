package configs

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// Configs holds all application configuration values.
type Configs struct {
	DATABASEHOST     string `yaml:"database_host"`
	DATABASEPORT     string `yaml:"database_port"`
	DATABASEUSER     string `yaml:"database_user"`
	DATABASEPASSWORD string `yaml:"database_password"`
	DATABASENAME     string `yaml:"database_name"`
	APIPORT          string `yaml:"api_port"`
}

// Load loads configuration from defaults, .env, YAML file, and environment variables.
func Load() (*Configs, error) {
	loadDotEnv(".env")

	cfg := defaults()

	path := os.Getenv("CONFIG_FILE")
	if path == "" {
		path = "./configs.yaml"
	}

	if err := mergeYAML(path, cfg); err != nil {
		if os.Getenv("CONFIG_FILE") != "" {
			return &Configs{}, err
		}
	}

	cfg.DATABASEHOST = env("DATABASE_HOST", cfg.DATABASEHOST)
	cfg.DATABASEPORT = env("DATABASE_PORT", cfg.DATABASEPORT)
	cfg.DATABASEUSER = env("DATABASE_USER", cfg.DATABASEUSER)
	cfg.DATABASEPASSWORD = env("DATABASE_PASSWORD", cfg.DATABASEPASSWORD)
	cfg.DATABASENAME = env("DATABASE_NAME", cfg.DATABASENAME)
	cfg.APIPORT = env("API_PORT", cfg.APIPORT)

	return cfg, nil
}

func defaults() *Configs {
	return &Configs{
		DATABASEHOST:     "localhost",
		DATABASEPORT:     "3306",
		DATABASEUSER:     "jobs",
		DATABASEPASSWORD: "jobs",
		DATABASENAME:     "jobs",
	}
}

func loadDotEnv(path string) {
	cleanPath := filepath.Clean(path)

	file, err := os.Open(cleanPath)
	if err != nil {
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.WithError(err).Warn("fail to close the env file.")
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || !strings.Contains(line, "=") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		key := strings.TrimSpace(parts[0])
		value := strings.Trim(strings.TrimSpace(parts[1]), "\"'")

		if os.Getenv(key) == "" {
			if err := os.Setenv(key, value); err != nil {
				continue
			}
		}
	}
}

func mergeYAML(path string, cfg *Configs) error {
	cleanPath := filepath.Clean(path)

	if strings.Contains(cleanPath, "..") {
		return errors.New("invalid config path")
	}

	bytes, err := os.ReadFile(cleanPath)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(bytes, cfg)
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
