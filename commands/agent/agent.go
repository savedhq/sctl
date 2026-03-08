package agent

import (
	"encoding/json"
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

			resp, r, err := cliCtx.Client.AgentsAPI.ListAgents(cliCtx.APICtx, workspaceID).Execute()
			if err != nil {
				return internal.PrintAPIError(err)
			}
			defer r.Body.Close()

			if len(resp) == 0 {
				color.Yellow("⚠ No agents found")
				return nil
			}

			for _, agent := range resp {
				color.Cyan("ID: %s", agent.GetId())
				fmt.Printf("  Name: %s\n", agent.GetName())
				fmt.Printf("  Status: %s\n\n", agent.GetStatus())
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	return cmd
}

func newAgentGetCmd() *cobra.Command {
	var workspaceID string
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "get <agent_id>",
		Short: "Get agent details",
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

			agentID, err := cliCtx.ResolveAgentID(workspaceID, args[0])
			if err != nil {
				return err
			}

			resp, r, err := cliCtx.Client.AgentsAPI.GetAgent(cliCtx.APICtx, workspaceID, agentID).Execute()
			if err != nil {
				return internal.PrintAPIError(err)
			}
			defer r.Body.Close()

			if jsonOutput {
				data, _ := json.MarshalIndent(resp, "", "  ")
				fmt.Println(string(data))
			} else {
				color.Cyan("ID: %s", resp.GetId())
				fmt.Printf("Name: %s\n", resp.GetName())
				fmt.Printf("Status: %s\n", resp.GetStatus())
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output as JSON")
	return cmd
}

func newAgentCreateCmd() *cobra.Command {
	var workspaceID, name string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new agent",
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

			req := saved.CreateAgentRequest{}
			if name != "" {
				req.SetName(name)
			}

			resp, r, err := cliCtx.Client.AgentsAPI.CreateAgent(cliCtx.APICtx, workspaceID).CreateAgentRequest(req).Execute()
			if err != nil {
				return internal.PrintAPIError(err)
			}
			defer r.Body.Close()

			color.Green("✓ Agent created")
			color.Cyan("ID: %s", resp.GetId())
			return nil
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

			agentID, err := cliCtx.ResolveAgentID(workspaceID, args[0])
			if err != nil {
				return err
			}

			_, r, err := cliCtx.Client.AgentsAPI.UpdateAgent(cliCtx.APICtx, workspaceID, agentID).Execute()
			if err != nil {
				return internal.PrintAPIError(err)
			}
			defer r.Body.Close()

			color.Green("✓ Agent updated")
			return nil
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

			agentID, err := cliCtx.ResolveAgentID(workspaceID, args[0])
			if err != nil {
				return err
			}

			r, err := cliCtx.Client.AgentsAPI.DeleteAgent(cliCtx.APICtx, workspaceID, agentID).Execute()
			if err != nil {
				return internal.PrintAPIError(err)
			}
			defer r.Body.Close()

			color.Green("✓ Agent deleted")
			return nil
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

			agentID, err := cliCtx.ResolveAgentID(workspaceID, args[0])
			if err != nil {
				return err
			}

			resp, r, err := cliCtx.Client.AgentsAPI.ResetAgentCredentials(cliCtx.APICtx, workspaceID, agentID).Execute()
			if err != nil {
				return internal.PrintAPIError(err)
			}
			defer r.Body.Close()

			color.Green("✓ Credentials reset")
			data, _ := json.MarshalIndent(resp, "", "  ")
			fmt.Println(string(data))
			return nil
		},
	}
	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	return cmd
}
