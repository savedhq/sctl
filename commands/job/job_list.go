package job

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/savedhq/sctl/internal"
	"github.com/savedhq/sctl/internal/render"
	saved "github.com/savedhq/sdk-go"
	"github.com/spf13/cobra"
)

// Job is a wrapper for the saved.ListJobs200ResponseInner type for rendering.
type Job saved.ListJobs200ResponseInner

// String implements the Stringer interface for Job.
func (j Job) String() string {
	return fmt.Sprintf("%s %s (%s)",
		color.CyanString("ID:"),
		j.Id,
		j.Name,
	)
}

// Jobs is a wrapper for a slice of saved.ListJobs200ResponseInner for rendering.
type Jobs []saved.ListJobs200ResponseInner

// String implements the Stringer interface for Jobs.
func (j Jobs) String() string {
	var b strings.Builder
	for _, job := range j {
		b.WriteString(fmt.Sprintf("%s %s\n  Name: %s\n  Type: %s\n  Enabled: %v\n\n",
			color.CyanString("ID:"),
			job.Id,
			job.Name,
			job.Type,
			job.Enabled,
		))
	}
	return strings.TrimSuffix(b.String(), "\n\n")
}

func newJobListCmd() *cobra.Command {
	var workspaceID string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all jobs",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := internal.GetCLIContext(cmd.Context())
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
				render.Message(color.YellowString("⚠ No jobs found"))
				return nil
			}

			render.Object(Jobs(resp))
			return nil
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

			resp, r, err := cliCtx.Client.JobsAPI.GetJob(cliCtx.APICtx, workspaceID, jobID).Execute()
			if err != nil {
				return internal.PrintAPIError(err)
			}
			defer r.Body.Close()

			render.Object(Job(*resp))
			return nil
		},
	}
	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	return cmd
}
