package schedule

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/savedhq/sctl/internal"
	saved "github.com/savedhq/sdk-go"
	"github.com/spf13/cobra"
)

func newScheduleCreateCmd() *cobra.Command {
	var (
		workspaceID string
		expression  string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new schedule",
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

			req := saved.CreateScheduleRequest{
				Expression: expression,
			}

			resp, r, err := cliCtx.Client.SchedulesAPI.CreateSchedule(cliCtx.APICtx, workspaceID).CreateScheduleRequest(req).Execute()
			if err != nil {
				return internal.PrintAPIError(err)
			}
			defer r.Body.Close()

			color.Green("✓ Schedule created")
			color.Cyan("ID: %s", resp.GetId())
			return nil
		},
	}

	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	cmd.Flags().StringVarP(&expression, "expression", "e", "", "Cron expression")
	cmd.MarkFlagRequired("expression")

	return cmd
}
