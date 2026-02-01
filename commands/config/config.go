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
		Run: func(cmd *cobra.Command, args []string) {
			internal.CheckErr(internal.InitConfig())
			color.Green("✓ Configuration initialized at ~/.sctl/config.yaml")
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set a configuration value",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			key, value := args[0], args[1]
			viper.Set(key, value)
			internal.CheckErr(viper.WriteConfig())
			color.Green("✓ Set %s = %s", key, value)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "get <key>",
		Short: "Get a configuration value",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			value := viper.Get(args[0])
			if value == nil {
				color.Yellow("⚠ %s is not set", args[0])
				return
			}
			fmt.Println(value)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List all configuration values",
		Run: func(cmd *cobra.Command, args []string) {
			cliCtx := internal.GetCLIContext(cmd.Context())
			internal.CheckErr(cliCtx.Err)

			settings := viper.AllSettings()
			if cliCtx.JSONOutput {
				internal.PrintJSON(settings)
			} else {
				if len(settings) == 0 {
					color.Yellow("⚠ No configuration set")
					return
				}
				for key, value := range settings {
					fmt.Printf("%s = %v\n", key, value)
				}
			}
		},
	})

	return cmd
}
