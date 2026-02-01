package job

import (
	"encoding/json"
	"fmt"

	"github.com/fatih/color"
	"github.com/savedhq/sctl/internal"
	"github.com/spf13/cobra"
)

func newJobListCmd() *cobra.Command {
	var workspaceID string
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all jobs",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := internal.GetCLIContext(cmd.Context())
			if cliCtx == nil {
				return fmt.Errorf("CLI context not initialized")
			}

			var err error
			workspaceID, err = cliCtx.ResolveWorkspaceID(workspaceID)
			if err != nil {
				return err
			}

			resp, r, err := cliCtx.Client.JobsAPI.ListJobs(cliCtx.APICtx, workspaceID).Execute()
			if err != nil {
				return internal.PrintAPIError(err)
			}
			defer r.Body.Close()

			if len(resp) == 0 {
				color.Yellow("⚠ No jobs found")
				return nil
			}

			if jsonOutput {
				data, _ := json.MarshalIndent(resp, "", "  ")
				fmt.Fprintln(cmd.OutOrStdout(), string(data))
			} else {
				for _, job := range resp {
					color.New(color.FgCyan).Fprintf(cmd.OutOrStdout(), "ID: %s\n", job.GetId())
					fmt.Fprintf(cmd.OutOrStdout(), "  Name: %s\n", job.GetName())
					fmt.Fprintf(cmd.OutOrStdout(), "  Type: %s\n", job.GetType())
					fmt.Fprintf(cmd.OutOrStdout(), "  Enabled: %v\n\n", job.GetEnabled())
				}
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output as JSON")
	return cmd
}

func newJobGetCmd() *cobra.Command {
	var workspaceID string
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "get <job_id>",
		Short: "Get job details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := internal.GetCLIContext(cmd.Context())
			if cliCtx == nil {
				return fmt.Errorf("CLI context not initialized")
			}

			var err error
			workspaceID, err = cliCtx.ResolveWorkspaceID(workspaceID)
			if err != nil {
				return err
			}

			jobID, err := cliCtx.ResolveJobID(workspaceID, args[0])
			if err != nil {
				return err
			}

			resp, r, err := cliCtx.Client.JobsAPI.GetJob(cliCtx.APICtx, workspaceID, jobID).Execute()
			if err != nil {
				return internal.PrintAPIError(err)
			}
			defer r.Body.Close()

			if jsonOutput {
				data, _ := json.MarshalIndent(resp, "", "  ")
				fmt.Fprintln(cmd.OutOrStdout(), string(data))
			} else {
				color.New(color.FgCyan).Fprintf(cmd.OutOrStdout(), "ID: %s\n", resp.GetId())
				fmt.Fprintf(cmd.OutOrStdout(), "Name: %s\n", resp.GetName())
				fmt.Fprintf(cmd.OutOrStdout(), "Type: %s\n", resp.GetType())
				fmt.Fprintf(cmd.OutOrStdout(), "Enabled: %v\n", resp.GetEnabled())
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output as JSON")
	return cmd
}
