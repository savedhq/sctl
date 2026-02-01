package schedule

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/savedhq/sctl/internal"
	"github.com/spf13/cobra"
)

func newScheduleListCmd() *cobra.Command {
	var (
		workspaceID string
		jsonOutput  bool
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all schedules",
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

			resp, r, err := cliCtx.Client.SchedulesAPI.ListSchedules(cliCtx.APICtx, workspaceID).Execute()
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

			if len(resp) == 0 {
				color.Yellow("⚠ No schedules found")
				return nil
			}

			for _, schedule := range resp {
				color.Cyan("ID: %s", schedule.GetId())
				fmt.Printf("  Cron Expression: %s\n", schedule.GetExpression())
				fmt.Printf("  Next Run: %s\n", schedule.GetNextRun())
				fmt.Printf("  Last Run: %s\n\n", schedule.GetLastRun())
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output as JSON")

	return cmd
}
