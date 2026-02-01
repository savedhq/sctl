package main

import (
	"fmt"
	"os"

	"github.com/savedhq/sctl/commands/agent"
	"github.com/savedhq/sctl/commands/auth"
	"github.com/savedhq/sctl/commands/backup"
	"github.com/savedhq/sctl/commands/billing"
	"github.com/savedhq/sctl/commands/config"
	"github.com/savedhq/sctl/commands/job"
	"github.com/savedhq/sctl/commands/schedule"
	"github.com/savedhq/sctl/commands/workspace"
	"github.com/savedhq/sctl/internal"
	"github.com/spf13/cobra"
)

const version = "1.0.0"

func main() {
	rootCmd := &cobra.Command{
		Use:   "sctl",
		Short: "Saved CLI - Encrypted backups and distributed storage",
		Long:  `A command-line interface for managing Saved workspaces, agents, jobs, backups, and billing.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Name() == "help" || cmd.Name() == "completion" || cmd.Name() == "version" {
				return nil
			}

			// Check if command explicitly skips auth
			if cmd.Annotations["auth"] == "none" {
				return nil
			}
			// Check if parent command skips auth (e.g. auth login -> auth)
			if cmd.Parent() != nil && cmd.Parent().Annotations["auth"] == "none" {
				return nil
			}

			cfg, err := internal.LoadConfig()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			cliCtx, err := internal.NewCLIContext(cfg)
			if err != nil {
				return err
			}

			cmd.SetContext(internal.WithCLIContext(cmd.Context(), cliCtx))
			return nil
		},
	}

	rootCmd.Version = version

	rootCmd.AddCommand(auth.NewAuthCmd())
	rootCmd.AddCommand(config.NewConfigCmd())
	rootCmd.AddCommand(workspace.NewWorkspaceCmd())
	rootCmd.AddCommand(agent.NewAgentCmd())
	rootCmd.AddCommand(job.NewJobCmd())
	rootCmd.AddCommand(schedule.NewScheduleCmd())
	rootCmd.AddCommand(backup.NewBackupCmd())
	rootCmd.AddCommand(billing.NewBillingCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
