package clients

import (
	"context"
	"crm-worker-go/config"
	"crm-worker-go/types"
	"fmt"
	"net/http"
)

type MasterDataClient interface {
	GetAssetType(ctx context.Context) *map[string]string
	GetSource(ctx context.Context) *map[string]string
	GetTypes(ctx context.Context) *map[string]string
	GetStores(ctx context.Context) *map[string]string
	GetStatuses(ctx context.Context) *map[string]string
	GetProvinces(ctx context.Context) *map[string]string
	GetGroups(ctx context.Context) *map[string]string
}

type masterDataClient struct {
	client *Client
}

func NewMasterDataClient(token string) MasterDataClient {
	return &masterDataClient{
		client: NewClient(config.ServiceConfig["masterDataUrl"], token),
	}
}

func (c *masterDataClient) findAssetTypes(ctx context.Context) ([]*types.IMasterData, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/asset-types", c.client.baseURL), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	var res []*types.IMasterData
	if err := c.client.sendRequest(req, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *masterDataClient) GetAssetType(ctx context.Context) *map[string]string {
	res, _ := c.findAssetTypes(ctx)
	result := make(map[string]string, len(res))
	for i := 0; i < len(res); i++ {
		result[res[i].Code] = res[i].Name
	}
	return &result
}

func (c *masterDataClient) findSource(ctx context.Context) ([]*types.IMasterData, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/crm/sources", c.client.baseURL), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	var res []*types.IMasterData
	if err := c.client.sendRequest(req, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *masterDataClient) GetSource(ctx context.Context) *map[string]string {
	res, _ := c.findSource(ctx)
	result := make(map[string]string, len(res))
	for i := 0; i < len(res); i++ {
		result[res[i].Code] = res[i].Name
	}
	return &result
}

func (c *masterDataClient) findType(ctx context.Context) ([]*types.IMasterData, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/crm/types", c.client.baseURL), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	var res []*types.IMasterData
	if err := c.client.sendRequest(req, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *masterDataClient) GetTypes(ctx context.Context) *map[string]string {
	res, _ := c.findType(ctx)
	result := make(map[string]string, len(res))
	for i := 0; i < len(res); i++ {
		result[res[i].Code] = res[i].Name
	}
	return &result
}

func (c *masterDataClient) findStores(ctx context.Context) ([]*types.IMasterData, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/stores", c.client.baseURL), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	var res []*types.IMasterData
	if err := c.client.sendRequest(req, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *masterDataClient) GetStores(ctx context.Context) *map[string]string {
	res, _ := c.findStores(ctx)
	result := make(map[string]string, len(res))
	for i := 0; i < len(res); i++ {
		result[res[i].Code] = res[i].Name
	}
	return &result
}

func (c *masterDataClient) findStatuses(ctx context.Context) ([]*types.IMasterData, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/crm/statuses", c.client.baseURL), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	var res []*types.IMasterData
	if err := c.client.sendRequest(req, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *masterDataClient) GetStatuses(ctx context.Context) *map[string]string {
	res, _ := c.findType(ctx)
	result := make(map[string]string, len(res))
	for i := 0; i < len(res); i++ {
		result[res[i].Code] = res[i].Name
	}
	return &result
}

func (c *masterDataClient) findProvinces(ctx context.Context) ([]*types.IMasterData, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/provinces", c.client.baseURL), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	var res []*types.IMasterData
	if err := c.client.sendRequest(req, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *masterDataClient) GetProvinces(ctx context.Context) *map[string]string {
	res, _ := c.findProvinces(ctx)
	result := make(map[string]string, len(res))
	for i := 0; i < len(res); i++ {
		result[res[i].Code] = res[i].Name
	}
	return &result
}

func (c *masterDataClient) findGroups(ctx context.Context) ([]*types.IMasterData, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/crm/groups", c.client.baseURL), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	var res []*types.IMasterData
	if err := c.client.sendRequest(req, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *masterDataClient) GetGroups(ctx context.Context) *map[string]string {
	res, _ := c.findGroups(ctx)
	result := make(map[string]string, len(res))
	for i := 0; i < len(res); i++ {
		result[res[i].Code] = res[i].Name
	}
	return &result
}
