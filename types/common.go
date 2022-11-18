package types

import "time"

type IMasterData struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type IEmployee struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
}

type ExportRequestData struct {
	Column       string `json:"column"`
	Start        string `json:"start"`
	End          string `json:"end"`
	StatusExport string `json:"statusExport"`
	SourceExport string `json:"sourceExport"`
	StoreCodes   string `json:"storeCodes"`
}

type IExportRequest struct {
	Type       string            `json:"type"`
	Resource   string            `json:"resource"`
	StoreCodes string            `json:"storeCodes"`
	Data       ExportRequestData `json:"data"`
	Status     string            `json:"status"`
	Column     string            `json:"column"`
	CreatedBy  string            `json:"createdBy"`
	CreatedAt  time.Time         `json:"createdAt"`
}

type CreateFileRequest struct {
	Url             string      `json:"url" :"url"`
	Info            interface{} `json:"info" :"info"`
	ExportRequestId string      `json:"exportRequestId" :"exportRequestId"`
}

type PayloadMessageExport struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

type RequestOrder struct {
	CustomerName string `json:"customer_name"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	AssetType    string `json:"assetType"`
	Detail       string `json:"detail"`
	Days         string `json:"days"`
	Bill         uint16 `json:"bill"`
	Code         string `json:"code"`
	Id           string `json:"id"`
	CreatedBy    string `json:"createdBy"`
	CustomerId   string `json:"customerId"`
}

type RequestMessageOrder struct {
	Order    RequestOrder  `json:"order"`
	Metadata interface{}   `json:"metadata"`
	Images   []interface{} `json:"images"`
}

type PayloadBorrowDisbursed struct {
	LoanAmount     int    `json:"loanAmount"`
	ModifiedAmount int    `json:"modifiedAmount"`
	ContractCode   string `json:"contractCode"`
}

type PayloadMessageDisbursed struct {
	SaleOppCode     string                  `json:"saleOppCode"`
	LoanPackageCode string                  `json:"loanPackageCode"`
	Lead            MessageDisbursedLead    `json:"lead"`
	SaleOpp         MessageDisbursedSaleOpp `json:"saleOpp"`
	IsManual        bool                    `json:"isManual"`
	LoanAmount      int                     `json:"loanAmount"`
	ModifiedAmount  int                     `json:"modifiedAmount"`
	ContractCode    string                  `json:"contractCode"`
}

type MessageDisbursedLead struct {
	FullName     string `json:"fullName"`
	Phone        string `json:"phone"`
	NationalId   string `json:"nationalId"`
	Account      string `json:"account"`
	AccountStore string `json:"accountStore"`
	CustomerId   string `json:"customerId"`
	Source       string `json:"source"`
}

type MessageDisbursedSaleOpp struct {
	Description     string    `json:"description"`
	ContractCode    string    `json:"contractCode"`
	AssetType       string    `json:"assetType"`
	LoanTerm        string    `json:"loanTerm"`
	DisbursedAmount int       `json:"disbursedAmount"`
	AccountStore    string    `json:"accountStore"`
	CreatedId       string    `json:"createdId"`
	Account         string    `json:"account"`
	DemandLoan      int       `json:"demandLoan"`
	DisbursedAt     time.Time `json:"disbursedAt"`
}
