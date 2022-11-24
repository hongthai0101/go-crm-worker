package entities

import (
	"crm-worker-go/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const CollectionLead = "Lead"

type Lead struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"ID,omitempty"`
	CustomerId string             `bson:"customerId" json:"customerId,omitempty"`
	FullName   string             `bson:"fullName" json:"fullName,omitempty"`
	Phone      string             `bson:"phone" json:"phone,omitempty"`
	Email      string             `bson:"email" json:"email,omitempty"`
	NationalId string             `bson:"nationalId" json:"nationalId,omitempty"`
	PassportId string             `bson:"passportId" json:"passportId,omitempty"`
	TaxId      string             `bson:"taxId" json:"taxId,omitempty"`
	Address    string             `bson:"address" json:"address,omitempty"`
	Province   string             `bson:"province" json:"province,omitempty"`
	District   string             `bson:"district" json:"district,omitempty"`
	Source     string             `bson:"source" json:"source,omitempty"`
	EmployeeBy string             `bson:"employeeBy" json:"employeeBy,omitempty"`
	StoreCode  string             `bson:"storeCode" json:"storeCode,omitempty"`
	Type       types.SaleOppType  `bson:"type" json:"type,omitempty"`
	Birthday   string             `bson:"birthday" json:"birthday,omitempty"`
	Gender     string             `bson:"gender" json:"gender,omitempty"`
	BaseEntity `bson:"inline"`
}
