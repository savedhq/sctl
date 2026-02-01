package job

import (
	"encoding/json"
	"fmt"

	"github.com/fatih/color"
	"github.com/savedhq/sctl/internal"
	"github.com/spf13/cobra"
)

func newJobUpdateCmd() *cobra.Command {
	var workspaceID string
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "update <job_id>",
		Short: "Update a job",
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

			_, r, err := cliCtx.Client.JobsAPI.UpdateJob(cliCtx.APICtx, workspaceID, jobID).Execute()
			if err != nil {
				return internal.PrintAPIError(err)
			}
			defer r.Body.Close()

			if jsonOutput {
				fmt.Println(`{"status": "success", "message": "Job updated"}`)
			} else {
				color.Green("✓ Job updated")
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output as JSON")
	return cmd
}

func newJobDeleteCmd() *cobra.Command {
	var workspaceID string
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "delete <job_id>",
		Short: "Delete a job",
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

			r, err := cliCtx.Client.JobsAPI.DeleteJob(cliCtx.APICtx, workspaceID, jobID).Execute()
			if err != nil {
				return internal.PrintAPIError(err)
			}
			defer r.Body.Close()

			if jsonOutput {
				fmt.Println(`{"status": "success", "message": "Job deleted"}`)
			} else {
				color.Green("✓ Job deleted")
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output as JSON")
	return cmd
}

func newJobTriggerCmd() *cobra.Command {
	var workspaceID string
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "trigger <job_id>",
		Short: "Trigger a job",
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

			resp, r, err := cliCtx.Client.JobsAPI.TriggerJob(cliCtx.APICtx, workspaceID, jobID).Execute()
			if err != nil {
				return internal.PrintAPIError(err)
			}
			defer r.Body.Close()

			if jsonOutput {
				data, err := json.MarshalIndent(resp, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal json output: %w", err)
				}
				fmt.Println(string(data))
			} else {
				color.Green("✓ Job triggered")
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output as JSON")
	return cmd
}
