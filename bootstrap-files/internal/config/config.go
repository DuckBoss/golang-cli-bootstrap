package config

import (
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/joho/godotenv"
)

// Config holds consolidated application configuration.
// Order of precedence (low to high): config file (TOML) → environment file (.env) → process env → CLI args.
type Config struct {
	// LogPath is the path to the log file.
	LogPath string
}

// fileConfig is the shape of the TOML config file.
type fileConfig struct {
	LogPath string `toml:"log_path"`
}

// envPrefix is the prefix for environment variables (e.g. APPNAME_LOG).
const envPrefix = "APPNAME_"

// loadConfigFile reads a TOML config file and returns a Config with non-zero values.
// Missing file is treated as empty (no error).
func loadConfigFile(path string) (*Config, error) {
	var f fileConfig
	_, err := toml.DecodeFile(path, &f)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil
		}
		return nil, err
	}
	cfg := &Config{}
	if f.LogPath != "" {
		cfg.LogPath = f.LogPath
	}
	return cfg, nil
}

// loadEnvFile reads a .env file via godotenv and returns a map with uppercase keys.
// Missing file is treated as empty (no error).
func loadEnvFile(path string) (map[string]string, error) {
	m, err := godotenv.Read(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[strings.ToUpper(strings.TrimSpace(k))] = strings.TrimSpace(v)
	}
	return out, nil
}

// applyFromEnvMap overlays non-empty values from a .env-derived map onto cfg.
// Map keys should be uppercase; known keys: LOG_PATH, LOG, APPNAME_LOG.
func (cfg *Config) applyFromEnvMap(m map[string]string) {
	if m == nil {
		return
	}
	for _, key := range []string{"APPNAME_LOG", "LOG_PATH", "LOG"} {
		if v := m[key]; v != "" {
			cfg.LogPath = v
			break
		}
	}
}

// ApplyEnv overlays process environment variables with prefix envPrefix (e.g. APPNAME_LOG).
func (cfg *Config) ApplyEnv() {
	for _, e := range os.Environ() {
		if !strings.HasPrefix(e, envPrefix) {
			continue
		}
		idx := strings.Index(e, "=")
		if idx <= 0 {
			continue
		}
		key := strings.ToUpper(e[:idx])
		val := strings.TrimSpace(e[idx+1:])
		switch key {
		case "APPNAME_LOG":
			if val != "" {
				cfg.LogPath = val
			}
		}
	}
}

// Consolidate builds the final config from TOML config file, .env file, process env, and CLI.
// configPath must be a TOML file; envPath must be a .env file.
// CLI values override when non-empty.
func Consolidate(configPath, envPath, defaultLogPath, cliLogPath string) (*Config, error) {
	cfg := &Config{LogPath: defaultLogPath}

	// 1) Config file (TOML, lowest priority)
	if configPath != "" {
		fileCfg, err := loadConfigFile(configPath)
		if err != nil {
			return nil, err
		}
		if fileCfg.LogPath != "" {
			cfg.LogPath = fileCfg.LogPath
		}
	}

	// 2) Environment file (.env)
	if envPath != "" {
		m, err := loadEnvFile(envPath)
		if err != nil {
			return nil, err
		}
		cfg.applyFromEnvMap(m)
	}

	// 3) Process environment
	cfg.ApplyEnv()

	// 4) CLI args (highest priority)
	if cliLogPath != "" {
		cfg.LogPath = cliLogPath
	}

	if cfg.LogPath == "" {
		cfg.LogPath = defaultLogPath
	}

	return cfg, nil
}
