package workspace

import (
	"encoding/json"
	"fmt"

	"github.com/fatih/color"
	"github.com/savedhq/sctl/internal"
	saved "github.com/savedhq/sdk-go"
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
				fmt.Fprintln(cmd.OutOrStdout(), color.YellowString("⚠ No workspaces found"))
				return nil
			}

			for _, ws := range resp {
				fmt.Fprintln(cmd.OutOrStdout(), color.CyanString("ID: %s", ws.GetId()))
				fmt.Fprintf(cmd.OutOrStdout(), "  Name: %s\n", ws.GetName())
				fmt.Fprintf(cmd.OutOrStdout(), "  Created: %s\n\n", ws.GetCreatedAt().Format("2006-01-02 15:04:05"))
			}
			return nil
		},
	}
}

func newWorkspaceGetCmd() *cobra.Command {
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

			jsonOutput, _ := cmd.Flags().GetBool("json")
			if jsonOutput {
				data, _ := json.MarshalIndent(resp, "", "  ")
				fmt.Fprintln(cmd.OutOrStdout(), string(data))
			} else {
				fmt.Fprintln(cmd.OutOrStdout(), color.CyanString("ID: %s", resp.GetId()))
				fmt.Fprintf(cmd.OutOrStdout(), "Name: %s\n", resp.GetName())
				fmt.Fprintf(cmd.OutOrStdout(), "Created: %s\n", resp.GetCreatedAt().Format("2006-01-02 15:04:05"))
			}
			return nil
		},
	}
	cmd.Flags().Bool("json", false, "Output as JSON")
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

			req := cliCtx.Client.WorkspacesAPI.CreateWorkspace(cliCtx.APICtx)
			req = req.CreateWorkspaceRequest(saved.CreateWorkspaceRequest{
				Name: name,
			})

			resp, r, err := req.Execute()
			if err != nil {
				return fmt.Errorf("API error: %v", err)
			}
			defer r.Body.Close()

			fmt.Fprintln(cmd.OutOrStdout(), color.GreenString("✓ Workspace created"))
			fmt.Fprintln(cmd.OutOrStdout(), color.CyanString("ID: %s", resp.GetId()))
			return nil
		},
	}
	cmd.Flags().StringVarP(&name, "name", "n", "", "Workspace name")
	return cmd
}

func newWorkspaceUpdateCmd() *cobra.Command {
	var name string
	cmd := &cobra.Command{
		Use:   "update <workspace_id>",
		Short: "Update a workspace",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := internal.GetCLIContext(cmd.Context())
			if cliCtx == nil {
				return fmt.Errorf("CLI context not initialized")
			}

			if !cmd.Flags().Changed("name") {
				fmt.Fprintln(cmd.OutOrStdout(), "No changes specified. Use --name to set a new name.")
				return nil
			}

			id, err := cliCtx.ResolveWorkspaceID(args[0])
			if err != nil {
				return err
			}

			req := cliCtx.Client.WorkspacesAPI.UpdateWorkspace(cliCtx.APICtx, id)
			req = req.UpdateWorkspaceRequest(saved.UpdateWorkspaceRequest{
				Name: name,
			})

			_, r, err := req.Execute()
			if err != nil {
				return fmt.Errorf("API error: %v", err)
			}
			defer r.Body.Close()

			fmt.Fprintln(cmd.OutOrStdout(), color.GreenString("✓ Workspace updated"))
			return nil
		},
	}
	cmd.Flags().StringVarP(&name, "name", "n", "", "Workspace name")
	return cmd
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

			fmt.Fprintln(cmd.OutOrStdout(), color.GreenString("✓ Workspace deleted"))
			return nil
		},
	}
}
