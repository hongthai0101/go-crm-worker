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
			EmployeeUrl:    GetEnv("EMPLOYEE_SERVICE", ""),
			FileManagerUrl: GetEnv("FILE_MANAGER_SERVICE", ""),
			MasterDataUrl:  GetEnv("MASTER_DATA_SERVICE", ""),
		},
		GCSConfig: &GCSConfig{
			ProjectId: GetEnv("GCS_PROJECT_ID", ""),
			Buket:     GetEnv("GBS_STORAGE_BUKET", ""),
		},
		DefaultDataConfig: &DefaultDataConfig{
			CreatedBy:   GetEnv("DEFAULT_DATA_CREATED_BY", ""),
			CreatedName: GetEnv("DEFAULT_DATA_CREATED_NAME", ""),
		},
		SubscriptionConfig: [3]*SubscriptionConfigItem{
			{
				Key:    GetEnv("SUBSCRIPTION_EXPORT", ""),
				Action: "ExportCrm",
			},
			{
				Key:    GetEnv("SUBSCRIPTION_LOAN_ORDER_CREATED", ""),
				Action: "OrderCreated",
			},
			{
				Key:    GetEnv("SUBSCRIPTION_DISBURSEMENT", ""),
				Action: "OrderDisbursed",
			},
		},
		TopicConfig: &TopicConfig{
			CustomerOrderUpdated: GetEnv("TOPIC_CUSTOMER_ORDER_UPDATED", ""),
			CustomerLoanUpdated:  GetEnv("TOPIC_LOAN_ORDER_UPDATED", ""),
		},
	}
}
