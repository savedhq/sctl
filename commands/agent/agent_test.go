package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/savedhq/sctl/internal"
	saved "github.com/savedhq/sdk-go"
)

// mockAgentsAPIService provides a mock implementation of the AgentsAPI for testing.
type mockAgentsAPIService struct {
	saved.AgentsAPI // Embed the interface to get default behavior if needed.

	// Mock functions to override specific API calls.
	listAgents            func(ctx context.Context, workspaceId string) saved.ApiListAgentsRequest
	getAgent              func(ctx context.Context, workspaceId string, agentId string) saved.ApiGetAgentRequest
	createAgent           func(ctx context.Context, workspaceId string) saved.ApiCreateAgentRequest
	updateAgent           func(ctx context.Context, workspaceId string, agentId string) saved.ApiUpdateAgentRequest
	deleteAgent           func(ctx context.Context, workspaceId string, agentId string) saved.ApiDeleteAgentRequest
	resetAgentCredentials func(ctx context.Context, workspaceId string, agentId string) saved.ApiResetAgentCredentialsRequest
}

func (m *mockAgentsAPIService) ListAgents(ctx context.Context, workspaceId string) saved.ApiListAgentsRequest {
	return m.listAgents(ctx, workspaceId)
}
func (m *mockAgentsAPIService) GetAgent(ctx context.Context, workspaceId string, agentId string) saved.ApiGetAgentRequest {
	return m.getAgent(ctx, workspaceId, agentId)
}
func (m *mockAgentsAPIService) CreateAgent(ctx context.Context, workspaceId string) saved.ApiCreateAgentRequest {
	return m.createAgent(ctx, workspaceId)
}
func (m *mockAgentsAPIService) UpdateAgent(ctx context.Context, workspaceId string, agentId string) saved.ApiUpdateAgentRequest {
	return m.updateAgent(ctx, workspaceId, agentId)
}
func (m *mockAgentsAPIService) DeleteAgent(ctx context.Context, workspaceId string, agentId string) saved.ApiDeleteAgentRequest {
	return m.deleteAgent(ctx, workspaceId, agentId)
}
func (m *mockAgentsAPIService) ResetAgentCredentials(ctx context.Context, workspaceId string, agentId string) saved.ApiResetAgentCredentialsRequest {
	return m.resetAgentCredentials(ctx, workspaceId, agentId)
}


// Mock request objects
type mockListAgentsRequest struct {
	saved.ApiListAgentsRequest
	ExecuteFunc func() ([]saved.Agent, *http.Response, error)
}
func (m *mockListAgentsRequest) Execute() ([]saved.Agent, *http.Response, error) {
	return m.ExecuteFunc()
}

type mockGetAgentRequest struct {
	saved.ApiGetAgentRequest
	ExecuteFunc func() (*saved.Agent, *http.Response, error)
}
func (m *mockGetAgentRequest) Execute() (*saved.Agent, *http.Response, error) {
	return m.ExecuteFunc()
}

type mockCreateAgentRequest struct {
	saved.ApiCreateAgentRequest
	req         saved.CreateAgentRequest
	ExecuteFunc func(saved.CreateAgentRequest) (*saved.Agent, *http.Response, error)
}
func (m *mockCreateAgentRequest) CreateAgentRequest(req saved.CreateAgentRequest) saved.ApiCreateAgentRequest {
	m.req = req
	return m
}
func (m *mockCreateAgentRequest) Execute() (*saved.Agent, *http.Response, error) {
	return m.ExecuteFunc(m.req)
}

type mockUpdateAgentRequest struct {
	saved.ApiUpdateAgentRequest
	req         saved.UpdateAgentRequest
	ExecuteFunc func(saved.UpdateAgentRequest) (*saved.Agent, *http.Response, error)
}
func (m *mockUpdateAgentRequest) UpdateAgentRequest(req saved.UpdateAgentRequest) saved.ApiUpdateAgentRequest {
	m.req = req
	return m
}
func (m *mockUpdateAgentRequest) Execute() (*saved.Agent, *http.Response, error) {
	return m.ExecuteFunc(m.req)
}

