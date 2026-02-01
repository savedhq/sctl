package job

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/savedhq/sctl/internal"
	saved "github.com/savedhq/sdk-go"
	"github.com/stretchr/testify/assert"
)

// newTestCLIContext creates a new CLIContext for testing purposes.
func newTestCLIContext(serverURL string) *internal.CLIContext {
	cfg := &internal.Config{
		ServerURL: serverURL,
		APIKey:    "test-key",
	}
	client, ctx := cfg.GetClient()
	return &internal.CLIContext{
		Config: cfg,
		Client: client,
		APICtx: ctx,
	}
}

func TestJobCommands(t *testing.T) {
	workspaces := []saved.ListWorkspaces200ResponseInner{
		{Id: saved.PtrString("ws-123"), Name: saved.PtrString("Workspace 1")},
	}
	jobs := []saved.ListJobs200ResponseInner{
		{Id: saved.PtrString("job-1"), Name: saved.PtrString("Job 1"), Type: saved.PtrString("worker"), Enabled: saved.PtrBool(true)},
		{Id: saved.PtrString("job-2"), Name: saved.PtrString("Job 2"), Type: saved.PtrString("agent"), Enabled: saved.PtrBool(false)},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/v1/workspaces":
			json.NewEncoder(w).Encode(workspaces)
		case r.Method == http.MethodGet && r.URL.Path == "/v1/workspaces/ws-123/jobs":
			json.NewEncoder(w).Encode(jobs)
		case r.Method == http.MethodGet && r.URL.Path == "/v1/workspaces/ws-123/jobs/job-1":
			json.NewEncoder(w).Encode(jobs[0])
		case r.Method == http.MethodPost && r.URL.Path == "/v1/workspaces/ws-123/jobs/worker":
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(saved.CreateWorkerJob201Response{Id: saved.PtrString("job-123")})
		case r.Method == http.MethodPatch && r.URL.Path == "/v1/workspaces/ws-123/jobs/job-1":
			w.WriteHeader(http.StatusOK)
		case r.Method == http.MethodDelete && r.URL.Path == "/v1/workspaces/ws-123/jobs/job-1":
			w.WriteHeader(http.StatusNoContent)
		case r.Method == http.MethodPost && r.URL.Path == "/v1/workspaces/ws-123/jobs/job-1/trigger":
			json.NewEncoder(w).Encode(saved.TriggerJob202Response{WorkflowId: saved.PtrString("run-123")})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	tests := []struct {
		name      string
		command   string
		args      []string
		expected  string
		expectErr bool
	}{
		{
			name:     "create worker job",
			command:  "create-worker",
			args:     []string{"--workspace", "Workspace 1", "--name", "test-worker", "--schedule", "0 0 * * *"},
			expected: "✓ Worker job created\nID: job-123\n",
		},
		{
			name:    "list jobs",
			command: "list",
			args:    []string{"--workspace", "Workspace 1"},
			expected: "ID: job-1\n" +
				"  Name: Job 1\n" +
				"  Type: worker\n" +
				"  Enabled: true\n\n" +
				"ID: job-2\n" +
				"  Name: Job 2\n" +
				"  Type: agent\n" +
				"  Enabled: false\n\n",
		},
		{
			name:    "get job",
			command: "get",
			args:    []string{"Job 1", "--workspace", "Workspace 1"},
			expected: "ID: job-1\n" +
				"Name: Job 1\n" +
				"Type: worker\n" +
				"Enabled: true\n",
		},
		{
			name:     "update job",
			command:  "update",
			args:     []string{"Job 1", "--workspace", "Workspace 1", "--name", "new-name", "--schedule", "0 1 * * *"},
			expected: "✓ Job updated\n",
		},
		{
			name:     "delete job",
			command:  "delete",
			args:     []string{"Job 1", "--workspace", "Workspace 1"},
			expected: "✓ Job deleted\n",
		},
		{
			name:    "trigger job",
			command: "trigger",
			args:    []string{"Job 1", "--workspace", "Workspace 1"},
			expected: "✓ Job triggered\n" +
				"{\n" +
				"  \"workflow_id\": \"run-123\"\n" +
				"}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cliCtx := newTestCLIContext(server.URL)
			cmd := NewJobCmd()
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)
			cmd.SetArgs(append([]string{tt.command}, tt.args...))
			ctx := internal.WithCLIContext(context.Background(), cliCtx)
			cmd.SetContext(ctx)

			err := cmd.Execute()

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, buf.String())
			}
		})
	}
}
