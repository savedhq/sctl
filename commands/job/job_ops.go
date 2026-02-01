package job

import (
	"encoding/json"

	"github.com/fatih/color"
	"github.com/savedhq/sctl/internal"
	"github.com/savedhq/sctl/internal/render"
	saved "github.com/savedhq/sdk-go"
	"github.com/spf13/cobra"
)

// JobTrigger is a wrapper for the saved.TriggerJob202Response type for rendering.
type JobTrigger saved.TriggerJob202Response

// String implements the Stringer interface for JobTrigger.
func (t JobTrigger) String() string {
	data, _ := json.MarshalIndent(t, "", "  ")
	return string(data)
}

func newJobUpdateCmd() *cobra.Command {
	var workspaceID, name, schedule, agentID string
	var enabled, disabled bool
	cmd := &cobra.Command{
		Use:   "update <job_id>",
		Short: "Update a job",
		Args:  cobra.ExactArgs(1),
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

			req := saved.UpdateJobRequest{}
			if name != "" {
				req.SetName(name)
			}
			if schedule != "" {
				req.SetSchedule(schedule)
			}
			if agentID != "" {
				// req.SetAgentId(agentID)
			}
			if enabled {
				req.SetEnabled(true)
			}
			if disabled {
				req.SetEnabled(false)
			}

			_, r, err := cliCtx.Client.JobsAPI.UpdateJob(cliCtx.APICtx, workspaceID, jobID).UpdateJobRequest(req).Execute()
			if err != nil {
				return internal.PrintAPIError(err)
			}
			defer r.Body.Close()

			render.Message(color.GreenString("✓ Job updated"))
			return nil
		},
	}
	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	cmd.Flags().StringVarP(&name, "name", "n", "", "Job name")
	cmd.Flags().StringVarP(&schedule, "schedule", "s", "", "Job schedule (cron format)")
	cmd.Flags().StringVarP(&agentID, "agent", "a", "", "Agent ID")
	cmd.Flags().BoolVar(&enabled, "enable", false, "Enable the job")
	cmd.Flags().BoolVar(&disabled, "disable", false, "Disable the job")
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

			render.Message(color.GreenString("✓ Job deleted"))
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

			render.Message(color.GreenString("✓ Job triggered"))
			render.Object(JobTrigger(*resp))
			return nil
		},
	}
	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	return cmd
}
