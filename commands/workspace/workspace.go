package workspace

import (
	"encoding/json"
	"fmt"

	"github.com/fatih/color"
	"github.com/savedhq/sctl/internal"
	"github.com/spf13/cobra"
)

func NewWorkspaceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "workspace",
		Aliases: []string{"ws", "workspaces"},
		Short:   "Manage workspaces",
	}

	cmd.AddCommand(newWorkspaceListCmd())
	cmd.AddCommand(newWorkspaceGetCmd())
	cmd.AddCommand(newWorkspaceCreateCmd())
	cmd.AddCommand(newWorkspaceUpdateCmd())
	cmd.AddCommand(newWorkspaceDeleteCmd())

	return cmd
}

func newWorkspaceListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all workspaces",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := internal.GetCLIContext(cmd.Context())
			if cliCtx == nil {
				return fmt.Errorf("CLI context not initialized")
			}

			resp, r, err := cliCtx.Client.WorkspacesAPI.ListWorkspaces(cliCtx.APICtx).Execute()
			if err != nil {
				return fmt.Errorf("API error: %v", err)
			}
			defer r.Body.Close()

			if len(resp) == 0 {
				color.Yellow("⚠ No workspaces found")
				return nil
			}

			for _, ws := range resp {
				color.Cyan("ID: %s", ws.GetId())
				fmt.Printf("  Name: %s\n", ws.GetName())
				fmt.Printf("  Created: %s\n\n", ws.GetCreatedAt().Format("2006-01-02 15:04:05"))
			}
			return nil
		},
	}
}

func newWorkspaceGetCmd() *cobra.Command {
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "get <workspace_id>",
		Short: "Get workspace details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := internal.GetCLIContext(cmd.Context())
			if cliCtx == nil {
				return fmt.Errorf("CLI context not initialized")
			}

			id, err := cliCtx.ResolveWorkspaceID(args[0])
			if err != nil {
				return err
			}

			resp, r, err := cliCtx.Client.WorkspacesAPI.GetWorkspace(cliCtx.APICtx, id).Execute()
			if err != nil {
				return fmt.Errorf("API error: %v", err)
			}
			defer r.Body.Close()

			if jsonOutput {
				data, _ := json.MarshalIndent(resp, "", "  ")
				fmt.Println(string(data))
			} else {
				color.Cyan("ID: %s", resp.GetId())
				fmt.Printf("Name: %s\n", resp.GetName())
				fmt.Printf("Created: %s\n", resp.GetCreatedAt().Format("2006-01-02 15:04:05"))
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output as JSON")
	return cmd
}

func newWorkspaceCreateCmd() *cobra.Command {
	var name string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new workspace",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := internal.GetCLIContext(cmd.Context())
			if cliCtx == nil {
				return fmt.Errorf("CLI context not initialized")
			}

			resp, r, err := cliCtx.Client.WorkspacesAPI.CreateWorkspace(cliCtx.APICtx).Execute()
			if err != nil {
				return fmt.Errorf("API error: %v", err)
			}
			defer r.Body.Close()

			color.Green("✓ Workspace created")
			color.Cyan("ID: %s", resp.GetId())
			return nil
		},
	}
	cmd.Flags().StringVarP(&name, "name", "n", "", "Workspace name")
	return cmd
}

func newWorkspaceUpdateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update <workspace_id>",
		Short: "Update a workspace",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := internal.GetCLIContext(cmd.Context())
			if cliCtx == nil {
				return fmt.Errorf("CLI context not initialized")
			}

			id, err := cliCtx.ResolveWorkspaceID(args[0])
			if err != nil {
				return err
			}

			_, r, err := cliCtx.Client.WorkspacesAPI.UpdateWorkspace(cliCtx.APICtx, id).Execute()
			if err != nil {
				return fmt.Errorf("API error: %v", err)
			}
			defer r.Body.Close()

			color.Green("✓ Workspace updated")
			return nil
		},
	}
}

func newWorkspaceDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <workspace_id>",
		Short: "Delete a workspace",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := internal.GetCLIContext(cmd.Context())
			if cliCtx == nil {
				return fmt.Errorf("CLI context not initialized")
			}

			id, err := cliCtx.ResolveWorkspaceID(args[0])
			if err != nil {
				return err
			}

			r, err := cliCtx.Client.WorkspacesAPI.DeleteWorkspace(cliCtx.APICtx, id).Execute()
			if err != nil {
				return fmt.Errorf("API error: %v", err)
			}
			defer r.Body.Close()

			color.Green("✓ Workspace deleted")
			return nil
		},
	}
}