type mockDeleteAgentRequest struct {
	saved.ApiDeleteAgentRequest
	ExecuteFunc func() (*http.Response, error)
}
func (m *mockDeleteAgentRequest) Execute() (*http.Response, error) {
	return m.ExecuteFunc()
}

type mockResetAgentCredentialsRequest struct {
	saved.ApiResetAgentCredentialsRequest
	ExecuteFunc func() (*saved.AgentCredentials, *http.Response, error)
}
func (m *mockResetAgentCredentialsRequest) Execute() (*saved.AgentCredentials, *http.Response, error) {
	return m.ExecuteFunc()
}

func TestAgentListCmd(t *testing.T) {
	// Mock data
	mockAgents := []saved.Agent{
		{Id: "agent-1", Name: "Agent 1", Status: "active"},
		{Id: "agent-2", Name: "Agent 2", Status: "inactive"},
	}

	// Mock API client
	mockClient := &saved.APIClient{
		AgentsAPI: &mockAgentsAPIService{
			listAgents: func(ctx context.Context, workspaceId string) saved.ApiListAgentsRequest {
				return &mockListAgentsRequest{
					ExecuteFunc: func() ([]saved.Agent, *http.Response, error) {
						return mockAgents, &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(""))}, nil
					},
				}
			},
		},
	}

	// Create a new CLI context with the mock client
	cliCtx := &internal.CLIContext{
		Client: mockClient,
		APICtx: context.Background(),
	}

	// Create a new command and set the context
	cmd := newAgentListCmd()
	ctx := internal.WithCLIContext(context.Background(), cliCtx)
	cmd.SetContext(ctx)

	// Test case 1: Standard output
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check the output
	expectedOutput := "ID: agent-1\n  Name: Agent 1\n  Status: active\n\nID: agent-2\n  Name: Agent 2\n  Status: inactive\n\n"
	if !strings.Contains(buf.String(), expectedOutput) {
		t.Errorf("expected output to contain %q, but got %q", expectedOutput, buf.String())
	}

	// Test case 2: JSON output
	buf.Reset()
	cmd.SetArgs([]string{"--json"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check the JSON output
	var agents []saved.Agent
	if err := json.Unmarshal(buf.Bytes(), &agents); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}
	if len(agents) != 2 {
		t.Fatalf("expected 2 agents, but got %d", len(agents))
	}
	if agents[0].GetId() != "agent-1" || agents[1].GetId() != "agent-2" {
		t.Errorf("unexpected agent IDs in JSON output")
	}
}

func TestAgentGetCmd(t *testing.T) {
	// Mock data
	mockAgent := &saved.Agent{
		Id: "agent-1", Name: "Agent 1", Status: "active",
	}

	// Mock API client
	mockClient := &saved.APIClient{
		AgentsAPI: &mockAgentsAPIService{
			getAgent: func(ctx context.Context, workspaceId string, agentId string) saved.ApiGetAgentRequest {
				return &mockGetAgentRequest{
					ExecuteFunc: func() (*saved.Agent, *http.Response, error) {
						return mockAgent, &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(""))}, nil
					},
				}
			},
		},
	}

	// Create a new CLI context with the mock client
	cliCtx := &internal.CLIContext{
		Client: mockClient,
		APICtx: context.Background(),
	}

	// Create a new command and set the context
	cmd := newAgentGetCmd()
	ctx := internal.WithCLIContext(context.Background(), cliCtx)
	cmd.SetContext(ctx)

	// Test case 1: Standard output
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"agent-1"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check the output
	expectedOutput := "ID: agent-1\nName: Agent 1\nStatus: active\n"
	if !strings.Contains(buf.String(), expectedOutput) {
		t.Errorf("expected output to contain %q, but got %q", expectedOutput, buf.String())
	}

	// Test case 2: JSON output
	buf.Reset()
	cmd.SetArgs([]string{"agent-1", "--json"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check the JSON output
	var agent saved.Agent
	if err := json.Unmarshal(buf.Bytes(), &agent); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}
	if agent.GetId() != "agent-1" {
		t.Errorf("unexpected agent ID in JSON output")
	}
}

