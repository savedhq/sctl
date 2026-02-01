package agent

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/savedhq/sctl/internal"
	saved "github.com/savedhq/sdk-go"
	"github.com/spf13/cobra"
)

func NewAgentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "agent",
		Aliases: []string{"agents"},
		Short:   "Manage agents",
	}

	cmd.AddCommand(newAgentListCmd())
	cmd.AddCommand(newAgentGetCmd())
	cmd.AddCommand(newAgentCreateCmd())
	cmd.AddCommand(newAgentUpdateCmd())
	cmd.AddCommand(newAgentDeleteCmd())
	cmd.AddCommand(newAgentCredentialsCmd())

	return cmd
}

func newAgentListCmd() *cobra.Command {
	var workspaceID string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all agents",
		Run: func(cmd *cobra.Command, args []string) {
			cliCtx := internal.GetCLIContext(cmd.Context())
			internal.CheckErr(cliCtx.Err)

			var err error
			workspaceID, err = cliCtx.ResolveWorkspaceID(workspaceID)
			internal.CheckErr(err)

			resp, r, err := cliCtx.Client.AgentsAPI.ListAgents(cliCtx.APICtx, workspaceID).Execute()
			internal.CheckErr(internal.PrintAPIError(err))
			defer r.Body.Close()

			if cliCtx.JSONOutput {
				internal.PrintJSON(resp)
			} else {
				if len(resp) == 0 {
					color.Yellow("⚠ No agents found")
					return
				}

				for _, agent := range resp {
					color.Cyan("ID: %s", agent.GetId())
					fmt.Printf("  Name: %s\n", agent.GetName())
					fmt.Printf("  Status: %s\n\n", agent.GetStatus())
				}
			}
		},
	}
	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	return cmd
}

func newAgentGetCmd() *cobra.Command {
	var workspaceID string
	cmd := &cobra.Command{
		Use:   "get <agent_id>",
		Short: "Get agent details",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cliCtx := internal.GetCLIContext(cmd.Context())
			internal.CheckErr(cliCtx.Err)

			var err error
			workspaceID, err = cliCtx.ResolveWorkspaceID(workspaceID)
			internal.CheckErr(err)

			agentID, err := cliCtx.ResolveAgentID(workspaceID, args[0])
			internal.CheckErr(err)

			resp, r, err := cliCtx.Client.AgentsAPI.GetAgent(cliCtx.APICtx, workspaceID, agentID).Execute()
			internal.CheckErr(internal.PrintAPIError(err))
			defer r.Body.Close()

			if cliCtx.JSONOutput {
				internal.PrintJSON(resp)
			} else {
				color.Cyan("ID: %s", resp.GetId())
				fmt.Printf("Name: %s\n", resp.GetName())
				fmt.Printf("Status: %s\n", resp.GetStatus())
			}
		},
	}
	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	return cmd
}

func newAgentCreateCmd() *cobra.Command {
	var workspaceID, name string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new agent",
		Run: func(cmd *cobra.Command, args []string) {
			cliCtx := internal.GetCLIContext(cmd.Context())
			internal.CheckErr(cliCtx.Err)

			var err error
			workspaceID, err = cliCtx.ResolveWorkspaceID(workspaceID)
			internal.CheckErr(err)

			req := saved.CreateAgentRequest{}
			if name != "" {
				req.SetName(name)
			}

			resp, r, err := cliCtx.Client.AgentsAPI.CreateAgent(cliCtx.APICtx, workspaceID).CreateAgentRequest(req).Execute()
			internal.CheckErr(internal.PrintAPIError(err))
			defer r.Body.Close()

			if cliCtx.JSONOutput {
				internal.PrintJSON(resp)
			} else {
				color.Green("✓ Agent created")
				color.Cyan("ID: %s", resp.GetId())
			}
		},
	}
	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	cmd.Flags().StringVarP(&name, "name", "n", "", "Agent name")
	return cmd
}

func newAgentUpdateCmd() *cobra.Command {
	var workspaceID string
	cmd := &cobra.Command{
		Use:   "update <agent_id>",
		Short: "Update an agent",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cliCtx := internal.GetCLIContext(cmd.Context())
			internal.CheckErr(cliCtx.Err)

			var err error
			workspaceID, err = cliCtx.ResolveWorkspaceID(workspaceID)
			internal.CheckErr(err)

			agentID, err := cliCtx.ResolveAgentID(workspaceID, args[0])
			internal.CheckErr(err)

			_, r, err := cliCtx.Client.AgentsAPI.UpdateAgent(cliCtx.APICtx, workspaceID, agentID).Execute()
			internal.CheckErr(internal.PrintAPIError(err))
			defer r.Body.Close()

			if cliCtx.JSONOutput {
				internal.PrintJSON(map[string]string{"status": "ok"})
			} else {
				color.Green("✓ Agent updated")
			}
		},
	}
	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	return cmd
}

func newAgentDeleteCmd() *cobra.Command {
	var workspaceID string
	cmd := &cobra.Command{
		Use:   "delete <agent_id>",
		Short: "Delete an agent",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cliCtx := internal.GetCLIContext(cmd.Context())
			internal.CheckErr(cliCtx.Err)

			var err error
			workspaceID, err = cliCtx.ResolveWorkspaceID(workspaceID)
			internal.CheckErr(err)

			agentID, err := cliCtx.ResolveAgentID(workspaceID, args[0])
			internal.CheckErr(err)

			r, err := cliCtx.Client.AgentsAPI.DeleteAgent(cliCtx.APICtx, workspaceID, agentID).Execute()
			internal.CheckErr(internal.PrintAPIError(err))
			defer r.Body.Close()

			if cliCtx.JSONOutput {
				internal.PrintJSON(map[string]string{"status": "ok"})
			} else {
				color.Green("✓ Agent deleted")
			}
		},
	}
	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	return cmd
}

func newAgentCredentialsCmd() *cobra.Command {
	var workspaceID string
	cmd := &cobra.Command{
		Use:   "reset-credentials <agent_id>",
		Short: "Reset agent credentials",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cliCtx := internal.GetCLIContext(cmd.Context())
			internal.CheckErr(cliCtx.Err)

			var err error
			workspaceID, err = cliCtx.ResolveWorkspaceID(workspaceID)
			internal.CheckErr(err)

			agentID, err := cliCtx.ResolveAgentID(workspaceID, args[0])
			internal.CheckErr(err)

			resp, r, err := cliCtx.Client.AgentsAPI.ResetAgentCredentials(cliCtx.APICtx, workspaceID, agentID).Execute()
			internal.CheckErr(internal.PrintAPIError(err))
			defer r.Body.Close()

			if cliCtx.JSONOutput {
				internal.PrintJSON(resp)
			} else {
				color.Green("✓ Credentials reset")
			}
		},
	}
	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	return cmd
}
