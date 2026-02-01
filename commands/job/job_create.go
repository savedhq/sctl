package job

import (
	"fmt"

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

			req := saved.CreateWorkerJobRequest{}
			if name != "" {
				req.SetName(name)
			}
			if schedule != "" {
				req.SetSchedule(schedule)
			}

			var resp saved.CreateWorkerJob201Response
			r, err := cliCtx.Client.JobsAPI.CreateWorkerJob(cliCtx.APICtx, workspaceID).CreateWorkerJobRequest(req).ExecuteWithBody(&resp)
			if err != nil {
				return internal.PrintAPIError(err)
			}
			defer r.Body.Close()

			color.Green("✓ Worker job created")
			color.Cyan("ID: %s", resp.GetId())
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
			if cliCtx == nil {
				return fmt.Errorf("CLI context not initialized")
			}

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

			var resp saved.CreateAgentJob201Response
			r, err := cliCtx.Client.JobsAPI.CreateAgentJob(cliCtx.APICtx, workspaceID).CreateAgentJobRequest(req).ExecuteWithBody(&resp)
			if err != nil {
				return internal.PrintAPIError(err)
			}
			defer r.Body.Close()

			color.Green("✓ Agent job created")
			color.Cyan("ID: %s", resp.GetId())
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
			if cliCtx == nil {
				return fmt.Errorf("CLI context not initialized")
			}

			var err error
			workspaceID, err = cliCtx.ResolveWorkspaceID(workspaceID)
			if err != nil {
				return err
			}

			req := saved.CreateManualJobRequest{}
			if name != "" {
				req.SetName(name)
			}

			var resp saved.CreateManualJob201Response
			r, err := cliCtx.Client.JobsAPI.CreateManualJob(cliCtx.APICtx, workspaceID).CreateManualJobRequest(req).ExecuteWithBody(&resp)
			if err != nil {
				return internal.PrintAPIError(err)
			}
			defer r.Body.Close()

			color.Green("✓ Manual job created")
			color.Cyan("ID: %s", resp.GetId())
			return nil
		},
	}
	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	cmd.Flags().StringVarP(&name, "name", "n", "", "Job name")
	return cmd
}
