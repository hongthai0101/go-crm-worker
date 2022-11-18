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

type EmployeeClient struct {
	Client *Client
}

func NewEmployeeClient() *EmployeeClient {
	return &EmployeeClient{
		Client: NewClient(config.GetConfig().ServiceConfig.EmployeeUrl),
	}
}

func (c *EmployeeClient) findByIds(ctx context.Context, ids []string) ([]*types.IEmployee, error) {
	postBody, _ := json.Marshal(map[string][]string{
		"ids": ids,
	})

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/employees/list", c.Client.baseURL), bytes.NewBuffer(postBody))
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	var employees []*types.IEmployee
	if err := c.Client.sendRequest(req, &employees); err != nil {
		return nil, err
	}
	return employees, nil
}

func (c *EmployeeClient) GetEmployees(ctx context.Context, ids []string) (*map[string]string, error) {
	employees, _ := c.findByIds(ctx, ids)
	result := make(map[string]string, len(ids))
	for i := 0; i < len(employees); i++ {
		result[employees[i].ID] = employees[i].DisplayName
	}
	return &result, nil
}
