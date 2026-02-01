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

			var resp []saved.ListJobs200ResponseInner
			r, err := cliCtx.Client.JobsAPI.ListJobs(cliCtx.APICtx, workspaceID).ExecuteWithBody(&resp)
			if err != nil {
				return internal.PrintAPIError(err)
			}
			defer r.Body.Close()

			if len(resp) == 0 {
				color.Yellow("⚠ No jobs found")
				return nil
			}

			for _, job := range resp {
				color.Cyan("ID: %s", job.GetId())
				fmt.Printf("  Name: %s\n", job.GetName())
				fmt.Printf("  Type: %s\n", job.GetType())
				fmt.Printf("  Enabled: %v\n\n", job.GetEnabled())
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
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

			var resp saved.ListJobs200ResponseInner
			r, err := cliCtx.Client.JobsAPI.GetJob(cliCtx.APICtx, workspaceID, jobID).ExecuteWithBody(&resp)
			if err != nil {
				return internal.PrintAPIError(err)
			}
			defer r.Body.Close()

			if jsonOutput {
				data, _ := json.MarshalIndent(resp, "", "  ")
				fmt.Println(string(data))
			} else {
				color.Cyan("ID: %s", resp.GetId())
				fmt.Printf("Name: %s\n", resp.GetName())
				fmt.Printf("Type: %s\n", resp.GetType())
				fmt.Printf("Enabled: %v\n", resp.GetEnabled())
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output as JSON")
	return cmd
}
