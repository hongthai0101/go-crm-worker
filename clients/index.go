package clients

type HttpClient struct {
	EmployeeClient    *EmployeeClient
	FileManagerClient *FileManagerClient
	MasterDataClient  *MasterDataClient
}

func NewHttpClient() *HttpClient {
	return &HttpClient{
		EmployeeClient:    NewEmployeeClient(),
		FileManagerClient: NewFileManagerClient(),
		MasterDataClient:  NewMasterDataClient(),
	}
}
