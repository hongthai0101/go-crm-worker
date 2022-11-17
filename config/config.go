package config

type Config struct {
	DB                 *DBConfig
	ServiceConfig      *ServiceConfig
	GCSConfig          *GCSConfig
	DefaultDataConfig  *DefaultDataConfig
	SubscriptionConfig [3]*SubscriptionConfigItem
	TopicConfig        *TopicConfig
}

type DBConfig struct {
	Host  string
	Port  string
	User  string
	Pass  string
	Name  string
	Debug bool
}

type ServiceConfig struct {
	EmployeeUrl    string
	FileManagerUrl string
	MasterDataUrl  string
}

type GCSConfig struct {
	ProjectId string
	Buket     string
}

type DefaultDataConfig struct {
	CreatedBy   string
	CreatedName string
}

type TopicConfig struct {
	CustomerOrderUpdated string
	CustomerLoanUpdated  string
}

type SubscriptionConfigItem struct {
	Key    string
	Action string
}

func GetConfig() *Config {
	return &Config{
		DB: &DBConfig{
			Host:  GetEnv("DB_HOST", ""),
			Port:  GetEnv("DB_PORT", ""),
			User:  GetEnv("DB_USER", ""),
			Pass:  GetEnv("DB_PASS", ""),
			Name:  GetEnv("DB_NAME", ""),
			Debug: false,
		},
		ServiceConfig: &ServiceConfig{
			EmployeeUrl:    GetEnv("EMPLOYEE_SERVICE", "https://gateway-dev.vietmoney.vn/bs/human-resource-mngmt/v1"),
			FileManagerUrl: GetEnv("FILE_MANAGER_SERVICE", "https://gateway-dev.vietmoney.vn/file-manager/v1"),
			MasterDataUrl:  GetEnv("MASTER_DATA_SERVICE", "https://gateway-dev.vietmoney.vn/rd/party/master-data/v1"),
		},
		GCSConfig: &GCSConfig{
			ProjectId: GetEnv("GCS_PROJECT_ID", "vietmoney-183803"),
			Buket:     GetEnv("GBS_STORAGE_BUKET", "hrm-vietmoney-dev"),
		},
		DefaultDataConfig: &DefaultDataConfig{
			CreatedBy:   GetEnv("DEFAULT_DATA_CREATED_BY", "82eac640-bac5-4350-86cf-8c1a1a274e1e"),
			CreatedName: GetEnv("DEFAULT_DATA_CREATED_NAME", "Vietmoney"),
		},
		SubscriptionConfig: [3]*SubscriptionConfigItem{
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
		},
		TopicConfig: &TopicConfig{
			CustomerOrderUpdated: GetEnv("TOPIC_CUSTOMER_ORDER_UPDATED", "dev.customer.order.updated"),
			CustomerLoanUpdated:  GetEnv("TOPIC_LOAN_ORDER_UPDATED", "dev.loan.order.updated"),
		},
	}
}
