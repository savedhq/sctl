package schedule

import (
	"fmt"

	"github.com/savedhq/sctl/internal"
	"github.com/spf13/cobra"
)

func NewScheduleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "schedule",
		Aliases: []string{"schedules"},
		Short:   "Manage schedules",
	}

	cmd.AddCommand(newScheduleListCmd())
	cmd.AddCommand(newScheduleGetCmd())
	cmd.AddCommand(newScheduleCreateCmd())
	cmd.AddCommand(newScheduleUpdateCmd())
	cmd.AddCommand(newScheduleDeleteCmd())
	cmd.AddCommand(newSchedulePauseCmd())
	cmd.AddCommand(newScheduleResumeCmd())

	return cmd
}

func newScheduleListCmd() *cobra.Command {
	var workspaceID string
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

			// TODO: Uncomment when SDK is updated with SchedulesAPI
			// resp, r, err := cliCtx.Client.SchedulesAPI.ListSchedules(cliCtx.APICtx, workspaceID).Execute()
			// if err != nil {
			// 	return internal.PrintAPIError(err)
			// }
			// defer r.Body.Close()

			// if len(resp) == 0 {
			// 	color.Yellow("⚠ No schedules found")
			// 	return nil
			// }

			// for _, schedule := range resp {
			// 	color.Cyan("ID: %s", schedule.GetId())
			// 	fmt.Printf("  Cron Expression: %s\n", schedule.GetCronExpression())
			// 	fmt.Printf("  Status: %s\n\n", schedule.GetStatus())
			// }
			return nil
		},
	}
	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	return cmd
}

func newScheduleGetCmd() *cobra.Command {
	var workspaceID string
	var jsonOutput bool
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

			// TODO: Uncomment when SDK is updated with SchedulesAPI
			// scheduleID, err := cliCtx.ResolveScheduleID(workspaceID, args[0])
			// if err != nil {
			// 	return err
			// }

			// resp, r, err := cliCtx.Client.SchedulesAPI.GetSchedule(cliCtx.APICtx, workspaceID, scheduleID).Execute()
			// if err != nil {
			// 	return internal.PrintAPIError(err)
			// }
			// defer r.Body.Close()

			// if jsonOutput {
			// 	data, _ := json.MarshalIndent(resp, "", "  ")
			// 	fmt.Println(string(data))
			// } else {
			// 	color.Cyan("ID: %s", resp.GetId())
			// 	fmt.Printf("  Cron Expression: %s\n", resp.GetCronExpression())
			// 	fmt.Printf("  Next Run: %s\n", resp.GetNextRun().String())
			// 	fmt.Printf("  Last Run: %s\n", resp.GetLastRun().String())
			// 	fmt.Printf("  Status: %s\n", resp.GetStatus())
			// }
			return nil
		},
	}
	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output as JSON")
	return cmd
}

func newScheduleCreateCmd() *cobra.Command {
	var workspaceID, cronExpression string
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

			// TODO: Uncomment when SDK is updated with SchedulesAPI
			// req := saved.CreateScheduleRequest{}
			// if cronExpression != "" {
			// 	req.SetCronExpression(cronExpression)
			// }

			// resp, r, err := cliCtx.Client.SchedulesAPI.CreateSchedule(cliCtx.APICtx, workspaceID).CreateScheduleRequest(req).Execute()
			// if err != nil {
			// 	return internal.PrintAPIError(err)
			// }
			// defer r.Body.Close()

			// color.Green("✓ Schedule created")
			// color.Cyan("ID: %s", resp.GetId())
			return nil
		},
	}
	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	cmd.Flags().StringVar(&cronExpression, "cron", "", "Cron expression")
	cmd.MarkFlagRequired("cron")
	return cmd
}

func newScheduleUpdateCmd() *cobra.Command {
	var workspaceID, cronExpression string
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

			// TODO: Uncomment when SDK is updated with SchedulesAPI
			// scheduleID, err := cliCtx.ResolveScheduleID(workspaceID, args[0])
			// if err != nil {
			// 	return err
			// }

			// req := saved.UpdateScheduleRequest{}
			// if cronExpression != "" {
			// 	req.SetCronExpression(cronExpression)
			// }

			// _, r, err := cliCtx.Client.SchedulesAPI.UpdateSchedule(cliCtx.APICtx, workspaceID, scheduleID).UpdateScheduleRequest(req).Execute()
			// if err != nil {
			// 	return internal.PrintAPIError(err)
			// }
			// defer r.Body.Close()

			// color.Green("✓ Schedule updated")
			return nil
		},
	}
	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	cmd.Flags().StringVar(&cronExpression, "cron", "", "Cron expression")
	return cmd
}

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

			// TODO: Uncomment when SDK is updated with SchedulesAPI
			// scheduleID, err := cliCtx.ResolveScheduleID(workspaceID, args[0])
			// if err != nil {
			// 	return err
			// }

			// r, err := cliCtx.Client.SchedulesAPI.DeleteSchedule(cliCtx.APICtx, workspaceID, scheduleID).Execute()
			// if err != nil {
			// 	return internal.PrintAPIError(err)
			// }
			// defer r.Body.Close()

			// color.Green("✓ Schedule deleted")
			return nil
		},
	}
	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	return cmd
}

func newSchedulePauseCmd() *cobra.Command {
	var workspaceID string
	cmd := &cobra.Command{
		Use:   "pause <schedule_id>",
		Short: "Pause a schedule",
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

			// TODO: Uncomment when SDK is updated with SchedulesAPI
			// scheduleID, err := cliCtx.ResolveScheduleID(workspaceID, args[0])
			// if err != nil {
			// 	return err
			// }

			// r, err := cliCtx.Client.SchedulesAPI.PauseSchedule(cliCtx.APICtx, workspaceID, scheduleID).Execute()
			// if err != nil {
			// 	return internal.PrintAPIError(err)
			// }
			// defer r.Body.Close()

			// color.Green("✓ Schedule paused")
			return nil
		},
	}
	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	return cmd
}

func newScheduleResumeCmd() *cobra.Command {
	var workspaceID string
	cmd := &cobra.Command{
		Use:   "resume <schedule_id>",
		Short: "Resume a schedule",
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

			// TODO: Uncomment when SDK is updated with SchedulesAPI
			// scheduleID, err := cliCtx.ResolveScheduleID(workspaceID, args[0])
			// if err != nil {
			// 	return err
			// }

			// r, err := cliCtx.Client.SchedulesAPI.ResumeSchedule(cliCtx.APICtx, workspaceID, scheduleID).Execute()
			// if err != nil {
			// 	return internal.PrintAPIError(err)
			// }
			// defer r.Body.Close()

			// color.Green("✓ Schedule resumed")
			return nil
		},
	}
	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	return cmd
}