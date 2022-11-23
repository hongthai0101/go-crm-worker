package clients

import (
	"crm-worker-go/config"
	"crm-worker-go/types"
	"crm-worker-go/utils"
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

func (c *MasterDataClient) findAssetTypes() ([]*types.IMasterData, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/asset-types", c.Client.baseURL), nil)
	if err != nil {
		utils.Logger.Error(err)
		return nil, err
	}

	req = req.WithContext(HttpCtx)
	var res []*types.IMasterData
	if err = c.Client.sendRequest(req, &res); err != nil {
		utils.Logger.Error(err)
		return nil, err
	}
	return res, nil
}

func (c *MasterDataClient) GetAssetType() *map[string]string {
	res, _ := c.findAssetTypes()
	result := make(map[string]string, len(res))
	for i := 0; i < len(res); i++ {
		result[res[i].Code] = res[i].Name
	}
	return &result
}

func (c *MasterDataClient) findSource() ([]*types.IMasterData, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/crm/sources", c.Client.baseURL), nil)
	if err != nil {
		utils.Logger.Error(err)
		return nil, err
	}

	req = req.WithContext(HttpCtx)
	var res []*types.IMasterData
	if err = c.Client.sendRequest(req, &res); err != nil {
		utils.Logger.Error(err)
		return nil, err
	}
	return res, nil
}

func (c *MasterDataClient) GetSource() *map[string]string {
	res, _ := c.findSource()
	result := make(map[string]string, len(res))
	for i := 0; i < len(res); i++ {
		result[res[i].Code] = res[i].Name
	}
	return &result
}

func (c *MasterDataClient) findType() ([]*types.IMasterData, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/crm/types", c.Client.baseURL), nil)
	if err != nil {
		utils.Logger.Error(err)
		return nil, err
	}

	req = req.WithContext(HttpCtx)
	var res []*types.IMasterData
	if err = c.Client.sendRequest(req, &res); err != nil {
		utils.Logger.Error(err)
		return nil, err
	}
	return res, nil
}

func (c *MasterDataClient) GetTypes() *map[string]string {
	res, _ := c.findType()
	result := make(map[string]string, len(res))
	for i := 0; i < len(res); i++ {
		result[res[i].Code] = res[i].Name
	}
	return &result
}

func (c *MasterDataClient) findStores() ([]*types.IMasterData, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/stores", c.Client.baseURL), nil)
	if err != nil {
		utils.Logger.Error(err)
		return nil, err
	}

	req = req.WithContext(HttpCtx)
	var res []*types.IMasterData
	if err = c.Client.sendRequest(req, &res); err != nil {
		utils.Logger.Error(err)
		return nil, err
	}
	return res, nil
}

func (c *MasterDataClient) GetStores() *map[string]string {
	res, _ := c.findStores()
	result := make(map[string]string, len(res))
	for i := 0; i < len(res); i++ {
		result[res[i].Code] = res[i].Name
	}
	return &result
}

func (c *MasterDataClient) findStatuses() ([]*types.IMasterData, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/crm/statuses", c.Client.baseURL), nil)
	if err != nil {
		utils.Logger.Error(err)
		return nil, err
	}

	req = req.WithContext(HttpCtx)
	var res []*types.IMasterData
	if err = c.Client.sendRequest(req, &res); err != nil {
		utils.Logger.Error(err)
		return nil, err
	}
	return res, nil
}

func (c *MasterDataClient) GetStatuses() *map[string]string {
	res, _ := c.findType()
	result := make(map[string]string, len(res))
	for i := 0; i < len(res); i++ {
		result[res[i].Code] = res[i].Name
	}
	return &result
}

func (c *MasterDataClient) findProvinces() ([]*types.IMasterData, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/provinces", c.Client.baseURL), nil)
	if err != nil {
		utils.Logger.Error(err)
		return nil, err
	}

	req = req.WithContext(HttpCtx)
	var res []*types.IMasterData
	if err = c.Client.sendRequest(req, &res); err != nil {
		utils.Logger.Error(err)
		return nil, err
	}
	return res, nil
}

func (c *MasterDataClient) GetProvinces() *map[string]string {
	res, _ := c.findProvinces()
	result := make(map[string]string, len(res))
	for i := 0; i < len(res); i++ {
		result[res[i].Code] = res[i].Name
	}
	return &result
}

func (c *MasterDataClient) findGroups() ([]*types.IMasterData, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/crm/groups", c.Client.baseURL), nil)
	if err != nil {
		utils.Logger.Error(err)
		return nil, err
	}

	req = req.WithContext(HttpCtx)
	var res []*types.IMasterData
	if err = c.Client.sendRequest(req, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *MasterDataClient) GetGroups() *map[string]string {
	res, _ := c.findGroups()
	result := make(map[string]string, len(res))
	for i := 0; i < len(res); i++ {
		result[res[i].Code] = res[i].Name
	}
	return &result
}
