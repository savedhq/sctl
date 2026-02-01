package schedule

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/savedhq/sctl/internal"
	"github.com/spf13/cobra"
)

func newScheduleGetCmd() *cobra.Command {
	var (
		workspaceID string
		jsonOutput  bool
	)

	cmd := &cobra.Command{
		Use:   "get <schedule_id>",
		Short: "Get schedule details",
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

			resp, r, err := cliCtx.Client.SchedulesAPI.GetSchedule(cliCtx.APICtx, workspaceID, scheduleID).Execute()
			if err != nil {
				return internal.PrintAPIError(err)
			}
			defer r.Body.Close()

			if jsonOutput {
				data, err := json.MarshalIndent(resp, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal json: %w", err)
				}
				fmt.Println(string(data))
				return nil
			}

			color.Cyan("ID: %s", resp.GetId())
			fmt.Printf("  Cron Expression: %s\n", resp.GetExpression())
			fmt.Printf("  Next Run: %s\n", resp.GetNextRun())
			fmt.Printf("  Last Run: %s\n", resp.GetLastRun())

			return nil
		},
	}

	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output as JSON")

	return cmd
}
