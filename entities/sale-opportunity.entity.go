package entities

import (
	"crm-worker-go/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const CollectionSaleOpportunities = "SaleOpportunities"

type AssetMedia struct {
	Url      string `bson:"url" json:"url,omitempty"`
	MimeType string `bson:"mimeType" json:"mimeType,omitempty"`
}

type Asset struct {
	Description string       `bson:"description" json:"description,omitempty"`
	Media       []AssetMedia `bson:"media" json:"media,omitempty"`
	AssetType   string       `bson:"assetType" json:"assetType,omitempty"`
	DemandLoan  interface{}  `bson:"demandLoan" json:"demandLoan,omitempty"`
	LoanTerm    interface{}  `bson:"loanTerm" json:"loanTerm,omitempty"`
}

type SourceRefs struct {
	Source     string      `bson:"source" json:"source,omitempty"`
	RefId      string      `bson:"ref_id" json:"refId,omitempty"`
	CustomerId interface{} `bson:"customerId" json:"customerId,omitempty"`
}

type SaleOpportunity struct {
	ID              primitive.ObjectID     `bson:"_id,omitempty" json:"id,omitempty"`
	SourceRefs      SourceRefs             `bson:"source_refs" json:"source_refs"`
	Code            string                 `bson:"code" json:"code,omitempty"`
	Status          types.SaleOppStatus    `bson:"status" json:"status,omitempty"`
	Type            types.SaleOppType      `bson:"type" json:"type,omitempty"`
	Source          string                 `bson:"source" json:"source,omitempty"`
	Group           types.SaleOppGroup     `bson:"group" json:"group,omitempty"`
	Assets          Asset                  `bson:"assets" json:"assets"`
	EmployeeBy      string                 `bson:"employeeBy" json:"employeeBy,omitempty"`
	StoreCode       string                 `bson:"storeCode" json:"storeCode,omitempty"`
	DisbursedAt     *time.Time             `bson:"disbursedAt" json:"disbursedAt,omitempty"`
	ContractCode    string                 `bson:"contractCode" json:"contractCode,omitempty"`
	Tags            []string               `bson:"tags" json:"tags,omitempty"`
	DisbursedAmount int                    `bson:"disbursedAmount" json:"disbursedAmount,omitempty"`
	LeadId          primitive.ObjectID     `bson:"leadId" json:"leadId,omitempty"`
	Hash            string                 `bson:"hash" json:"hash,omitempty"`
	Metadata        map[string]interface{} `bson:"metadata" json:"metadata,omitempty"`
	BaseEntity      `bson:"inline"`
}
