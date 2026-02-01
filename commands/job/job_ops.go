package job

import (
	"github.com/fatih/color"
	"github.com/savedhq/sctl/internal"
	"github.com/spf13/cobra"
)

func newJobUpdateCmd() *cobra.Command {
	var workspaceID string
	cmd := &cobra.Command{
		Use:   "update <job_id>",
		Short: "Update a job",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cliCtx := internal.GetCLIContext(cmd.Context())
			internal.CheckErr(cliCtx.Err)

			var err error
			workspaceID, err = cliCtx.ResolveWorkspaceID(workspaceID)
			internal.CheckErr(err)

			jobID, err := cliCtx.ResolveJobID(workspaceID, args[0])
			internal.CheckErr(err)

			_, r, err := cliCtx.Client.JobsAPI.UpdateJob(cliCtx.APICtx, workspaceID, jobID).Execute()
			internal.CheckErr(internal.PrintAPIError(err))
			defer r.Body.Close()

			if cliCtx.JSONOutput {
				internal.PrintJSON(map[string]string{"status": "ok"})
			} else {
				color.Green("✓ Job updated")
			}
		},
	}
	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	return cmd
}

func newJobDeleteCmd() *cobra.Command {
	var workspaceID string
	cmd := &cobra.Command{
		Use:   "delete <job_id>",
		Short: "Delete a job",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cliCtx := internal.GetCLIContext(cmd.Context())
			internal.CheckErr(cliCtx.Err)

			var err error
			workspaceID, err = cliCtx.ResolveWorkspaceID(workspaceID)
			internal.CheckErr(err)

			jobID, err := cliCtx.ResolveJobID(workspaceID, args[0])
			internal.CheckErr(err)

			r, err := cliCtx.Client.JobsAPI.DeleteJob(cliCtx.APICtx, workspaceID, jobID).Execute()
			internal.CheckErr(internal.PrintAPIError(err))
			defer r.Body.Close()

			if cliCtx.JSONOutput {
				internal.PrintJSON(map[string]string{"status": "ok"})
			} else {
				color.Green("✓ Job deleted")
			}
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
		Run: func(cmd *cobra.Command, args []string) {
			cliCtx := internal.GetCLIContext(cmd.Context())
			internal.CheckErr(cliCtx.Err)

			var err error
			workspaceID, err = cliCtx.ResolveWorkspaceID(workspaceID)
			internal.CheckErr(err)

			jobID, err := cliCtx.ResolveJobID(workspaceID, args[0])
			internal.CheckErr(err)

			resp, r, err := cliCtx.Client.JobsAPI.TriggerJob(cliCtx.APICtx, workspaceID, jobID).Execute()
			internal.CheckErr(internal.PrintAPIError(err))
			defer r.Body.Close()

			if cliCtx.JSONOutput {
				internal.PrintJSON(resp)
			} else {
				color.Green("✓ Job triggered")
			}
		},
	}
	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	return cmd
}
