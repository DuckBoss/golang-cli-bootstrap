package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConsolidate(t *testing.T) {
	tests := []struct {
		name           string
		configPath     string
		envPath        string
		defaultLogPath string
		cliLogPath     string
		setupToml      string // if non-empty, create temp dir and write this as appname.toml; configPath will be set to that path
		wantLogPath    string
		wantErr        bool
	}{
		{
			name:           "default log path",
			configPath:     "",
			envPath:        "",
			defaultLogPath: "/var/log/appname/appname.log",
			cliLogPath:     "",
			wantLogPath:    "/var/log/appname/appname.log",
			wantErr:        false,
		},
		{
			name:           "CLI overrides",
			configPath:     "",
			envPath:        "",
			defaultLogPath: "/default.log",
			cliLogPath:     "/cli.log",
			wantLogPath:    "/cli.log",
			wantErr:        false,
		},
		{
			name:           "missing config file treated as empty",
			configPath:     "/nonexistent/appname.toml",
			envPath:        "",
			defaultLogPath: "/var/log/appname/appname.log",
			cliLogPath:     "",
			wantLogPath:    "/var/log/appname/appname.log",
			wantErr:        false,
		},
		{
			name:           "config file then CLI",
			defaultLogPath: "/default.log",
			cliLogPath:     "/cli.log",
			setupToml:      "log_path = \"/from-file.log\"\n",
			wantLogPath:    "/cli.log",
			wantErr:        false,
		},
		{
			name:           "config file only",
			defaultLogPath: "/default.log",
			cliLogPath:     "",
			setupToml:      "log_path = \"/from-file.log\"\n",
			wantLogPath:    "/from-file.log",
			wantErr:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configPath := tt.configPath
			if tt.setupToml != "" {
				dir := t.TempDir()
				tomlPath := filepath.Join(dir, "appname.toml")
				if err := os.WriteFile(tomlPath, []byte(tt.setupToml), 0644); err != nil {
					t.Fatal(err)
				}
				configPath = tomlPath
			}
			cfg, err := Consolidate(configPath, tt.envPath, tt.defaultLogPath, tt.cliLogPath)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Consolidate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && cfg.LogPath != tt.wantLogPath {
				t.Errorf("LogPath = %q, want %q", cfg.LogPath, tt.wantLogPath)
			}
		})
	}
}

func TestApplyFromEnvMap(t *testing.T) {
	tests := []struct {
		name        string
		env         map[string]string
		wantLogPath string
	}{
		{
			name:        "APPNAME_LOG",
			env:         map[string]string{"APPNAME_LOG": "/appname.log"},
			wantLogPath: "/appname.log",
		},
		{
			name:        "LOG_PATH",
			env:         map[string]string{"LOG_PATH": "/env-log.log"},
			wantLogPath: "/env-log.log",
		},
		{
			name:        "prefer APPNAME_LOG over LOG_PATH",
			env:         map[string]string{"APPNAME_LOG": "/appname.log", "LOG_PATH": "/logpath.log"},
			wantLogPath: "/appname.log",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{}
			cfg.applyFromEnvMap(tt.env)
			if cfg.LogPath != tt.wantLogPath {
				t.Errorf("LogPath = %q, want %q", cfg.LogPath, tt.wantLogPath)
			}
		})
	}
}
