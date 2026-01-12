package internal

import (
	"context"
	"fmt"
	"os"

	saved "github.com/savedhq/sdk-go"
	"github.com/spf13/viper"
)

type Config struct {
	APIKey      string
	WorkspaceID string
	ServerURL   string
	Debug       bool
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	homeDir, err := os.UserHomeDir()
	if err == nil {
		viper.AddConfigPath(homeDir + "/.sctl")
	}
	viper.AddConfigPath(".")

	viper.SetDefault("server_url", "https://api.saved.sh")
	viper.SetDefault("debug", false)

	viper.SetEnvPrefix("SAVED")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	return &Config{
		APIKey:      viper.GetString("api_key"),
		WorkspaceID: viper.GetString("workspace_id"),
		ServerURL:   viper.GetString("server_url"),
		Debug:       viper.GetBool("debug"),
	}, nil
}

func (c *Config) GetClient() (*saved.APIClient, context.Context) {
	cfg := saved.NewConfiguration()
	cfg.Debug = c.Debug

	// Use custom server URL if provided, otherwise default to production
	if c.ServerURL != "" {
		cfg.Servers = saved.ServerConfigurations{{URL: c.ServerURL}}
	}

	client := saved.NewAPIClient(cfg)
	ctx := context.WithValue(context.Background(), saved.ContextAccessToken, c.APIKey)

	return client, ctx
}

func (c *Config) Validate() error {
	if c.APIKey == "" {
		return fmt.Errorf("API key not configured. Set SAVED_API_KEY or use 'sctl config set api_key <key>'")
	}
	return nil
}

func InitConfig() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configDir := home + "/.sctl"
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	configFile := configDir + "/config.yaml"
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		f, err := os.Create(configFile)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = f.WriteString("# Saved CLI Configuration\n")
		_, err = f.WriteString("# api_key: your-api-key-here\n")
		_, err = f.WriteString("# workspace_id: your-default-workspace-id\n")
		_, err = f.WriteString("# server_url: https://api.saved.sh\n")
		return err
	}

	return nil
}
