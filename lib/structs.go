package coolconf

type AppConfig struct {
	Profiling             string `json:"PROFILING" envconfig:"PROFILING" default:"0"`
	CacheExpiresInSeconds int    `json:"CACHE_EXPIRES" envconfig:"CACHE_EXPIRES" default:"600"`
	ManagerAuthEnabled    bool   `json:"MANAGER_AUTH_ENABLED" envconfig:"MANAGER_AUTH_ENABLED" default:"true"`

	VerticalHeader      string `json:"VERTICAL_HEADER" envconfig:"VERTICAL_HEADER" default:"X-Vertical"`
	SystemTokenHeader   string `json:"SYSTEM_TOKEN_HEADER" envconfig:"SYSTEM_TOKEN_HEADER" default:"X-Requested-With"`
	AuthorizationHeader string `json:"AUTHORIZATION_HEADER" envconfig:"AUTHORIZATION_HEADER" default:"Authorization"`

	SystemToken string `json:"SYSTEM_TOKEN" envconfig:"SYSTEM_TOKEN" default:"systoken"`

	DatabaseURL             string `json:"DATABASE_URL" envconfig:"DATABASE_URL"`
	DatabaseMaxConns        int    `json:"DATABASE_MAX_CONNS" envconfig:"DATABASE_MAX_CONNS"`
	DatabaseMaxConnLifetime int    `json:"DATABASE_MAX_CONN_LIFETIME" envconfig:"DATABASE_MAX_CONN_LIFETIME"`
	DatabaseMaxConnIdleTime int    `json:"DATABASE_MAX_IDLE_TIME" envconfig:"DATABASE_MAX_IDLE_TIME"`

	LogDestination string `json:"LOG_DEST" envconfig:"LOG_DEST" default:"stdout"`
	LogFormat      string `json:"LOG_FORMAT" envconfig:"LOG_FORMAT" default:"text"`
	LogLevel       string `json:"LOG_LEVEL" envconfig:"LOG_LEVEL" default:"info"`

	Integrator          string `json:"INTERGRATOR" envconfig:"INTERGRATOR" default:"mundipagg"`
	MundipaggSecret     string `json:"MUNDIPAGG_SECRET" envconfig:"MUNDIPAGG_SECRET"`
	MundipaggBaseURL    string `json:"MUNDIPAGG_BASE_URL" envconfig:"MUNDIPAGG_BASE_URL" default:"https://api.mundipagg.com/core/v1"`
	WebhookUser         string `json:"WEBHOOK_USER" envconfig:"WEBHOOK_USER"`
	WebhookPassword     string `json:"WEBHOOK_PASSWORD" envconfig:"WEBHOOK_PASSWORD"`
	WebhookProxyFile    string `json:"WEBHOOK_PROXY_FILE" envconfig:"WEBHOOK_PROXY_FILE" default:"/tmp/webhook.txt"`
	WebhookProxyEnabled bool   `json:"WEBHOOK_PROXY_ENABLED" envconfig:"WEBHOOK_PROXY_ENABLED"`

	CardTokenURL string `json:"CARDTOKEN_URL" envconfig:"CARDTOKEN_URL"`
	NFeURL       string `json:"NFE_URL" envconfig:"NFE_URL" default:"http://tecprime.com.br:9000"`
	NFeEnabled   string `json:"NFE_ENABLED" envconfig:"NFE_ENABLED" default:"0"`

	AccountsClientSystemToken string `json:"ACCOUNTS_CLIENT_SYSTEM_TOKEN" envconfig:"ACCOUNTS_CLIENT_SYSTEM_TOKEN"`
	AccountsClientHost        string `json:"ACCOUNTS_CLIENT_HOST" envconfig:"ACCOUNTS_CLIENT_HOST"`

	SMTPUser   string `json:"SMTP_USER" envconfig:"SMTP_USER"`
	SMTPSecret string `json:"SMTP_SECRET" envconfig:"SMTP_SECRET"`
	SMTPSender string `json:"SMTP_SENDER" envconfig:"SMTP_SENDER"`
	SMTPFrom   string `json:"SMTP_FROM" envconfig:"SMTP_FROM"`
	SMTPServer string `json:"SMTP_SERVER" envconfig:"SMTP_SERVER"`
	SMTPPort   string `json:"SMTP_PORT" envconfig:"SMTP_PORT"`
}