func TestAgentCreateCmd(t *testing.T) {
	// Mock data
	mockAgent := &saved.Agent{
		Id: "agent-1", Name: "New Agent", Status: "active",
	}

	// Mock API client
	mockClient := &saved.APIClient{
		AgentsAPI: &mockAgentsAPIService{
			createAgent: func(ctx context.Context, workspaceId string) saved.ApiCreateAgentRequest {
				return &mockCreateAgentRequest{
					ExecuteFunc: func(req saved.CreateAgentRequest) (*saved.Agent, *http.Response, error) {
						mockAgent.SetName(req.GetName())
						return mockAgent, &http.Response{StatusCode: 201, Body: io.NopCloser(strings.NewReader(""))}, nil
					},
				}
			},
		},
	}

	// Create a new CLI context with the mock client
	cliCtx := &internal.CLIContext{
		Client: mockClient,
		APICtx: context.Background(),
	}

	// Create a new command and set the context
	cmd := newAgentCreateCmd()
	ctx := internal.WithCLIContext(context.Background(), cliCtx)
	cmd.SetContext(ctx)

	// Test case 1: Standard output
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"--name", "New Agent"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check the output
	expectedOutput := "✓ Agent created\nID: agent-1\n"
	if !strings.Contains(buf.String(), expectedOutput) {
		t.Errorf("expected output to contain %q, but got %q", expectedOutput, buf.String())
	}

	// Test case 2: JSON output
	buf.Reset()
	cmd.SetArgs([]string{"--name", "New Agent", "--json"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check the JSON output
	var agent saved.Agent
	if err := json.Unmarshal(buf.Bytes(), &agent); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}
	if agent.GetId() != "agent-1" || agent.GetName() != "New Agent" {
		t.Errorf("unexpected agent data in JSON output")
	}
}

func TestAgentUpdateCmd(t *testing.T) {
	// Mock data
	mockAgent := &saved.Agent{
		Id: "agent-1", Name: "Updated Agent", Status: "active",
	}

	// Mock API client
	mockClient := &saved.APIClient{
		AgentsAPI: &mockAgentsAPIService{
			updateAgent: func(ctx context.Context, workspaceId string, agentId string) saved.ApiUpdateAgentRequest {
				return &mockUpdateAgentRequest{
					ExecuteFunc: func(req saved.UpdateAgentRequest) (*saved.Agent, *http.Response, error) {
						mockAgent.SetName(req.GetName())
						return mockAgent, &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(""))}, nil
					},
				}
			},
		},
	}

	// Create a new CLI context with the mock client
	cliCtx := &internal.CLIContext{
		Client: mockClient,
		APICtx: context.Background(),
	}

	// Create a new command and set the context
	cmd := newAgentUpdateCmd()
	ctx := internal.WithCLIContext(context.Background(), cliCtx)
	cmd.SetContext(ctx)

	// Test case 1: Standard output
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"agent-1", "--name", "Updated Agent"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check the output
	expectedOutput := "✓ Agent updated\n"
	if !strings.Contains(buf.String(), expectedOutput) {
		t.Errorf("expected output to contain %q, but got %q", expectedOutput, buf.String())
	}

	// Test case 2: JSON output
	buf.Reset()
	cmd.SetArgs([]string{"agent-1", "--name", "Updated Agent", "--json"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check the JSON output
	var agent saved.Agent
	if err := json.Unmarshal(buf.Bytes(), &agent); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}
	if agent.GetId() != "agent-1" || agent.GetName() != "Updated Agent" {
		t.Errorf("unexpected agent data in JSON output")
	}
}
