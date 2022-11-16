package clients

type HttpClient struct {
	EmployeeClient    *EmployeeClient
	FileManagerClient *FileManagerClient
	MasterDataClient  *MasterDataClient
}

func NewHttpClient(string2 string) *HttpClient {
	return &HttpClient{
		EmployeeClient:    NewEmployeeClient(string2),
		FileManagerClient: NewFileManagerClient(string2),
		MasterDataClient:  NewMasterDataClient(string2),
	}
}
