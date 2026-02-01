package backup

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/savedhq/sctl/internal"
	"github.com/savedhq/sdk-go/api"
	"github.com/stretchr/testify/assert"
)

type MockBackupsAPI struct {
	ListBackupsFunc  func(ctx context.Context, workspaceId string, jobId string) api.ApiListBackupsRequest
	GetBackupFunc    func(ctx context.Context, workspaceId string, jobId string, backupId string) api.ApiGetBackupRequest
	DeleteBackupFunc func(ctx context.Context, workspaceId string, jobId string, backupId string) api.ApiDeleteBackupRequest
	DownloadBackupFunc func(ctx context.Context, workspaceId string, jobId string, backupId string) api.ApiDownloadBackupRequest
}

func (m *MockBackupsAPI) ListBackups(ctx context.Context, workspaceId string, jobId string) api.ApiListBackupsRequest {
	return m.ListBackupsFunc(ctx, workspaceId, jobId)
}

func (m *MockBackupsAPI) GetBackup(ctx context.Context, workspaceId string, jobId string, backupId string) api.ApiGetBackupRequest {
	return m.GetBackupFunc(ctx, workspaceId, jobId, backupId)
}

func (m *MockBackupsAPI) DeleteBackup(ctx context.Context, workspaceId string, jobId string, backupId string) api.ApiDeleteBackupRequest {
	return m.DeleteBackupFunc(ctx, workspaceId, jobId, backupId)
}

func (m *MockBackupsAPI) DownloadBackup(ctx context.Context, workspaceId string, jobId string, backupId string) api.ApiDownloadBackupRequest {
	return m.DownloadBackupFunc(ctx, workspaceId, jobId, backupId)
}

type MockJobOperationsAPI struct {
    RequestBackupFunc func(ctx context.Context, workspaceId string, jobId string) api.ApiRequestBackupRequest
}

func (m *MockJobOperationsAPI) RequestBackup(ctx context.Context, workspaceId string, jobId string) api.ApiRequestBackupRequest {
    return m.RequestBackupFunc(ctx, workspaceId, jobId)
}

func TestBackupListCmd(t *testing.T) {
	cmd := newBackupListCmd()
	assert.NotNil(t, cmd)

	// Mock API client
	mockClient := &internal.APIClient{
		BackupsAPI: &MockBackupsAPI{
			ListBackupsFunc: func(ctx context.Context, workspaceId string, jobId string) api.ApiListBackupsRequest {
				return api.ApiListBackupsRequest{
					ApiService: nil,
					ctx:        ctx,
					workspaceId: workspaceId,
					jobId:      jobId,
					Execute: func(r api.ApiListBackupsRequest) ([]api.Backup, *http.Response, error) {
						return []api.Backup{
							{Id: api.PtrString("backup1"), Status: api.PtrString("completed"), CreatedAt: api.PtrTime(time.Now())},
						}, &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewBufferString(""))}, nil
					},
				}
			},
		},
	}

	ctx := internal.SetCLIContext(context.Background(), &internal.CLIContext{
		Client:      mockClient,
		WorkspaceID: "ws1",
	})
	cmd.SetContext(ctx)

	// Test human-readable output
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"job1"})
	err := cmd.Execute()
	assert.NoError(t, err)
	assert.Contains(t, out.String(), "ID: backup1")

	// Test JSON output
	out.Reset()
	cmd.SetArgs([]string{"job1", "--json"})
	err = cmd.Execute()
	assert.NoError(t, err)
	var backups []api.Backup
	err = json.Unmarshal(out.Bytes(), &backups)
	assert.NoError(t, err)
	assert.Len(t, backups, 1)
	assert.Equal(t, "backup1", *backups[0].Id)
}

func TestBackupGetCmd(t *testing.T) {
    cmd := newBackupGetCmd()
    assert.NotNil(t, cmd)

    // Mock API client
    mockClient := &internal.APIClient{
        BackupsAPI: &MockBackupsAPI{
            GetBackupFunc: func(ctx context.Context, workspaceId string, jobId string, backupId string) api.ApiGetBackupRequest {
                return api.ApiGetBackupRequest{
                    ApiService: nil,
                    ctx:        ctx,
                    workspaceId: workspaceId,
                    jobId:      jobId,
                    backupId:   backupId,
                    Execute: func(r api.ApiGetBackupRequest) (*api.Backup, *http.Response, error) {
                        return &api.Backup{
                            Id: api.PtrString("backup1"),
                            Status: api.PtrString("completed"),
                            CreatedAt: api.PtrTime(time.Now())},
                            &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewBufferString(""))}, nil
                    },
                }
            },
        },
    }

    ctx := internal.SetCLIContext(context.Background(), &internal.CLIContext{
        Client:      mockClient,
        WorkspaceID: "ws1",
    })
    cmd.SetContext(ctx)

    // Test human-readable output
    var out bytes.Buffer
    cmd.SetOut(&out)
    cmd.SetArgs([]string{"job1", "backup1"})
    err := cmd.Execute()
    assert.NoError(t, err)
    assert.Contains(t, out.String(), "ID: backup1")

    // Test JSON output
    out.Reset()
    cmd.SetArgs([]string{"job1", "backup1", "--json"})
    err = cmd.Execute()
    assert.NoError(t, err)
    var backup api.Backup
    err = json.Unmarshal(out.Bytes(), &backup)
    assert.NoError(t, err)
    assert.Equal(t, "backup1", *backup.Id)
}

