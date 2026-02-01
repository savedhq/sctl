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
				fmt.Fprintln(cmd.OutOrStdout(), color.YellowString("⚠ No agents found"))
				return nil
			}

			for _, agent := range resp {
				fmt.Fprintln(cmd.OutOrStdout(), color.CyanString("ID: %s", agent.GetId()))
				fmt.Fprintf(cmd.OutOrStdout(), "  Name: %s\n", agent.GetName())
				fmt.Fprintf(cmd.OutOrStdout(), "  Status: %s\n\n", agent.GetStatus())
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
				fmt.Fprintln(cmd.OutOrStdout(), string(data))
			} else {
				fmt.Fprintln(cmd.OutOrStdout(), color.CyanString("ID: %s", resp.GetId()))
				fmt.Fprintf(cmd.OutOrStdout(), "Name: %s\n", resp.GetName())
				fmt.Fprintf(cmd.OutOrStdout(), "Status: %s\n", resp.GetStatus())
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

			fmt.Fprintln(cmd.OutOrStdout(), color.GreenString("✓ Agent created"))
			fmt.Fprintln(cmd.OutOrStdout(), color.CyanString("ID: %s", resp.GetId()))
			return nil
		},
	}
	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	cmd.Flags().StringVarP(&name, "name", "n", "", "Agent name")
	return cmd
}

func newAgentUpdateCmd() *cobra.Command {
	var workspaceID, name string
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

			req := saved.UpdateAgentRequest{}
			if name != "" {
				req.SetName(name)
			}

			_, r, err := cliCtx.Client.AgentsAPI.UpdateAgent(cliCtx.APICtx, workspaceID, agentID).UpdateAgentRequest(req).Execute()
			if err != nil {
				return internal.PrintAPIError(err)
			}
			defer r.Body.Close()

			fmt.Fprintln(cmd.OutOrStdout(), color.GreenString("✓ Agent updated"))
			return nil
		},
	}
	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	cmd.Flags().StringVarP(&name, "name", "n", "", "Agent name")
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

			fmt.Fprintln(cmd.OutOrStdout(), color.GreenString("✓ Agent deleted"))
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

			fmt.Fprintln(cmd.OutOrStdout(), color.GreenString("✓ Credentials reset"))
			data, _ := json.MarshalIndent(resp, "", "  ")
			fmt.Fprintln(cmd.OutOrStdout(), string(data))
			return nil
		},
	}
	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	return cmd
}
