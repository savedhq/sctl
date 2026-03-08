package job

import (
	"github.com/spf13/cobra"
)

func NewJobCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "job",
		Aliases: []string{"jobs"},
		Short:   "Manage jobs",
	}

	cmd.AddCommand(newJobListCmd())
	cmd.AddCommand(newJobGetCmd())
	cmd.AddCommand(newJobCreateWorkerCmd())
	cmd.AddCommand(newJobCreateAgentCmd())
	cmd.AddCommand(newJobCreateManualCmd())
	cmd.AddCommand(newJobUpdateCmd())
	cmd.AddCommand(newJobDeleteCmd())
	cmd.AddCommand(newJobTriggerCmd())

	return cmd
}
