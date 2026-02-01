package backup

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/savedhq/sctl/internal"
	"github.com/savedhq/sctl/internal/render"
	saved "github.com/savedhq/sdk-go"
	"github.com/spf13/cobra"
)

// Backup is a wrapper for the saved.RequestBackup201Response type for rendering.
type Backup saved.RequestBackup201Response

// String implements the Stringer interface for Backup.
func (b Backup) String() string {
	return fmt.Sprintf("%s %s\n  Status: %s\n  Created: %s",
		color.CyanString("ID:"),
		b.Id,
		b.Status,
		b.CreatedAt.Format(time.RFC3339),
	)
}

// Backups is a wrapper for a slice of saved.RequestBackup201Response for rendering.
type Backups []saved.RequestBackup201Response

// String implements the Stringer interface for Backups.
func (b Backups) String() string {
	var sb strings.Builder
	for _, backup := range b {
		sb.WriteString(fmt.Sprintf("%s %s\n  Status: %s\n  Created: %s\n\n",
			color.CyanString("ID:"),
			backup.Id,
			backup.Status,
			backup.CreatedAt.Format(time.RFC3339),
		))
	}
	return strings.TrimSuffix(sb.String(), "\n\n")
}

// BackupRequest is a wrapper for the saved.RequestBackup201Response type for rendering.
type BackupRequest saved.RequestBackup201Response

// String implements the Stringer interface for BackupRequest.
func (b BackupRequest) String() string {
	data, _ := json.MarshalIndent(b, "", "  ")
	return string(data)
}

// BackupDownload is a wrapper for the saved.DownloadBackup200Response type for rendering.
type BackupDownload saved.DownloadBackup200Response

// String implements the Stringer interface for BackupDownload.
func (b BackupDownload) String() string {
	data, _ := json.MarshalIndent(b, "", "  ")
	return string(data)
}

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
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := internal.GetCLIContext(cmd.Context())
			var err error
			workspaceID, err = cliCtx.ResolveWorkspaceID(workspaceID)
			if err != nil {
				return err
			}

			if len(args) > 0 {
				jobID = args[0]
			}
			jobID, err = cliCtx.ResolveJobID(workspaceID, jobID)
			if err != nil {
				return err
			}

			resp, r, err := cliCtx.Client.BackupsAPI.ListBackups(cliCtx.APICtx, workspaceID, jobID).Execute()
			if err != nil {
				return internal.PrintAPIError(err)
			}
			defer r.Body.Close()

			if len(resp) == 0 {
				render.Message(color.YellowString("⚠ No backups found"))
				return nil
			}

			render.Object(Backups(resp))
			return nil
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
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := internal.GetCLIContext(cmd.Context())
			var err error
			workspaceID, err = cliCtx.ResolveWorkspaceID(workspaceID)
			if err != nil {
				return err
			}

			jobID, err = cliCtx.ResolveJobID(workspaceID, args[0])
			if err != nil {
				return err
			}

			backupID := args[1]

			resp, r, err := cliCtx.Client.BackupsAPI.GetBackup(cliCtx.APICtx, workspaceID, jobID, backupID).Execute()
			if err != nil {
				return internal.PrintAPIError(err)
			}
			defer r.Body.Close()

			render.Object(Backup(*resp))
			return nil
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
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := internal.GetCLIContext(cmd.Context())
			var err error
			workspaceID, err = cliCtx.ResolveWorkspaceID(workspaceID)
			if err != nil {
				return err
			}

			jobID, err := cliCtx.ResolveJobID(workspaceID, args[0])
			if err != nil {
				return err
			}
			backupID := args[1]

			r, err := cliCtx.Client.BackupsAPI.DeleteBackup(cliCtx.APICtx, workspaceID, jobID, backupID).Execute()
			if err != nil {
				return internal.PrintAPIError(err)
			}
			defer r.Body.Close()

			render.Message(color.GreenString("✓ Backup deleted"))
			return nil
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
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := internal.GetCLIContext(cmd.Context())
			var err error
			workspaceID, err = cliCtx.ResolveWorkspaceID(workspaceID)
			if err != nil {
				return err
			}

			jobID, err = cliCtx.ResolveJobID(workspaceID, args[0])
			if err != nil {
				return err
			}

			resp, r, err := cliCtx.Client.JobOperationsAPI.RequestBackup(cliCtx.APICtx, workspaceID, jobID).Execute()
			if err != nil {
				return internal.PrintAPIError(err)
			}
			defer r.Body.Close()

			render.Message(color.GreenString("✓ Backup requested"))
			render.Object(BackupRequest(*resp))
			return nil
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
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := internal.GetCLIContext(cmd.Context())
			var err error
			workspaceID, err = cliCtx.ResolveWorkspaceID(workspaceID)
			if err != nil {
				return err
			}

			jobID, err := cliCtx.ResolveJobID(workspaceID, args[0])
			if err != nil {
				return err
			}
			backupID := args[1]

			resp, r, err := cliCtx.Client.BackupsAPI.DownloadBackup(cliCtx.APICtx, workspaceID, jobID, backupID).Execute()
			if err != nil {
				return internal.PrintAPIError(err)
			}
			defer r.Body.Close()

			render.Message(color.GreenString("✓ Backup download URL generated"))
			render.Object(BackupDownload(*resp))
			return nil
		},
	}
	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	return cmd
}
