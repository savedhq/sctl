package agent

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/savedhq/sctl/internal"
	"github.com/savedhq/sctl/internal/render"
	saved "github.com/savedhq/sdk-go"
	"github.com/spf13/cobra"
)

// Agent is a wrapper for the saved.ListAgents200ResponseInner type for rendering.
type Agent saved.ListAgents200ResponseInner

// String implements the Stringer interface for Agent.
func (a Agent) String() string {
	return fmt.Sprintf("%s %s (%s)",
		color.CyanString("ID:"),
		a.Id,
		a.Name,
	)
}

// Agents is a wrapper for a slice of saved.ListAgents200ResponseInner for rendering.
type Agents []saved.ListAgents200ResponseInner

// String implements the Stringer interface for Agents.
func (a Agents) String() string {
	var b strings.Builder
	for _, agent := range a {
		b.WriteString(fmt.Sprintf("%s %s\n  Name: %s\n  Status: %s\n\n",
			color.CyanString("ID:"),
			agent.Id,
			agent.Name,
			agent.Status,
		))
	}
	return strings.TrimSuffix(b.String(), "\n\n")
}

// AgentCredentials is a wrapper for the saved.AgentCredentials type for rendering.
type AgentCredentials saved.AgentCredentials

// String implements the Stringer interface for AgentCredentials.
func (c AgentCredentials) String() string {
	data, _ := json.MarshalIndent(c, "", "  ")
	return string(data)
}

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
				render.Message(color.YellowString("⚠ No agents found"))
				return nil
			}

			render.Object(Agents(resp))
			return nil
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
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := internal.GetCLIContext(cmd.Context())
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

			render.Object(Agent(*resp))
			return nil
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
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := internal.GetCLIContext(cmd.Context())
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

			render.Message(color.GreenString("✓ Agent created: %s", resp.GetId()))
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
			var err error
			workspaceID, err = cliCtx.ResolveWorkspaceID(workspaceID)
			if err != nil {
				return err
			}

			agentID, err := cliCtx.ResolveAgentID(workspaceID, args[0])
			if err != nil {
				return err
			}

			req := cliCtx.Client.AgentsAPI.UpdateAgent(cliCtx.APICtx, workspaceID, agentID)
			if name != "" {
				req = req.UpdateAgentRequest(saved.UpdateAgentRequest{Name: &name})
			}

			_, r, err := req.Execute()
			if err != nil {
				return internal.PrintAPIError(err)
			}
			defer r.Body.Close()

			render.Message(color.GreenString("✓ Agent updated"))
			return nil
		},
	}
	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	cmd.Flags().StringVarP(&name, "name", "n", "", "New agent name")
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

			render.Message(color.GreenString("✓ Agent deleted"))
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

			render.Message(color.GreenString("✓ Credentials reset"))
			render.Object(AgentCredentials(*resp))
			return nil
		},
	}
	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	return cmd
}
