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

type FileManagerClient struct {
	Client *Client
}

func NewFileManagerClient() *FileManagerClient {
	return &FileManagerClient{
		Client: NewClient(config.GetConfig().ServiceConfig.FileManagerUrl),
	}
}

func (c *FileManagerClient) FindExportRequests(ctx context.Context, id string) (*types.IExportRequest, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/export-requests/%s", c.Client.baseURL, id), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	var res types.IExportRequest
	if err := c.Client.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *FileManagerClient) UpdateExportRequestFailure(ctx context.Context, id string) error {
	update, _ := json.Marshal(map[string]string{"status": "failure"})

	req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/export-requests/%s", c.Client.baseURL, id), bytes.NewBuffer(update))
	if err != nil {
		return err
	}

	req = req.WithContext(ctx)
	var res types.IExportRequest
	if err := c.Client.sendRequest(req, &res); err != nil {
		return err
	}
	return nil
}

func (c *FileManagerClient) CreateFile(
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

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/files", c.Client.baseURL), bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req = req.WithContext(ctx)
	res := struct{}{}
	if err = c.Client.sendRequest(req, &res); err != nil {
		return err
	}
	return nil
}
