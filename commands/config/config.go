package config

import (
	"fmt"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/savedhq/sctl/internal"
	"github.com/savedhq/sctl/internal/render"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ConfigValue represents a single configuration value for rendering.
type ConfigValue struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

// String implements the Stringer interface for ConfigValue.
func (c ConfigValue) String() string {
	return fmt.Sprintf("%v", c.Value)
}

// ConfigSettings represents all configuration settings for rendering.
type ConfigSettings map[string]interface{}

// String implements the Stringer interface for ConfigSettings.
func (s ConfigSettings) String() string {
	var b strings.Builder
	keys := make([]string, 0, len(s))
	for k := range s {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		b.WriteString(fmt.Sprintf("%s = %v\n", k, s[k]))
	}
	return strings.TrimSuffix(b.String(), "\n")
}

func NewConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:         "config",
		Short:       "Manage CLI configuration",
		Annotations: map[string]string{"auth": "none"},
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "init",
		Short: "Initialize configuration file",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := internal.InitConfig(); err != nil {
				return err
			}
			render.Message(color.GreenString("✓ Configuration initialized at ~/.sctl/config.yaml"))
			return nil
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set a configuration value",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			key, value := args[0], args[1]
			viper.Set(key, value)
			if err := viper.WriteConfig(); err != nil {
				return err
			}
			render.Message(color.GreenString("✓ Set %s = %s", key, value))
			return nil
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "get <key>",
		Short: "Get a configuration value",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]
			value := viper.Get(key)
			if value == nil {
				render.Message(color.YellowString("⚠ %s is not set", key))
				return nil
			}
			render.Object(ConfigValue{Key: key, Value: value})
			return nil
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List all configuration values",
		RunE: func(cmd *cobra.Command, args []string) error {
			settings := viper.AllSettings()
			if len(settings) == 0 {
				render.Message(color.YellowString("⚠ No configuration set"))
				return nil
			}
			render.Object(ConfigSettings(settings))
			return nil
		},
	})

	return cmd
}
