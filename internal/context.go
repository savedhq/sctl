package internal

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	saved "github.com/savedhq/sdk-go"
)

type contextKey string

const (
	configKey contextKey = "config"
	clientKey contextKey = "client"
)

type CLIContext struct {
	Config     *Config
	Client     *saved.APIClient
	APICtx     context.Context
	JSONOutput bool
	Err        error
}

func NewCLIContext(cfg *Config) (*CLIContext, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	client, apiCtx := cfg.GetClient()

	return &CLIContext{
		Config: cfg,
		Client: client,
		APICtx: apiCtx,
	}, nil
}

func (c *CLIContext) GetWorkspaceID() string {
	return c.Config.WorkspaceID
}

func WithCLIContext(ctx context.Context, cliCtx *CLIContext) context.Context {
	return context.WithValue(ctx, configKey, cliCtx)
}

func GetCLIContext(ctx context.Context) *CLIContext {
	if cliCtx, ok := ctx.Value(configKey).(*CLIContext); ok {
		return cliCtx
	}
	return nil
}

func (c *CLIContext) ResolveWorkspaceID(flagValue string) (string, error) {
	// 1. Determine the identifier (Flag > Config)
	nameOrID := flagValue
	if nameOrID == "" {
		nameOrID = c.Config.WorkspaceID
	}
	if nameOrID == "" {
		return "", fmt.Errorf("workspace_id required (use --workspace or set in config)")
	}

	// 2. If it's a UUID, return it directly
	if _, err := uuid.Parse(nameOrID); err == nil {
		return nameOrID, nil
	}

	// 3. Try to find by name
	resp, r, err := c.Client.WorkspacesAPI.ListWorkspaces(c.APICtx).Execute()
	if err != nil {
		return "", fmt.Errorf("failed to list workspaces: %w", err)
	}
	defer r.Body.Close()

	for _, ws := range resp {
		if ws.GetName() == nameOrID {
			return ws.GetId(), nil
		}
	}

	return "", fmt.Errorf("workspace '%s' not found", nameOrID)
}

func (c *CLIContext) ResolveJobID(workspaceID, nameOrID string) (string, error) {
	if _, err := uuid.Parse(nameOrID); err == nil {
		return nameOrID, nil
	}

	// Try to find by name
	resp, r, err := c.Client.JobsAPI.ListJobs(c.APICtx, workspaceID).Execute()
	if err != nil {
		return "", fmt.Errorf("failed to list jobs: %w", err)
	}
	defer r.Body.Close()

	for _, job := range resp {
		if job.GetName() == nameOrID {
			return job.GetId(), nil
		}
	}

	return "", fmt.Errorf("job '%s' not found", nameOrID)
}

func (c *CLIContext) ResolveAgentID(workspaceID, nameOrID string) (string, error) {
	if _, err := uuid.Parse(nameOrID); err == nil {
		return nameOrID, nil
	}

	// Try to find by name
	resp, r, err := c.Client.AgentsAPI.ListAgents(c.APICtx, workspaceID).Execute()
	if err != nil {
		return "", fmt.Errorf("failed to list agents: %w", err)
	}
	defer r.Body.Close()

	for _, agent := range resp {
		if agent.GetName() == nameOrID {
			return agent.GetId(), nil
		}
	}

	return "", fmt.Errorf("agent '%s' not found", nameOrID)
}
