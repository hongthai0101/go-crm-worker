package clients

import (
	"context"
	"crm-worker-go/config"
	"crm-worker-go/types"
	"fmt"
	"net/http"
)

type MasterDataClient struct {
	client *Client
}

func NewMasterDataClient(token string) *MasterDataClient {
	return &MasterDataClient{
		client: NewClient(config.GetConfig().ServiceConfig.MasterDataUrl, token),
	}
}

func (c *MasterDataClient) findAssetTypes(ctx context.Context) ([]*types.IMasterData, error) {
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

func (c *MasterDataClient) GetAssetType(ctx context.Context) *map[string]string {
	res, _ := c.findAssetTypes(ctx)
	result := make(map[string]string, len(res))
	for i := 0; i < len(res); i++ {
		result[res[i].Code] = res[i].Name
	}
	return &result
}

func (c *MasterDataClient) findSource(ctx context.Context) ([]*types.IMasterData, error) {
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

func (c *MasterDataClient) GetSource(ctx context.Context) *map[string]string {
	res, _ := c.findSource(ctx)
	result := make(map[string]string, len(res))
	for i := 0; i < len(res); i++ {
		result[res[i].Code] = res[i].Name
	}
	return &result
}

func (c *MasterDataClient) findType(ctx context.Context) ([]*types.IMasterData, error) {
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

func (c *MasterDataClient) GetTypes(ctx context.Context) *map[string]string {
	res, _ := c.findType(ctx)
	result := make(map[string]string, len(res))
	for i := 0; i < len(res); i++ {
		result[res[i].Code] = res[i].Name
	}
	return &result
}

func (c *MasterDataClient) findStores(ctx context.Context) ([]*types.IMasterData, error) {
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

func (c *MasterDataClient) GetStores(ctx context.Context) *map[string]string {
	res, _ := c.findStores(ctx)
	result := make(map[string]string, len(res))
	for i := 0; i < len(res); i++ {
		result[res[i].Code] = res[i].Name
	}
	return &result
}

func (c *MasterDataClient) findStatuses(ctx context.Context) ([]*types.IMasterData, error) {
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

func (c *MasterDataClient) GetStatuses(ctx context.Context) *map[string]string {
	res, _ := c.findType(ctx)
	result := make(map[string]string, len(res))
	for i := 0; i < len(res); i++ {
		result[res[i].Code] = res[i].Name
	}
	return &result
}

func (c *MasterDataClient) findProvinces(ctx context.Context) ([]*types.IMasterData, error) {
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

func (c *MasterDataClient) GetProvinces(ctx context.Context) *map[string]string {
	res, _ := c.findProvinces(ctx)
	result := make(map[string]string, len(res))
	for i := 0; i < len(res); i++ {
		result[res[i].Code] = res[i].Name
	}
	return &result
}

func (c *MasterDataClient) findGroups(ctx context.Context) ([]*types.IMasterData, error) {
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

func (c *MasterDataClient) GetGroups(ctx context.Context) *map[string]string {
	res, _ := c.findGroups(ctx)
	result := make(map[string]string, len(res))
	for i := 0; i < len(res); i++ {
		result[res[i].Code] = res[i].Name
	}
	return &result
}
