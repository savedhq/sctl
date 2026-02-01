package schedule

import (
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
