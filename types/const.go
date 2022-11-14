package types

const (
	Pending                    = "pending"
	InProgress                 = "inProgress"
	Completed                  = "completed"
	Failure                    = "failure"
	DDMMYYYY                   = "02-01-2006"
	YYMMDD                     = "060102"
	GroupOld                   = "OLD"
	GroupNew                   = "NEW"
	SaleOppStatusNew           = "NEW"
	SaleOppStatusSuccess       = "SUCCESS"
	SaleOppStatusPending       = "PENDING"
	SaleOppStatusConsulting    = "CONSULTING"
	SaleOppStatusDealt         = "DEALT"
	SaleOppStatusDenied        = "DENIED"
	SaleOppStatusCancel        = "CANCEL"
	SaleOppStatusUnContactable = "UNCONTACTABLE"
	SaleOppTypeBorrower        = "BORROWER"
	SaleOppTypePartner         = "PARTNER"
	SaleOppTypeInvestment      = "INVESTMENT"
)

const (
	TopicSubscriptionTypeOrderCreated = "customer.order.created"
	TopicSubscriptionTypeOrderUpdated = "customer.order.updated"
	TopicSubscriptionTypeLoanUpdated  = "loan.order.updated"
)
