package job

import (
	"github.com/fatih/color"
	"github.com/savedhq/sctl/internal"
	"github.com/savedhq/sctl/internal/render"
	saved "github.com/savedhq/sdk-go"
	"github.com/spf13/cobra"
)

func newJobCreateWorkerCmd() *cobra.Command {
	var workspaceID, name, schedule string
	cmd := &cobra.Command{
		Use:   "create-worker",
		Short: "Create a worker job",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := internal.GetCLIContext(cmd.Context())
			var err error
			workspaceID, err = cliCtx.ResolveWorkspaceID(workspaceID)
			if err != nil {
				return err
			}

			req := saved.CreateWorkerJobRequest{}
			if name != "" {
				req.SetName(name)
			}
			if schedule != "" {
				req.SetSchedule(schedule)
			}

			resp, r, err := cliCtx.Client.JobsAPI.CreateWorkerJob(cliCtx.APICtx, workspaceID).CreateWorkerJobRequest(req).Execute()
			if err != nil {
				return internal.PrintAPIError(err)
			}
			defer r.Body.Close()

			render.Message(color.GreenString("✓ Worker job created: %s", resp.GetId()))
			return nil
		},
	}
	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	cmd.Flags().StringVarP(&name, "name", "n", "", "Job name")
	cmd.Flags().StringVarP(&schedule, "schedule", "s", "", "Job schedule (cron format)")
	return cmd
}

func newJobCreateAgentCmd() *cobra.Command {
	var workspaceID, name, schedule, agentID string
	cmd := &cobra.Command{
		Use:   "create-agent",
		Short: "Create an agent job",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := internal.GetCLIContext(cmd.Context())
			var err error
			workspaceID, err = cliCtx.ResolveWorkspaceID(workspaceID)
			if err != nil {
				return err
			}

			req := saved.CreateAgentJobRequest{}
			if name != "" {
				req.SetName(name)
			}
			if schedule != "" {
				req.SetSchedule(schedule)
			}
			if agentID != "" {
				req.SetAgentId(agentID)
			}

			resp, r, err := cliCtx.Client.JobsAPI.CreateAgentJob(cliCtx.APICtx, workspaceID).CreateAgentJobRequest(req).Execute()
			if err != nil {
				return internal.PrintAPIError(err)
			}
			defer r.Body.Close()

			render.Message(color.GreenString("✓ Agent job created: %s", resp.GetId()))
			return nil
		},
	}
	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	cmd.Flags().StringVarP(&name, "name", "n", "", "Job name")
	cmd.Flags().StringVarP(&schedule, "schedule", "s", "", "Job schedule (cron format)")
	cmd.Flags().StringVarP(&agentID, "agent", "a", "", "Agent ID")
	return cmd
}

func newJobCreateManualCmd() *cobra.Command {
	var workspaceID, name string
	cmd := &cobra.Command{
		Use:   "create-manual",
		Short: "Create a manual job",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := internal.GetCLIContext(cmd.Context())
			var err error
			workspaceID, err = cliCtx.ResolveWorkspaceID(workspaceID)
			if err != nil {
				return err
			}

			req := saved.CreateManualJobRequest{}
			if name != "" {
				req.SetName(name)
			}

			resp, r, err := cliCtx.Client.JobsAPI.CreateManualJob(cliCtx.APICtx, workspaceID).CreateManualJobRequest(req).Execute()
			if err != nil {
				return internal.PrintAPIError(err)
			}
			defer r.Body.Close()

			render.Message(color.GreenString("✓ Manual job created: %s", resp.GetId()))
			return nil
		},
	}
	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	cmd.Flags().StringVarP(&name, "name", "n", "", "Job name")
	return cmd
}
