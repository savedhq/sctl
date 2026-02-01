package workspace

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	saved "github.com/savedhq/sdk-go"
	"github.com/savedhq/sctl/internal"
	"github.com/savedhq/sctl/internal/render"
	"github.com/spf13/cobra"
)

// Workspace is a wrapper for the saved.ListWorkspaces200ResponseInner type for rendering.
type Workspace saved.ListWorkspaces200ResponseInner

// String implements the Stringer interface for Workspace.
func (w Workspace) String() string {
	return fmt.Sprintf("%s %s (%s)",
		color.CyanString("ID:"),
		w.Id,
		w.Name,
	)
}

// Workspaces is a wrapper for a slice of saved.ListWorkspaces200ResponseInner for rendering.
type Workspaces []saved.ListWorkspaces200ResponseInner

// String implements the Stringer interface for Workspaces.
func (w Workspaces) String() string {
	var b strings.Builder
	for _, ws := range w {
		b.WriteString(fmt.Sprintf("%s %s\n  Name: %s\n  Created: %s\n\n",
			color.CyanString("ID:"),
			ws.Id,
			ws.Name,
			ws.CreatedAt.Format(time.RFC3339),
		))
	}
	return strings.TrimSuffix(b.String(), "\n\n")
}

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
			resp, r, err := cliCtx.Client.WorkspacesAPI.ListWorkspaces(cliCtx.APICtx).Execute()
			if err != nil {
				return internal.PrintAPIError(err)
			}
			defer r.Body.Close()

			if len(resp) == 0 {
				render.Message(color.YellowString("⚠ No workspaces found"))
				return nil
			}

			render.Object(Workspaces(resp))
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
			id, err := cliCtx.ResolveWorkspaceID(args[0])
			if err != nil {
				return err
			}

			resp, r, err := cliCtx.Client.WorkspacesAPI.GetWorkspace(cliCtx.APICtx, id).Execute()
			if err != nil {
				return internal.PrintAPIError(err)
			}
			defer r.Body.Close()

			render.Object(Workspace(*resp))
			return nil
		},
	}
	return cmd
}

func newWorkspaceCreateCmd() *cobra.Command {
	var name string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new workspace",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := internal.GetCLIContext(cmd.Context())
			req := cliCtx.Client.WorkspacesAPI.CreateWorkspace(cliCtx.APICtx)
			if name != "" {
				req = req.CreateWorkspaceRequest(saved.CreateWorkspaceRequest{Name: name})
			}

			resp, r, err := req.Execute()
			if err != nil {
				return internal.PrintAPIError(err)
			}
			defer r.Body.Close()

			render.Message(color.GreenString("✓ Workspace created: %s", resp.GetId()))
			return nil
		},
	}
	cmd.Flags().StringVarP(&name, "name", "n", "", "Workspace name (optional)")
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
			id, err := cliCtx.ResolveWorkspaceID(args[0])
			if err != nil {
				return err
			}

			req := cliCtx.Client.WorkspacesAPI.UpdateWorkspace(cliCtx.APICtx, id)
			if name != "" {
				req = req.UpdateWorkspaceRequest(saved.UpdateWorkspaceRequest{Name: name})
			}

			_, r, err := req.Execute()
			if err != nil {
				return internal.PrintAPIError(err)
			}
			defer r.Body.Close()

			render.Message(color.GreenString("✓ Workspace updated"))
			return nil
		},
	}
	cmd.Flags().StringVarP(&name, "name", "n", "", "New workspace name")
	return cmd
}

func newWorkspaceDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <workspace_id>",
		Short: "Delete a workspace",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := internal.GetCLIContext(cmd.Context())
			id, err := cliCtx.ResolveWorkspaceID(args[0])
			if err != nil {
				return err
			}

			r, err := cliCtx.Client.WorkspacesAPI.DeleteWorkspace(cliCtx.APICtx, id).Execute()
			if err != nil {
				return internal.PrintAPIError(err)
			}
			defer r.Body.Close()

			render.Message(color.GreenString("✓ Workspace deleted"))
			return nil
		},
	}
}
