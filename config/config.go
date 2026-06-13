package config

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

const (
	// EnvProduction production environment identifier.
	EnvProduction = "production"
	// EnvStaging staging environment identifier.
	EnvStaging = "staging"
	// EnvTesting testing environment identifier.
	EnvTesting = "testing"
)

// Config holds all application configuration values.
type Config struct {
	DBCONNECTION            string `yaml:"db_connection"`
	DBHOST                  string `yaml:"db_host"`
	DBPORT                  string `yaml:"db_port"`
	DBDATABASE              string `yaml:"db_database"`
	DBUSERNAME              string `yaml:"db_username"`
	DBPASSWORD              string `yaml:"db_password"`
	APIPORT                 string `yaml:"api_port"`
	NAME                    string `yaml:"name"`
	ENV                     string `yaml:"env"`
	LOGLEVEL                string `yaml:"logLevel"`
	MAINTENANCEMODE         bool   `yaml:"maintenanceMode"`
	GRACEFULSHUTDOWNTIMEOUT int    `yaml:"gracefulShutdownTimeout"`
	READHEADERTIMEOUT       int    `yaml:"readHeaderTimeout"`
	SSLMODE                 string `yamle:"sslMode"`
}

// Load loads configuration from defaults, .env, YAML file, and environment variables.
func Load() (*Config, error) {
	loadDotEnv(".env")

	cfg := defaults()

	path := os.Getenv("CONFIG_FILE")
	if path == "" {
		path = "./config.yaml"
	}

	if err := mergeYAML(path, cfg); err != nil {
		if os.Getenv("CONFIG_FILE") != "" {
			return &Config{}, err
		}
	}

	cfg.APIPORT = env("API_PORT", cfg.APIPORT)
	cfg.DBCONNECTION = env("DB_CONNECTION", cfg.DBCONNECTION)
	cfg.DBHOST = env("DB_HOST", cfg.DBHOST)
	cfg.DBPORT = env("DB_PORT", cfg.DBPORT)
	cfg.DBDATABASE = env("DB_DATABASE", cfg.DBDATABASE)
	cfg.DBUSERNAME = env("DB_USERNAME", cfg.DBUSERNAME)
	cfg.DBPASSWORD = env("DB_PASSWORD", cfg.DBPASSWORD)
	cfg.SSLMODE = env("SSL_MODE", cfg.SSLMODE)
	return cfg, nil
}

func defaults() *Config {
	return &Config{
		DBCONNECTION:            "mysql",
		DBHOST:                  "jobs-db",
		DBPORT:                  "3306",
		DBDATABASE:              "jobs",
		DBUSERNAME:              "jobs",
		DBPASSWORD:              "jobs",
		NAME:                    "Search Engine",
		ENV:                     EnvStaging,
		LOGLEVEL:                "info",
		GRACEFULSHUTDOWNTIMEOUT: 15,
		READHEADERTIMEOUT:       60,
		SSLMODE:                 "disable",
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

func mergeYAML(path string, cfg *Config) error {
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

// DSN constructs the Data Source Name for a MySQL connection.
func DSN(cfg *Config) string {
	return cfg.DBUSERNAME + ":" + cfg.DBPASSWORD + "@tcp(" + cfg.DBHOST + ":" + cfg.DBPORT + ")/" + cfg.DBDATABASE + "?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci"
}
