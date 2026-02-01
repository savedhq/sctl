package job

import (
	"github.com/fatih/color"
	"github.com/savedhq/sctl/internal"
	saved "github.com/savedhq/sdk-go"
	"github.com/spf13/cobra"
)

func newJobCreateWorkerCmd() *cobra.Command {
	var workspaceID, name, schedule string
	cmd := &cobra.Command{
		Use:   "create-worker",
		Short: "Create a worker job",
		Run: func(cmd *cobra.Command, args []string) {
			cliCtx := internal.GetCLIContext(cmd.Context())
			internal.CheckErr(cliCtx.Err)

			var err error
			workspaceID, err = cliCtx.ResolveWorkspaceID(workspaceID)
			internal.CheckErr(err)

			req := saved.CreateWorkerJobRequest{}
			if name != "" {
				req.SetName(name)
			}
			if schedule != "" {
				req.SetSchedule(schedule)
			}

			resp, r, err := cliCtx.Client.JobsAPI.CreateWorkerJob(cliCtx.APICtx, workspaceID).CreateWorkerJobRequest(req).Execute()
			internal.CheckErr(internal.PrintAPIError(err))
			defer r.Body.Close()

			if cliCtx.JSONOutput {
				internal.PrintJSON(resp)
			} else {
				color.Green("✓ Worker job created")
				color.Cyan("ID: %s", resp.GetId())
			}
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
		Run: func(cmd *cobra.Command, args []string) {
			cliCtx := internal.GetCLIContext(cmd.Context())
			internal.CheckErr(cliCtx.Err)

			var err error
			workspaceID, err = cliCtx.ResolveWorkspaceID(workspaceID)
			internal.CheckErr(err)

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
			internal.CheckErr(internal.PrintAPIError(err))
			defer r.Body.Close()

			if cliCtx.JSONOutput {
				internal.PrintJSON(resp)
			} else {
				color.Green("✓ Agent job created")
				color.Cyan("ID: %s", resp.GetId())
			}
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
		Run: func(cmd *cobra.Command, args []string) {
			cliCtx := internal.GetCLIContext(cmd.Context())
			internal.CheckErr(cliCtx.Err)

			var err error
			workspaceID, err = cliCtx.ResolveWorkspaceID(workspaceID)
			internal.CheckErr(err)

			req := saved.CreateManualJobRequest{}
			if name != "" {
				req.SetName(name)
			}

			resp, r, err := cliCtx.Client.JobsAPI.CreateManualJob(cliCtx.APICtx, workspaceID).CreateManualJobRequest(req).Execute()
			internal.CheckErr(internal.PrintAPIError(err))
			defer r.Body.Close()

			if cliCtx.JSONOutput {
				internal.PrintJSON(resp)
			} else {
				color.Green("✓ Manual job created")
				color.Cyan("ID: %s", resp.GetId())
			}
		},
	}
	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	cmd.Flags().StringVarP(&name, "name", "n", "", "Job name")
	return cmd
}
