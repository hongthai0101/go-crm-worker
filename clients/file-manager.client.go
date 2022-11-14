package clients

import (
	"bytes"
	"context"
	"crm-worker-go/config"
	"crm-worker-go/types"
	"encoding/json"
	"fmt"
	"net/http"
)

type FileManagerClient interface {
	FindExportRequests(ctx context.Context, id string) (*types.IExportRequest, error)
	UpdateExportRequestFailure(ctx context.Context, id string) error
	CreateFile(ctx context.Context, exportRequestId string, url string, info interface{}) error
}

type fileManagerClient struct {
	client *Client
}

func NewFileManagerClient(token string) FileManagerClient {
	return &fileManagerClient{
		client: NewClient(config.ServiceConfig["fileManagerUrl"], token),
	}
}

func (c *fileManagerClient) FindExportRequests(ctx context.Context, id string) (*types.IExportRequest, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/export-requests/%s", c.client.baseURL, id), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	var res types.IExportRequest
	if err := c.client.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *fileManagerClient) UpdateExportRequestFailure(ctx context.Context, id string) error {
	update, _ := json.Marshal(map[string]string{"status": "failure"})

	req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/export-requests/%s", c.client.baseURL, id), bytes.NewBuffer(update))
	if err != nil {
		return err
	}

	req = req.WithContext(ctx)
	var res types.IExportRequest
	if err := c.client.sendRequest(req, &res); err != nil {
		return err
	}
	return nil
}

func (c *fileManagerClient) CreateFile(
	ctx context.Context,
	exportRequestId string,
	url string,
	info interface{},
) error {
	body, _ := json.Marshal(types.CreateFileRequest{
		Url:             url,
		Info:            info,
		ExportRequestId: exportRequestId,
	})

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/files", c.client.baseURL), bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req = req.WithContext(ctx)
	var res types.IExportRequest
	if err := c.client.sendRequest(req, &res); err != nil {
		return err
	}
	return nil
}
