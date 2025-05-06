package config

import (
	"fmt"
	"os"

	workflows "gitlab.com/clubhub.ai1/go-libraries/saga/client"

	"gitlab.com/clubhub.ai1/gommon/config"
)

type Configuration struct {
	MembersApi            MembersApi               `mapstructure:"members_api"`
	DeUnaApi              DeUnaApi                 `mapstructure:"deuna_api"`
	BillApi               BillApi                  `mapstructure:"bill_api"`
	MailService           MailService              `mapstructure:"mail_service"`
	NotificationLinkEmail NotificationLinkEmail    `mapstructure:"notification_link_email"`
	Server                Server                   `mapstructure:"server"`
	App                   App                      `mapstructure:"app"`
	Otel                  Otel                     `mapstructure:"otel"`
	TemporalConfig        workflows.TemporalConfig `mapstructure:"temporal"`
	Environment           string                   `mapstructure:"environment"`
	Db                    DB                       `mapstructure:"db"`
	Aws                   Aws                      `mapstructure:"aws"`
	DynamoDB              DynamoDB                 `mapstructure:"dynamodb"`
}

type Aws struct {
	Region               string `mapstructure:"region"                          required:"true"`
	AccessKeyID          string `mapstructure:"access_key_id"                   required:"true"`
	SecretAccessKey      string `mapstructure:"secret_access_key"            required:"true"`
	EndpointURL          string `mapstructure:"endpoint_url"                   required:"true"`
	RoleARN              string `mapstructure:"role_arn"                       required:"false"`
	ReceiptPaymentBucket string `mapstructure:"receipt_payment_bucket"         required:"true"`
	CDN_URL              string `mapstructure:"cdn_url"                        required:"true"`
}

type DynamoDB struct {
	PaymentsEventStoreTable string `mapstructure:"order_event_store_dynamodb_table" required:"true"`
}

type MembersApi struct {
	URL string `required:"true" mapstructure:"url"`
}

type BillApi struct {
	URL string `required:"true" mapstructure:"url"`
}

type MailService struct {
	URL string `required:"true" mapstructure:"url"`
}

type NotificationLinkEmail struct {
	Wallet         string `required:"true" mapstructure:"wallet"`
	ContactSupport string `required:"true" mapstructure:"contact_support"`
}

type DeUnaApi struct {
	URL          string `required:"true" mapstructure:"url"`
	ApiKey       string `required:"true" mapstructure:"api_key"`
	PublicKey    string `required:"true" mapstructure:"public_key"`
	PublicKeyPEM string `required:"true" mapstructure:"public_key_pem"`
}
type Server struct {
	Port     int    `required:"true" mapstructure:"port"`
	Host     string `required:"true" mapstructure:"host"`
	BasePath string `required:"true" mapstructure:"base_path"`
}

type App struct {
	ServiceName      string `mapstructure:"service_name" required:"true"`
	Postfix          string `mapstructure:"postfix" required:"true"`
	Country          string `mapstructure:"country" required:"true"`
	LoggerDebugMode  bool   `mapstructure:"logger_debug_mode" required:"true"`
	EventStoreDynamo bool   `mapstructure:"event_store_dynamo" required:"true"`
	ClubhubMainHost  string `mapstructure:"clubhub_main_host" required:"true"`
}

type DB struct {
	Host     string `mapstructure:"host" required:"true"`
	Port     int    `mapstructure:"port" required:"true"`
	User     string `mapstructure:"user" required:"true"`
	Password string `mapstructure:"password" required:"true"`
	DB       string `mapstructure:"name" required:"true"`
	PoolSize int    `mapstructure:"pool_size" required:"false"`
}

type Otel struct {
	ServiceName               string `mapstructure:"service_name" required:"true"`
	ExporterOtlpHeadersKey    string `mapstructure:"exporter_otlp_headers_key" required:"true"`
	ExporterOtlpEndpoint      string `mapstructure:"exporter_otlp_endpoint" required:"true"`
	AttributeValueLengthLimit string `mapstructure:"attribute_value_length_limit" required:"true"`
	ExporterProvider          string `mapstructure:"exporter_provider" required:"false"`
	Provider                  string `mapstructure:"provider" required:"true"`
	InternalCollectorGrpcURL  string `mapstructure:"internal_collector_grpc_url" required:"true"`
}

var configuration Configuration

func Config() Configuration {
	return configuration
}

func Environments() error {
	if os.Getenv("ENVIRONMENT") == "local" {
		if err := config.GenEnvsFromFile(".env"); err != nil {
			return fmt.Errorf("error reading config file %w", err)
		}
	}

	if errR := config.ReadConfigFromEnv(&configuration); errR != nil {
		return fmt.Errorf("error getting configurations from env %w", errR)
	}

	return nil
}

func IsLocal() bool {
	return Config().Environment == "local"
}
