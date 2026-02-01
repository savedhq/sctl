package schedule

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/savedhq/sctl/internal"
	"github.com/spf13/cobra"
)

func newScheduleDeleteCmd() *cobra.Command {
	var workspaceID string
	cmd := &cobra.Command{
		Use:   "delete <schedule_id>",
		Short: "Delete a schedule",
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

			r, err := cliCtx.Client.SchedulesAPI.DeleteSchedule(cliCtx.APICtx, workspaceID, scheduleID).Execute()
			if err != nil {
				return internal.PrintAPIError(err)
			}
			defer r.Body.Close()

			color.Green("✓ Schedule deleted")
			return nil
		},
	}
	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	return cmd
}