func TestBackupDeleteCmd(t *testing.T) {
    cmd := newBackupDeleteCmd()
    assert.NotNil(t, cmd)

    // Mock API client
    mockClient := &internal.APIClient{
        BackupsAPI: &MockBackupsAPI{
            DeleteBackupFunc: func(ctx context.Context, workspaceId string, jobId string, backupId string) api.ApiDeleteBackupRequest {
                return api.ApiDeleteBackupRequest{
                    ApiService: nil,
                    ctx:        ctx,
                    workspaceId: workspaceId,
                    jobId:      jobId,
                    backupId:   backupId,
                    Execute: func(r api.ApiDeleteBackupRequest) (*http.Response, error) {
                        return &http.Response{StatusCode: 204, Body: ioutil.NopCloser(bytes.NewBufferString(""))}, nil
                    },
                }
            },
        },
    }

    ctx := internal.SetCLIContext(context.Background(), &internal.CLIContext{
        Client:      mockClient,
        WorkspaceID: "ws1",
    })
    cmd.SetContext(ctx)

    // Test human-readable output
    var out bytes.Buffer
    cmd.SetOut(&out)
    cmd.SetArgs([]string{"job1", "backup1"})
    err := cmd.Execute()
    assert.NoError(t, err)
    assert.Contains(t, out.String(), "Backup deleted")

    // Test JSON output
    out.Reset()
    cmd.SetArgs([]string{"job1", "backup1", "--json"})
    err = cmd.Execute()
    assert.NoError(t, err)
    var resp map[string]string
    err = json.Unmarshal(out.Bytes(), &resp)
    assert.NoError(t, err)
    assert.Equal(t, "deleted", resp["status"])
    assert.Equal(t, "backup1", resp["id"])
}

func TestBackupRequestCmd(t *testing.T) {
    cmd := newBackupRequestCmd()
    assert.NotNil(t, cmd)

    // Mock API client
    mockClient := &internal.APIClient{
        JobOperationsAPI: &MockJobOperationsAPI{
            RequestBackupFunc: func(ctx context.Context, workspaceId string, jobId string) api.ApiRequestBackupRequest {
                return api.ApiRequestBackupRequest{
                    ApiService: nil,
                    ctx:        ctx,
                    workspaceId: workspaceId,
                    jobId:      jobId,
                    Execute: func(r api.ApiRequestBackupRequest) (*api.Backup, *http.Response, error) {
                        return &api.Backup{
                            Id: api.PtrString("new-backup"),
                            Status: api.PtrString("pending"),
                            CreatedAt: api.PtrTime(time.Now())},
                            &http.Response{StatusCode: 202, Body: ioutil.NopCloser(bytes.NewBufferString(""))}, nil
                    },
                }
            },
        },
    }

    ctx := internal.SetCLIContext(context.Background(), &internal.CLIContext{
        Client:      mockClient,
        WorkspaceID: "ws1",
    })
    cmd.SetContext(ctx)

    // Test human-readable output
    var out bytes.Buffer
    cmd.SetOut(&out)
    cmd.SetArgs([]string{"job1"})
    err := cmd.Execute()
    assert.NoError(t, err)
    assert.Contains(t, out.String(), "Backup requested")
    assert.Contains(t, out.String(), "ID: new-backup")

    // Test JSON output
    out.Reset()
    cmd.SetArgs([]string{"job1", "--json"})
    err = cmd.Execute()
    assert.NoError(t, err)
    var backup api.Backup
    err = json.Unmarshal(out.Bytes(), &backup)
    assert.NoError(t, err)
    assert.Equal(t, "new-backup", *backup.Id)
}

func TestBackupDownloadCmd(t *testing.T) {
    cmd := newBackupDownloadCmd()
    assert.NotNil(t, cmd)

    // Mock API client
    mockClient := &internal.APIClient{
        BackupsAPI: &MockBackupsAPI{
            DownloadBackupFunc: func(ctx context.Context, workspaceId string, jobId string, backupId string) api.ApiDownloadBackupRequest {
                return api.ApiDownloadBackupRequest{
                    ApiService: nil,
                    ctx:        ctx,
                    workspaceId: workspaceId,
                    jobId:      jobId,
                    backupId:   backupId,
                    Execute: func(r api.ApiDownloadBackupRequest) (*api.BackupDownload, *http.Response, error) {
                        return &api.BackupDownload{Url: api.PtrString("https://example.com/backup.zip")}, &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewBufferString(""))}, nil
                    },
                }
            },
        },
    }

    ctx := internal.SetCLIContext(context.Background(), &internal.CLIContext{
        Client:      mockClient,
        WorkspaceID: "ws1",
    })
    cmd.SetContext(ctx)

    // Test human-readable output
    var out bytes.Buffer
    cmd.SetOut(&out)
    cmd.SetArgs([]string{"job1", "backup1"})
    err := cmd.Execute()
    assert.NoError(t, err)
    assert.Contains(t, out.String(), "Backup download URL generated")
    assert.Contains(t, out.String(), "https://example.com/backup.zip")

    // Test JSON output
    out.Reset()
    cmd.SetArgs([]string{"job1", "backup1", "--json"})
    err = cmd.Execute()
    assert.NoError(t, err)
    var download api.BackupDownload
    err = json.Unmarshal(out.Bytes(), &download)
    assert.NoError(t, err)
    assert.Equal(t, "https://example.com/backup.zip", *download.Url)
}
