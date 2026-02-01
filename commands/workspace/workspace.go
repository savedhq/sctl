package workspace

import (
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
		Run: func(cmd *cobra.Command, args []string) {
			cliCtx := internal.GetCLIContext(cmd.Context())
			internal.CheckErr(cliCtx.Err)

			resp, r, err := cliCtx.Client.WorkspacesAPI.ListWorkspaces(cliCtx.APICtx).Execute()
			internal.CheckErr(internal.PrintAPIError(err))
			defer r.Body.Close()

			if cliCtx.JSONOutput {
				internal.PrintJSON(resp)
			} else {
				if len(resp) == 0 {
					color.Yellow("⚠ No workspaces found")
					return
				}

				for _, ws := range resp {
					color.Cyan("ID: %s", ws.GetId())
					fmt.Printf("  Name: %s\n", ws.GetName())
					fmt.Printf("  Created: %s\n\n", ws.GetCreatedAt().Format("2006-01-02 15:04:05"))
				}
			}
		},
	}
}

func newWorkspaceGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <workspace_id>",
		Short: "Get workspace details",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cliCtx := internal.GetCLIContext(cmd.Context())
			internal.CheckErr(cliCtx.Err)

			id, err := cliCtx.ResolveWorkspaceID(args[0])
			internal.CheckErr(err)

			resp, r, err := cliCtx.Client.WorkspacesAPI.GetWorkspace(cliCtx.APICtx, id).Execute()
			internal.CheckErr(internal.PrintAPIError(err))
			defer r.Body.Close()

			if cliCtx.JSONOutput {
				internal.PrintJSON(resp)
			} else {
				color.Cyan("ID: %s", resp.GetId())
				fmt.Printf("Name: %s\n", resp.GetName())
				fmt.Printf("Created: %s\n", resp.GetCreatedAt().Format("2006-01-02 15:04:05"))
			}
		},
	}
	return cmd
}

func newWorkspaceCreateCmd() *cobra.Command {
	var name string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new workspace",
		Run: func(cmd *cobra.Command, args []string) {
			cliCtx := internal.GetCLIContext(cmd.Context())
			internal.CheckErr(cliCtx.Err)

			resp, r, err := cliCtx.Client.WorkspacesAPI.CreateWorkspace(cliCtx.APICtx).Execute()
			internal.CheckErr(internal.PrintAPIError(err))
			defer r.Body.Close()

			if cliCtx.JSONOutput {
				internal.PrintJSON(resp)
			} else {
				color.Green("✓ Workspace created")
				color.Cyan("ID: %s", resp.GetId())
			}
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
		Run: func(cmd *cobra.Command, args []string) {
			cliCtx := internal.GetCLIContext(cmd.Context())
			internal.CheckErr(cliCtx.Err)

			id, err := cliCtx.ResolveWorkspaceID(args[0])
			internal.CheckErr(err)

			_, r, err := cliCtx.Client.WorkspacesAPI.UpdateWorkspace(cliCtx.APICtx, id).Execute()
			internal.CheckErr(internal.PrintAPIError(err))
			defer r.Body.Close()

			if cliCtx.JSONOutput {
				internal.PrintJSON(map[string]string{"status": "ok"})
			} else {
				color.Green("✓ Workspace updated")
			}
		},
	}
}

func newWorkspaceDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <workspace_id>",
		Short: "Delete a workspace",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cliCtx := internal.GetCLIContext(cmd.Context())
			internal.CheckErr(cliCtx.Err)

			id, err := cliCtx.ResolveWorkspaceID(args[0])
			internal.CheckErr(err)

			r, err := cliCtx.Client.WorkspacesAPI.DeleteWorkspace(cliCtx.APICtx, id).Execute()
			internal.CheckErr(internal.PrintAPIError(err))
			defer r.Body.Close()

			if cliCtx.JSONOutput {
				internal.PrintJSON(map[string]string{"status": "ok"})
			} else {
				color.Green("✓ Workspace deleted")
			}
		},
	}
}
