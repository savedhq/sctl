package backup

import (
	"encoding/json"
	"fmt"

	"github.com/fatih/color"
	"github.com/savedhq/sctl/internal"
	"github.com/spf13/cobra"
)

func NewBackupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "backup",
		Aliases: []string{"backups"},
		Short:   "Manage backups",
	}

	cmd.AddCommand(newBackupListCmd())
	cmd.AddCommand(newBackupGetCmd())
	cmd.AddCommand(newBackupDeleteCmd())
	cmd.AddCommand(newBackupRequestCmd())
	cmd.AddCommand(newBackupDownloadCmd())

	return cmd
}

func newBackupListCmd() *cobra.Command {
	var workspaceID string
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "list <job_id>",
		Short: "List backups for a job",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := internal.GetCLIContext(cmd.Context())
			if cliCtx == nil {
				return fmt.Errorf("CLI context not initialized")
			}

			if workspaceID == "" {
				workspaceID = cliCtx.GetWorkspaceID()
			}
			if workspaceID == "" {
				return fmt.Errorf("workspace_id required (use --workspace or set in config)")
			}

			jobID := args[0]

			resp, r, err := cliCtx.Client.BackupsAPI.ListBackups(cliCtx.APICtx, workspaceID, jobID).Execute()
			if err != nil {
				return fmt.Errorf("API error: %v", err)
			}
			defer r.Body.Close()

			if jsonOutput {
				data, _ := json.MarshalIndent(resp, "", "  ")
				fmt.Println(string(data))
				return nil
			}

			if len(resp) == 0 {
				color.Yellow("⚠ No backups found")
				return nil
			}

			for _, backup := range resp {
				color.Cyan("ID: %s", backup.GetId())
				fmt.Printf("  Status: %s\n", backup.GetStatus())
				fmt.Printf("  Created: %s\n\n", backup.GetCreatedAt().Format("2006-01-02 15:04:05"))
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output as JSON")
	return cmd
}

func newBackupGetCmd() *cobra.Command {
	var workspaceID string
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "get <job_id> <backup_id>",
		Short: "Get backup details",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := internal.GetCLIContext(cmd.Context())
			if cliCtx == nil {
				return fmt.Errorf("CLI context not initialized")
			}

			if workspaceID == "" {
				workspaceID = cliCtx.GetWorkspaceID()
			}
			if workspaceID == "" {
				return fmt.Errorf("workspace_id required (use --workspace or set in config)")
			}

			jobID := args[0]
			backupID := args[1]

			resp, r, err := cliCtx.Client.BackupsAPI.GetBackup(cliCtx.APICtx, workspaceID, jobID, backupID).Execute()
			if err != nil {
				return fmt.Errorf("API error: %v", err)
			}
			defer r.Body.Close()

			if jsonOutput {
				data, _ := json.MarshalIndent(resp, "", "  ")
				fmt.Println(string(data))
			} else {
				color.Cyan("ID: %s", resp.GetId())
				fmt.Printf("Status: %s\n", resp.GetStatus())
				fmt.Printf("Created: %s\n", resp.GetCreatedAt().Format("2006-01-02 15:04:05"))
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output as JSON")
	return cmd
}

func newBackupDeleteCmd() *cobra.Command {
	var workspaceID string
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "delete <job_id> <backup_id>",
		Short: "Delete a backup",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := internal.GetCLIContext(cmd.Context())
			if cliCtx == nil {
				return fmt.Errorf("CLI context not initialized")
			}

			if workspaceID == "" {
				workspaceID = cliCtx.GetWorkspaceID()
			}
			if workspaceID == "" {
				return fmt.Errorf("workspace_id required (use --workspace or set in config)")
			}

			jobID := args[0]
			backupID := args[1]

			r, err := cliCtx.Client.BackupsAPI.DeleteBackup(cliCtx.APICtx, workspaceID, jobID, backupID).Execute()
			if err != nil {
				return fmt.Errorf("API error: %v", err)
			}
			defer r.Body.Close()

			if jsonOutput {
				data := map[string]string{"status": "deleted", "id": backupID}
				jsonData, _ := json.MarshalIndent(data, "", "  ")
				fmt.Println(string(jsonData))
			} else {
				color.Green("✓ Backup deleted")
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output as JSON")
	return cmd
}

func newBackupRequestCmd() *cobra.Command {
	var workspaceID string
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "request <job_id>",
		Short: "Request a new backup",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := internal.GetCLIContext(cmd.Context())
			if cliCtx == nil {
				return fmt.Errorf("CLI context not initialized")
			}

			if workspaceID == "" {
				workspaceID = cliCtx.GetWorkspaceID()
			}
			if workspaceID == "" {
				return fmt.Errorf("workspace_id required (use --workspace or set in config)")
			}

			jobID := args[0]

			resp, r, err := cliCtx.Client.JobOperationsAPI.RequestBackup(cliCtx.APICtx, workspaceID, jobID).Execute()
			if err != nil {
				return fmt.Errorf("API error: %v", err)
			}
			defer r.Body.Close()

			if jsonOutput {
				data, _ := json.MarshalIndent(resp, "", "  ")
				fmt.Println(string(data))
			} else {
				color.Green("✓ Backup requested")
				color.Cyan("ID: %s", resp.GetId())
				fmt.Printf("Status: %s\n", resp.GetStatus())
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output as JSON")
	return cmd
}

func newBackupDownloadCmd() *cobra.Command {
	var workspaceID string
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "download <job_id> <backup_id>",
		Short: "Download a backup",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := internal.GetCLIContext(cmd.Context())
			if cliCtx == nil {
				return fmt.Errorf("CLI context not initialized")
			}

			if workspaceID == "" {
				workspaceID = cliCtx.GetWorkspaceID()
			}
			if workspaceID == "" {
				return fmt.Errorf("workspace_id required (use --workspace or set in config)")
			}

			jobID := args[0]
			backupID := args[1]

			resp, r, err := cliCtx.Client.BackupsAPI.DownloadBackup(cliCtx.APICtx, workspaceID, jobID, backupID).Execute()
			if err != nil {
				return fmt.Errorf("API error: %v", err)
			}
			defer r.Body.Close()

			if jsonOutput {
				data, _ := json.MarshalIndent(resp, "", "  ")
				fmt.Println(string(data))
			} else {
				color.Green("✓ Backup download URL generated")
				fmt.Println(resp.GetUrl())
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output as JSON")
	return cmd
}
