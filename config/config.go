package config

type Config struct {
	DB *DBConfig
}

type DBConfig struct {
	Host  string
	Port  string
	User  string
	Pass  string
	Name  string
	Debug bool
}

func GetConfigDB() *DBConfig {
	return &DBConfig{
		Host:  GetEnv("DB_HOST", ""),
		Port:  GetEnv("DB_PORT", ""),
		User:  GetEnv("DB_USER", ""),
		Pass:  GetEnv("DB_PASS", ""),
		Name:  GetEnv("DB_NAME", ""),
		Debug: false,
	}
}

var ServiceConfig = map[string]string{
	"employeeUrl":    GetEnv("EMPLOYEE_SERVICE", "https://gateway-dev.vietmoney.vn/bs/human-resource-mngmt/v1"),
	"fileManagerUrl": GetEnv("FILE_MANAGER_SERVICE", "https://gateway-dev.vietmoney.vn/file-manager/v1"),
	"masterDataUrl":  GetEnv("MASTER_DATA_SERVICE", "https://gateway-dev.vietmoney.vn/rd/party/master-data/v1"),
}

var GCSConfig = map[string]string{
	"projectId": GetEnv("GCS_PROJECT_ID", "vietmoney-183803"),
	"buket":     GetEnv("GBS_STORAGE_BUKET", "hrm-vietmoney-dev"),
}

var DefaultDataConfig = map[string]string{
	"createdBy":   GetEnv("DEFAULT_DATA_CREATED_BY", "82eac640-bac5-4350-86cf-8c1a1a274e1e"),
	"createdName": GetEnv("DEFAULT_DATA_CREATED_NAME", "Vietmoney"),
}

type SubscriptionConfigItem struct {
	Key    string
	Action string
}

var SubscriptionConfig = []SubscriptionConfigItem{
	{
		Key:    GetEnv("SUBSCRIPTION_EXPORT", "dev.file_manager.data.export.crm"),
		Action: "ExportCrm",
	},
	{
		Key:    GetEnv("SUBSCRIPTION_LOAN_ORDER_CREATED", "dev.loan.order.created.local"),
		Action: "OrderCreated",
	},
	{
		Key:    GetEnv("SUBSCRIPTION_DISBURSEMENT", "dev.loan.order.disbursed.crm.local"),
		Action: "OrderDisbursed",
	},
}

var TopicConfig = map[string]string{
	"customerOrderUpdated": GetEnv("TOPIC_CUSTOMER_ORDER_UPDATED", "dev.customer.order.updated"),
	"customerLoanUpdated":  GetEnv("TOPIC_LOAN_ORDER_UPDATED", "dev.loan.order.updated"),
}
