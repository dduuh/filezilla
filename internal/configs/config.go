package configs

import (
	"os"
	"time"

	"github.com/spf13/viper"
)

const (
	defaultHTTPPort        = "8082"
	defaultMaxHeaderBytes  = 1
	defaultReadWriteTimeot = 10 * time.Second
	defaultPSQLUser        = "postgres"
	defaultPSQLPassword    = "Backham2410+"
	defaultPSQLHost        = "localhost"
	defaultPSQLPort        = 5435
	defaultPSQLDBName      = "postgres"
	defaultPSQLSSLMode     = "disable"
	defaultMongoUri        = "mongodb://localhost:27017"
	defaultMongoUser       = "David3410"
	defaultMongoPassword   = "eGd1"
	defaultMongoDb         = "logs_db"

	defaultAccessTtl  = 15 * time.Minute
	defaultRefreshTtl = 24 * time.Hour * 30

	envLocal = "local"
)

type (
	Config struct {
		HttpCfg     HTTPConfig
		PostgresCfg PostgreSQLConfig
		JWTCfg      JWTConfig
		MongoCfg    MongoDBConfig
	}

	HTTPConfig struct {
		Host           string        `mapstructure:"host"`
		Port           string        `mapstructure:"port"`
		ReadTimeout    time.Duration `mapstructure:"readTimeout"`
		WriteTimeout   time.Duration `mapstructure:"writeTimeout"`
		MaxHeaderBytes int           `mapstructure:"maxHeaderBytes"`
	}

	JWTConfig struct {
		SigningKey string `mapstructure:"signingKey"`
		AccessTtl  string `mapstructure:"accessTtl"`
		RefreshTtl string `mapstructure:"refreshTtl"`
	}

	PostgreSQLConfig struct {
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		DBName   string `mapstructure:"dbname"`
		SSLMode  string `mapstructure:"sslmode"`
	}

	MongoDBConfig struct {
		Uri      string `mapstructure:"uri"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		Database string `mapstructure:"database"`
	}
)

func Init(configDir string) (*Config, error) {
	populateDefaultFiles()

	if err := parseConfigFiles(configDir, os.Getenv("APP_ENV")); err != nil {
		return nil, err
	}

	var cfg Config
	if err := unmarshal(&cfg); err != nil {
		return nil, err
	}

	setFromEnv(&cfg)

	return &cfg, nil
}

func setFromEnv(cfg *Config) {
	cfg.HttpCfg.Port = os.Getenv("HTTP_PORT")
	cfg.HttpCfg.Host = os.Getenv("HTTP_HOST")
	cfg.PostgresCfg.User = os.Getenv("POSTGRES_USER")
	cfg.PostgresCfg.Password = os.Getenv("POSTGRES_PASSWORD")
	cfg.PostgresCfg.Host = os.Getenv("POSTGRES_HOST")
	cfg.PostgresCfg.Port = os.Getenv("POSTGRES_PORT")
	cfg.PostgresCfg.DBName = os.Getenv("POSTGRES_DBNAME")
	cfg.PostgresCfg.SSLMode = os.Getenv("POSTGRES_SSLMODE")
	cfg.JWTCfg.AccessTtl = os.Getenv("ACCESS_TOKEN_TTL")
	cfg.JWTCfg.RefreshTtl = os.Getenv("REFRESH_TOKEN_TTL")
	cfg.MongoCfg.Uri = os.Getenv("MONGO_INITDB_ROOT_URI")
	cfg.MongoCfg.Username = os.Getenv("MONGO_INITDB_ROOT_USERNAME")
	cfg.MongoCfg.Password = os.Getenv("MONGO_INITDB_ROOT_PASSWORD")
	cfg.MongoCfg.Database = os.Getenv("MONGO_INITDB_ROOT_DB")
}

func unmarshal(cfg *Config) error {
	if err := viper.UnmarshalKey("auth", &cfg.JWTCfg); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("postgres_db", &cfg.PostgresCfg); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("mongodb", &cfg.MongoCfg); err != nil {
		return err
	}

	return viper.UnmarshalKey("http", &cfg.HttpCfg)
}

func parseConfigFiles(folder, env string) error {
	viper.AddConfigPath(folder)
	viper.SetConfigName("main")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if env == envLocal {
		return nil
	}

	viper.SetConfigName(env)

	return viper.MergeInConfig()
}

func populateDefaultFiles() {
	viper.SetDefault("http.port", defaultHTTPPort)
	viper.SetDefault("http.maxHeaderBytes", defaultMaxHeaderBytes)
	viper.SetDefault("http.readTimeout", defaultReadWriteTimeot)
	viper.SetDefault("http.writeTimeout", defaultReadWriteTimeot)
	viper.SetDefault("postgres_db.user", defaultPSQLUser)
	viper.SetDefault("postgres_db.password", defaultPSQLPassword)
	viper.SetDefault("postgres_db.host", defaultPSQLHost)
	viper.SetDefault("postgres_db.port", defaultPSQLPort)
	viper.SetDefault("postgres_db.dbname", defaultPSQLDBName)
	viper.SetDefault("postgres_db.sslmode", defaultPSQLSSLMode)
	viper.SetDefault("auth.accessTtl", defaultAccessTtl)
	viper.SetDefault("auth.refreshTtl", defaultRefreshTtl)
	viper.SetDefault("mongodb.uri", defaultMongoUri)
	viper.SetDefault("mongodb.username", defaultMongoUser)
	viper.SetDefault("mongodb.password", defaultMongoPassword)
	viper.SetDefault("mongodb.database", defaultMongoDb)
}
