package schedule

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/savedhq/sctl/internal"
	saved "github.com/savedhq/sdk-go"
	"github.com/spf13/cobra"
)

func newScheduleUpdateCmd() *cobra.Command {
	var (
		workspaceID string
		expression  string
	)

	cmd := &cobra.Command{
		Use:   "update <schedule_id>",
		Short: "Update a schedule",
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

			scheduleID := args[0]

			req := saved.UpdateScheduleRequest{}

			if cmd.Flags().Changed("expression") {
				req.SetExpression(expression)
			} else {
				color.Yellow("⚠ No fields to update. Use --expression to set a new cron expression.")
				return nil
			}

			_, r, err := cliCtx.Client.SchedulesAPI.UpdateSchedule(cliCtx.APICtx, workspaceID, scheduleID).UpdateScheduleRequest(req).Execute()
			if err != nil {
				return internal.PrintAPIError(err)
			}
			defer r.Body.Close()

			color.Green("✓ Schedule updated")
			return nil
		},
	}

	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	cmd.Flags().StringVarP(&expression, "expression", "e", "", "Cron expression")

	return cmd
}
