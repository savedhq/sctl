package job

import (
	"encoding/json"
	"fmt"

	"github.com/fatih/color"
	"github.com/savedhq/sctl/internal"
	saved "github.com/savedhq/sdk-go"
	"github.com/spf13/cobra"
)

func newJobUpdateCmd() *cobra.Command {
	var workspaceID, name, schedule string
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

			req := saved.UpdateJobRequest{}
			if name != "" {
				req.SetName(name)
			}
			if schedule != "" {
				req.SetSchedule(schedule)
			}

			_, r, err := cliCtx.Client.JobsAPI.UpdateJob(cliCtx.APICtx, workspaceID, jobID).UpdateJobRequest(req).Execute()
			if err != nil {
				return internal.PrintAPIError(err)
			}
			defer r.Body.Close()

			color.Green("✓ Job updated")
			return nil
		},
	}
	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	cmd.Flags().StringVarP(&name, "name", "n", "", "Job name")
	cmd.Flags().StringVarP(&schedule, "schedule", "s", "", "Job schedule (cron format)")
	return cmd
}

func newJobDeleteCmd() *cobra.Command {
	var workspaceID string
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

			color.Green("✓ Job deleted")
			return nil
		},
	}
	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	return cmd
}

func newJobTriggerCmd() *cobra.Command {
	var workspaceID string
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

			color.Green("✓ Job triggered")
			data, _ := json.MarshalIndent(resp, "", "  ")
			fmt.Fprintln(cmd.OutOrStdout(), string(data))
			return nil
		},
	}
	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	return cmd
}
