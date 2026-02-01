package backup

import (
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
	var workspaceID, jobID string
	cmd := &cobra.Command{
		Use:   "list <job_id>",
		Short: "List backups for a job",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cliCtx := internal.GetCLIContext(cmd.Context())
			internal.CheckErr(cliCtx.Err)

			var err error
			workspaceID, err = cliCtx.ResolveWorkspaceID(workspaceID)
			internal.CheckErr(err)

			if len(args) > 0 {
				jobID, err = cliCtx.ResolveJobID(workspaceID, args[0])
				internal.CheckErr(err)
			}
			if jobID == "" {
				internal.CheckErr(fmt.Errorf("job_id required"))
			}

			resp, r, err := cliCtx.Client.BackupsAPI.ListBackups(cliCtx.APICtx, workspaceID, jobID).Execute()
			internal.CheckErr(internal.PrintAPIError(err))
			defer r.Body.Close()

			if cliCtx.JSONOutput {
				internal.PrintJSON(resp)
			} else {
				if len(resp) == 0 {
					color.Yellow("⚠ No backups found")
					return
				}

				for _, backup := range resp {
					color.Cyan("ID: %s", backup.GetId())
					fmt.Printf("  Status: %s\n", backup.GetStatus())
					fmt.Printf("  Created: %s\n\n", backup.GetCreatedAt().Format("2006-01-02 15:04:05"))
				}
			}
		},
	}
	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	cmd.Flags().StringVarP(&jobID, "job", "j", "", "Job ID")
	return cmd
}

func newBackupGetCmd() *cobra.Command {
	var workspaceID, jobID string
	cmd := &cobra.Command{
		Use:   "get <job_id> <backup_id>",
		Short: "Get backup details",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			cliCtx := internal.GetCLIContext(cmd.Context())
			internal.CheckErr(cliCtx.Err)

			var err error
			workspaceID, err = cliCtx.ResolveWorkspaceID(workspaceID)
			internal.CheckErr(err)

			jobID, err = cliCtx.ResolveJobID(workspaceID, args[0])
			internal.CheckErr(err)
			backupID := args[1]

			resp, r, err := cliCtx.Client.BackupsAPI.GetBackup(cliCtx.APICtx, workspaceID, jobID, backupID).Execute()
			internal.CheckErr(internal.PrintAPIError(err))
			defer r.Body.Close()

			if cliCtx.JSONOutput {
				internal.PrintJSON(resp)
			} else {
				color.Cyan("ID: %s", resp.GetId())
				fmt.Printf("Status: %s\n", resp.GetStatus())
				fmt.Printf("Created: %s\n", resp.GetCreatedAt().Format("2006-01-02 15:04:05"))
			}
		},
	}
	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	return cmd
}

func newBackupDeleteCmd() *cobra.Command {
	var workspaceID string
	cmd := &cobra.Command{
		Use:   "delete <job_id> <backup_id>",
		Short: "Delete a backup",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			cliCtx := internal.GetCLIContext(cmd.Context())
			internal.CheckErr(cliCtx.Err)

			var err error
			workspaceID, err = cliCtx.ResolveWorkspaceID(workspaceID)
			internal.CheckErr(err)

			jobID, err := cliCtx.ResolveJobID(workspaceID, args[0])
			internal.CheckErr(err)
			backupID := args[1]

			r, err := cliCtx.Client.BackupsAPI.DeleteBackup(cliCtx.APICtx, workspaceID, jobID, backupID).Execute()
			internal.CheckErr(internal.PrintAPIError(err))
			defer r.Body.Close()

			if cliCtx.JSONOutput {
				internal.PrintJSON(map[string]string{"status": "ok"})
			} else {
				color.Green("✓ Backup deleted")
			}
		},
	}
	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	return cmd
}

func newBackupRequestCmd() *cobra.Command {
	var workspaceID, jobID string
	cmd := &cobra.Command{
		Use:   "request <job_id>",
		Short: "Request a new backup",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cliCtx := internal.GetCLIContext(cmd.Context())
			internal.CheckErr(cliCtx.Err)

			var err error
			workspaceID, err = cliCtx.ResolveWorkspaceID(workspaceID)
			internal.CheckErr(err)

			jobID, err = cliCtx.ResolveJobID(workspaceID, args[0])
			internal.CheckErr(err)

			resp, r, err := cliCtx.Client.JobOperationsAPI.RequestBackup(cliCtx.APICtx, workspaceID, jobID).Execute()
			internal.CheckErr(internal.PrintAPIError(err))
			defer r.Body.Close()

			if cliCtx.JSONOutput {
				internal.PrintJSON(resp)
			} else {
				color.Green("✓ Backup requested")
			}
		},
	}
	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	return cmd
}

func newBackupDownloadCmd() *cobra.Command {
	var workspaceID string
	cmd := &cobra.Command{
		Use:   "download <job_id> <backup_id>",
		Short: "Download a backup",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			cliCtx := internal.GetCLIContext(cmd.Context())
			internal.CheckErr(cliCtx.Err)

			var err error
			workspaceID, err = cliCtx.ResolveWorkspaceID(workspaceID)
			internal.CheckErr(err)

			jobID, err := cliCtx.ResolveJobID(workspaceID, args[0])
			internal.CheckErr(err)
			backupID := args[1]

			resp, r, err := cliCtx.Client.BackupsAPI.DownloadBackup(cliCtx.APICtx, workspaceID, jobID, backupID).Execute()
			internal.CheckErr(internal.PrintAPIError(err))
			defer r.Body.Close()

			if cliCtx.JSONOutput {
				internal.PrintJSON(resp)
			} else {
				color.Green("✓ Backup download URL generated")
			}
		},
	}
	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	return cmd
}
