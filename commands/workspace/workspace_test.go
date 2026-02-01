package workspace_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fatih/color"
	"github.com/savedhq/sctl/commands/workspace"
	"github.com/savedhq/sctl/internal"
	saved "github.com/savedhq/sdk-go"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// newTestCLIContext creates a new CLIContext for testing.
func newTestCLIContext(client *saved.APIClient, apiCtx context.Context) context.Context {
	cliCtx := &internal.CLIContext{
		Config: &internal.Config{},
		Client: client,
		APICtx: apiCtx,
	}
	return internal.WithCLIContext(context.Background(), cliCtx)
}

// setupMockServer creates a new httptest.Server and a corresponding API client.
func setupMockServer(t *testing.T, handler http.Handler) (*httptest.Server, *saved.APIClient, context.Context) {
	server := httptest.NewServer(handler)

	cfg := saved.NewConfiguration()
	cfg.Servers = saved.ServerConfigurations{{URL: server.URL}}
	client := saved.NewAPIClient(cfg)
	ctx := context.Background()

	return server, client, ctx
}

// executeCmd runs the command and returns the output.
func executeCmd(ctx context.Context, cmd *cobra.Command, args ...string) (string, error) {
	buf := new(bytes.Buffer)
	rootCmd := &cobra.Command{Use: "sctl"}
	rootCmd.AddCommand(cmd)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(append([]string{cmd.Use}, args...))
	rootCmd.SetContext(ctx)

	err := rootCmd.Execute()
	return buf.String(), err
}

func TestWorkspaceList(t *testing.T) {
	color.NoColor = true // Disable color for testing.

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v1/workspaces", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[
			{
				"id": "ws_1",
				"name": "Workspace 1",
				"created_at": "2006-01-02T15:04:05Z"
			},
			{
				"id": "ws_2",
				"name": "Workspace 2",
				"created_at": "2006-01-02T15:04:05Z"
			}
		]`))
	})

	server, client, ctx := setupMockServer(t, handler)
	defer server.Close()

	ctx = newTestCLIContext(client, ctx)
	output, err := executeCmd(ctx, workspace.NewWorkspaceCmd(), "list")

	assert.NoError(t, err)
	assert.Contains(t, output, "ID: ws_1")
	assert.Contains(t, output, "Name: Workspace 1")
	assert.Contains(t, output, "ID: ws_2")
	assert.Contains(t, output, "Name: Workspace 2")
}

func TestWorkspaceGet(t *testing.T) {
	color.NoColor = true // Disable color for testing.
	workspaceID := "11111111-1111-1111-1111-111111111111"

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v1/workspaces/"+workspaceID, r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": "` + workspaceID + `",
			"name": "Workspace 1",
			"created_at": "2006-01-02T15:04:05Z"
		}`))
	})

	server, client, ctx := setupMockServer(t, handler)
	defer server.Close()

	ctx = newTestCLIContext(client, ctx)

	// Test with normal output
	output, err := executeCmd(ctx, workspace.NewWorkspaceCmd(), "get", workspaceID)
	assert.NoError(t, err)
	assert.Contains(t, output, "ID: "+workspaceID)
	assert.Contains(t, output, "Name: Workspace 1")

	// Test with JSON output
	output, err = executeCmd(ctx, workspace.NewWorkspaceCmd(), "get", workspaceID, "--json")
	assert.NoError(t, err)
	assert.JSONEq(t, `{
		"id": "`+workspaceID+`",
		"name": "Workspace 1",
		"created_at": "2006-01-02T15:04:05Z"
	}`, output)
}

func TestWorkspaceCreate(t *testing.T) {
	color.NoColor = true // Disable color for testing.

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v1/workspaces", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{
			"id": "ws_new",
			"name": "New Workspace",
			"created_at": "2006-01-02T15:04:05Z"
		}`))
	})

	server, client, ctx := setupMockServer(t, handler)
	defer server.Close()

	ctx = newTestCLIContext(client, ctx)
	output, err := executeCmd(ctx, workspace.NewWorkspaceCmd(), "create", "--name", "New Workspace")

	assert.NoError(t, err)
	assert.Contains(t, output, "✓ Workspace created")
	assert.Contains(t, output, "ID: ws_new")
}

func TestWorkspaceUpdate(t *testing.T) {
	color.NoColor = true // Disable color for testing.
	workspaceID := "11111111-1111-1111-1111-111111111111"

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v1/workspaces/"+workspaceID, r.URL.Path)
		assert.Equal(t, http.MethodPatch, r.Method)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(``))
	})

	server, client, ctx := setupMockServer(t, handler)
	defer server.Close()

	ctx = newTestCLIContext(client, ctx)
	output, err := executeCmd(ctx, workspace.NewWorkspaceCmd(), "update", workspaceID, "--name", "Updated Workspace")

	assert.NoError(t, err)
	assert.Contains(t, output, "✓ Workspace updated")
}

func TestWorkspaceDelete(t *testing.T) {
	color.NoColor = true // Disable color for testing.
	workspaceID := "11111111-1111-1111-1111-111111111111"

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v1/workspaces/"+workspaceID, r.URL.Path)
		assert.Equal(t, http.MethodDelete, r.Method)

		w.WriteHeader(http.StatusNoContent)
	})

	server, client, ctx := setupMockServer(t, handler)
	defer server.Close()

	ctx = newTestCLIContext(client, ctx)
	output, err := executeCmd(ctx, workspace.NewWorkspaceCmd(), "delete", workspaceID)

	assert.NoError(t, err)
	assert.Contains(t, output, "✓ Workspace deleted")
}
