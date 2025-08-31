package config

type Config struct {
	LogLevel       string `env:"LOG_LEVEL"`
	Environment    string `env:"ENVIRONMENT"`
	ServiceName    string `env:"SERVICE_NAME"`
	ServiceVersion string `env:"SERVICE_VERSION"`
	HTTPPort       int    `env:"HTTP_PORT"`
	MetricPort     int    `env:"METRIC_PORT"`
	Security       *Security
	DB             *DB
	Redis          *Redis
}

type Security struct {
	APIKey string `env:"API_KEY"`
}

type DB struct {
	DriverName             string `env:"DB_DRIVER_NAME"`
	ConnectionURL          string `env:"DB_CONNECTION_URL" json:"-"` //zap ignore
	MigrationURL           string `env:"DB_MIGRATION_URL" json:"-"`  //zap ignore
	MaxOpenConnNumber      int    `env:"DB_MAX_OPEN_CONN_NUMBER"`
	MaxIdleConnNumber      int    `env:"DB_MAX_IDLE_CONN_NUMBER"`
	ConnMaxLifeTimeSeconds int64  `env:"DB_CONN_MAX_LIFE_TIME_SECONDS"`
}

type Redis struct {
	ConnectionURL       string `env:"REDIS_CONNECTION_URL"`
	Password            string `env:"REDIS_PASSWORD" json:"-"`
	DB                  int    `env:"REDIS_DB"`
	PoolSize            int    `env:"REDIS_POOL_SIZE"`
	DialTimeoutSeconds  int    `env:"REDIS_DIAL_TIMEOUT_SECONDS"`
	ReadTimeoutSeconds  int    `env:"REDIS_READ_TIMEOUT_SECONDS"`
	WriteTimeoutSeconds int    `env:"REDIS_WRITE_TIMEOUT_SECONDS"`
}
