package clients

import (
	"bytes"
	"crm-worker-go/config"
	"crm-worker-go/types"
	"crm-worker-go/utils"
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

func (c *FileManagerClient) FindExportRequests(id string) (*types.IExportRequest, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/export-requests/%s", c.Client.baseURL, id), nil)
	if err != nil {
		utils.Logger.Error(err)
		return nil, err
	}

	req = req.WithContext(HttpCtx)
	var res types.IExportRequest
	if err = c.Client.sendRequest(req, &res); err != nil {
		utils.Logger.Error(err)
		return nil, err
	}

	return &res, nil
}

func (c *FileManagerClient) UpdateExportRequestFailure(id string) error {
	update, _ := json.Marshal(map[string]string{"status": "failure"})

	req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/export-requests/%s", c.Client.baseURL, id), bytes.NewBuffer(update))
	if err != nil {
		utils.Logger.Error(err)
		return err
	}

	req = req.WithContext(HttpCtx)
	var res types.IExportRequest
	if err = c.Client.sendRequest(req, &res); err != nil {
		utils.Logger.Error(err)
		return err
	}
	return nil
}

func (c *FileManagerClient) CreateFile(
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
		utils.Logger.Error(err)
		return err
	}

	req = req.WithContext(HttpCtx)
	res := struct{}{}
	if err = c.Client.sendRequest(req, &res); err != nil {
		utils.Logger.Error(err)
		return err
	}
	return nil
}
