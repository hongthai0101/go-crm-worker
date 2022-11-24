package types

const (
	DDMMYYYY string = "02-01-2006"
	YYMMDD   string = "060102"
	YYMM     string = "0601"
)

type ExportRequestStatus string

const (
	Pending    ExportRequestStatus = "pending"
	InProgress ExportRequestStatus = "inProgress"
	Completed  ExportRequestStatus = "completed"
	Failure    ExportRequestStatus = "failure"
)

type SaleOppGroup string

const (
	GroupOld SaleOppGroup = "OLD"
	GroupNew SaleOppGroup = "NEW"
)

type SaleOppStatus string

const (
	SaleOppStatusNew           SaleOppStatus = "NEW"
	SaleOppStatusSuccess       SaleOppStatus = "SUCCESS"
	SaleOppStatusPending       SaleOppStatus = "PENDING"
	SaleOppStatusConsulting    SaleOppStatus = "CONSULTING"
	SaleOppStatusDealt         SaleOppStatus = "DEALT"
	SaleOppStatusDenied        SaleOppStatus = "DENIED"
	SaleOppStatusCancel        SaleOppStatus = "CANCEL"
	SaleOppStatusUnContactable SaleOppStatus = "UNCONTACTABLE"
)

type SaleOppType string

const (
	SaleOppTypeBorrower   SaleOppType = "BORROWER"
	SaleOppTypePartner    SaleOppType = "PARTNER"
	SaleOppTypeInvestment SaleOppType = "INVESTMENT"
)

type TopicSubscriptionType string

const (
	TopicSubscriptionTypeOrderCreated TopicSubscriptionType = "customer.order.created"
	TopicSubscriptionTypeOrderUpdated TopicSubscriptionType = "customer.order.updated"
	TopicSubscriptionTypeLoanUpdated  TopicSubscriptionType = "loan.order.updated"
)

type ExportRequestType string

const (
	ExportRequestTypeSaleOpp ExportRequestType = "saleOpportunity"
	ExportRequestTypeLead    ExportRequestType = "lead"
)

type PolicyResource string

const (
	PolicyResourceSaleOpportunities PolicyResource = "SALES_OPPORTUNITY"
	PolicyResourceLead              PolicyResource = "LEAD"
)

type AuthorizationAction string

const (
	AuthorizationActionReadAny = "read:any"
)
