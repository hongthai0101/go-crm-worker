package clients

import (
	"context"
	"crm-worker-go/config"
	"crm-worker-go/types"
	"fmt"
	"net/http"
)

type MasterDataClient struct {
	Client *Client
}

func NewMasterDataClient() *MasterDataClient {
	return &MasterDataClient{
		Client: NewClient(config.GetConfig().ServiceConfig.MasterDataUrl),
	}
}

func (c *MasterDataClient) findAssetTypes(ctx context.Context) ([]*types.IMasterData, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/asset-types", c.Client.baseURL), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	var res []*types.IMasterData
	if err := c.Client.sendRequest(req, &res); err != nil {
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
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/crm/sources", c.Client.baseURL), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	var res []*types.IMasterData
	if err := c.Client.sendRequest(req, &res); err != nil {
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
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/crm/types", c.Client.baseURL), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	var res []*types.IMasterData
	if err := c.Client.sendRequest(req, &res); err != nil {
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
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/stores", c.Client.baseURL), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	var res []*types.IMasterData
	if err := c.Client.sendRequest(req, &res); err != nil {
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
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/crm/statuses", c.Client.baseURL), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	var res []*types.IMasterData
	if err := c.Client.sendRequest(req, &res); err != nil {
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
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/provinces", c.Client.baseURL), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	var res []*types.IMasterData
	if err := c.Client.sendRequest(req, &res); err != nil {
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
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/crm/groups", c.Client.baseURL), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	var res []*types.IMasterData
	if err := c.Client.sendRequest(req, &res); err != nil {
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
