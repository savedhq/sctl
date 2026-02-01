package job

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/savedhq/sctl/internal"
	"github.com/spf13/cobra"
)

func newJobListCmd() *cobra.Command {
	var workspaceID string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all jobs",
		Run: func(cmd *cobra.Command, args []string) {
			cliCtx := internal.GetCLIContext(cmd.Context())
			internal.CheckErr(cliCtx.Err)

			var err error
			workspaceID, err = cliCtx.ResolveWorkspaceID(workspaceID)
			internal.CheckErr(err)

			resp, r, err := cliCtx.Client.JobsAPI.ListJobs(cliCtx.APICtx, workspaceID).Execute()
			internal.CheckErr(internal.PrintAPIError(err))
			defer r.Body.Close()

			if cliCtx.JSONOutput {
				internal.PrintJSON(resp)
			} else {
				if len(resp) == 0 {
					color.Yellow("⚠ No jobs found")
					return
				}

				for _, job := range resp {
					color.Cyan("ID: %s", job.GetId())
					fmt.Printf("  Name: %s\n", job.GetName())
					fmt.Printf("  Type: %s\n", job.GetType())
					fmt.Printf("  Enabled: %v\n\n", job.GetEnabled())
				}
			}
		},
	}
	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	return cmd
}

func newJobGetCmd() *cobra.Command {
	var workspaceID string
	cmd := &cobra.Command{
		Use:   "get <job_id>",
		Short: "Get job details",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cliCtx := internal.GetCLIContext(cmd.Context())
			internal.CheckErr(cliCtx.Err)

			var err error
			workspaceID, err = cliCtx.ResolveWorkspaceID(workspaceID)
			internal.CheckErr(err)

			jobID, err := cliCtx.ResolveJobID(workspaceID, args[0])
			internal.CheckErr(err)

			resp, r, err := cliCtx.Client.JobsAPI.GetJob(cliCtx.APICtx, workspaceID, jobID).Execute()
			internal.CheckErr(internal.PrintAPIError(err))
			defer r.Body.Close()

			if cliCtx.JSONOutput {
				internal.PrintJSON(resp)
			} else {
				color.Cyan("ID: %s", resp.GetId())
				fmt.Printf("Name: %s\n", resp.GetName())
				fmt.Printf("Type: %s\n", resp.GetType())
				fmt.Printf("Enabled: %v\n", resp.GetEnabled())
			}
		},
	}
	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	return cmd
}
