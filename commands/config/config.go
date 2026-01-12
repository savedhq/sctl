package config

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/savedhq/sctl/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

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
			color.Green("✓ Configuration initialized at ~/.sctl/config.yaml")
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
			color.Green("✓ Set %s = %s", key, value)
			return nil
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "get <key>",
		Short: "Get a configuration value",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			value := viper.Get(args[0])
			if value == nil {
				color.Yellow("⚠ %s is not set", args[0])
				return nil
			}
			fmt.Println(value)
			return nil
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List all configuration values",
		RunE: func(cmd *cobra.Command, args []string) error {
			settings := viper.AllSettings()
			if len(settings) == 0 {
				color.Yellow("⚠ No configuration set")
				return nil
			}
			for key, value := range settings {
				fmt.Printf("%s = %v\n", key, value)
			}
			return nil
		},
	})

	return cmd
}
