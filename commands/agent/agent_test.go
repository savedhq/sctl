package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/savedhq/sctl/internal"
	saved "github.com/savedhq/sdk-go"
	"github.com/stretchr/testify/assert"
)

type AgentCredentials struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func TestAgentCommands(t *testing.T) {
	var handler http.HandlerFunc

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if handler != nil {
			handler(w, r)
		}
	}))
	defer mockServer.Close()

	cliCtx := &internal.CLIContext{
		Config: &internal.Config{
			ServerURL: mockServer.URL,
			APIKey:    "test-token",
		},
	}
	client, apiCtx := cliCtx.Config.GetClient()
	cliCtx.Client = client
	cliCtx.APICtx = apiCtx

	execute := func(args ...string) (string, error) {
		cmd := NewAgentCmd()
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)

		cmd.SetArgs(args)

		ctx := internal.WithCLIContext(context.Background(), cliCtx)
		cmd.SetContext(ctx)

		err := cmd.Execute()
		return buf.String(), err
	}

	t.Run("create", func(t *testing.T) {
		agentName := "test-agent"
		agentID := uuid.New().String()
		workspaceID := uuid.New().String()
		cliCtx.Config.WorkspaceID = workspaceID

		mockAgent := saved.ListAgents200ResponseInner{}
		mockAgent.SetId(agentID)
		mockAgent.SetName(agentName)
		mockAgent.SetStatus("offline")
		mockAgent.SetCreatedAt(time.Now())
		mockAgent.SetUpdatedAt(time.Now())

		handler = func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "/v1/workspaces/"+workspaceID+"/agents", r.URL.Path)

			var reqBody saved.CreateAgentRequest
			err := json.NewDecoder(r.Body).Decode(&reqBody)
			assert.NoError(t, err)
			assert.Equal(t, agentName, reqBody.GetName())

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(mockAgent)
		}

		output, err := execute("create", "--name", agentName, "-w", workspaceID)

		assert.NoError(t, err)
		assert.Contains(t, output, agentID)
	})

	t.Run("list", func(t *testing.T) {
		workspaceID := uuid.New().String()
		cliCtx.Config.WorkspaceID = workspaceID

		mockAgents := []saved.ListAgents200ResponseInner{
			{Id: saved.PtrString(uuid.New().String()), Name: saved.PtrString("agent-1"), Status: saved.PtrString("offline")},
			{Id: saved.PtrString(uuid.New().String()), Name: saved.PtrString("agent-2"), Status: saved.PtrString("online")},
		}

		handler = func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, "/v1/workspaces/"+workspaceID+"/agents", r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(mockAgents)
		}

		output, err := execute("list", "-w", workspaceID)

		assert.NoError(t, err)
		assert.Contains(t, output, "agent-1")
		assert.Contains(t, output, "agent-2")
	})

	t.Run("get", func(t *testing.T) {
		agentID := uuid.New().String()
		workspaceID := uuid.New().String()
		cliCtx.Config.WorkspaceID = workspaceID

		mockAgent := saved.ListAgents200ResponseInner{}
		mockAgent.SetId(agentID)
		mockAgent.SetName("test-agent")
		mockAgent.SetStatus("offline")

		handler = func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, "/v1/workspaces/"+workspaceID+"/agents/"+agentID, r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(mockAgent)
		}

		output, err := execute("get", agentID, "-w", workspaceID)

		assert.NoError(t, err)
		assert.Contains(t, output, agentID)
		assert.Contains(t, output, "test-agent")
	})

	t.Run("update", func(t *testing.T) {
		agentID := uuid.New().String()
		workspaceID := uuid.New().String()
		cliCtx.Config.WorkspaceID = workspaceID
		updatedName := "updated-agent-name"

		handler = func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPatch, r.Method)
			assert.Equal(t, "/v1/workspaces/"+workspaceID+"/agents/"+agentID, r.URL.Path)

			var reqBody saved.UpdateAgentRequest
			err := json.NewDecoder(r.Body).Decode(&reqBody)
			assert.NoError(t, err)
			assert.Equal(t, updatedName, reqBody.GetName())

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(saved.ListAgents200ResponseInner{})
		}

		output, err := execute("update", agentID, "--name", updatedName, "-w", workspaceID)

		assert.NoError(t, err)
		assert.Contains(t, output, "Agent updated")
	})

	t.Run("delete", func(t *testing.T) {
		agentID := uuid.New().String()
		workspaceID := uuid.New().String()
		cliCtx.Config.WorkspaceID = workspaceID

		handler = func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodDelete, r.Method)
			assert.Equal(t, "/v1/workspaces/"+workspaceID+"/agents/"+agentID, r.URL.Path)
			w.WriteHeader(http.StatusNoContent)
		}

		output, err := execute("delete", agentID, "-w", workspaceID)

		assert.NoError(t, err)
		assert.Contains(t, output, "Agent deleted")
	})

	t.Run("reset-credentials", func(t *testing.T) {
		agentID := uuid.New().String()
		workspaceID := uuid.New().String()
		cliCtx.Config.WorkspaceID = workspaceID

		mockCredentials := map[string]interface{}{
			"client_id":     "test-client-id",
			"client_secret": "test-client-secret",
		}

		handler = func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "/v1/workspaces/"+workspaceID+"/agents/"+agentID+"/credentials", r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(mockCredentials)
		}

		output, err := execute("reset-credentials", agentID, "-w", workspaceID)

		assert.NoError(t, err)
		assert.Contains(t, output, "Credentials reset")

		var creds AgentCredentials
		// The output contains more than just the JSON, so we need to find the start of the JSON object.
		jsonStart := bytes.IndexByte([]byte(output), '{')
		err = json.Unmarshal([]byte(output[jsonStart:]), &creds)
		assert.NoError(t, err)

		assert.Equal(t, "test-client-id", creds.ClientID)
		assert.Equal(t, "test-client-secret", creds.ClientSecret)
	})
}
