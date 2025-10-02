package config

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	Database  DatabaseConfig  `mapstructure:"database"`
	JWT       JWTConfig       `mapstructure:"jwt"`
	WebSocket WebSocketConfig `mapstructure:"websocket"`
	CORS      CORSConfig      `mapstructure:"cors"`
	Logging   LoggingConfig   `mapstructure:"logging"`
	App       App             `mapstructure:"app"`
	Storage   StorageConfig   `mapstructure:"storage"`
	Minio     MinioConfig     `mapstructure:"minio"`
	Qiniu     QiniuConfig     `mapstructure:"qiniu"`
}

type ServerConfig struct {
	Port         string        `mapstructure:"port"`
	Host         string        `mapstructure:"host"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

type DatabaseConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	DBName          string        `mapstructure:"dbname"`
	SSLMode         string        `mapstructure:"sslmode"`
	Timezone        string        `mapstructure:"timezone"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

type JWTConfig struct {
	Secret      string `mapstructure:"secret"`
	ExpireHours int    `mapstructure:"expire_hours"`
	Issuer      string `mapstructure:"issuer"`
}

type WebSocketConfig struct {
	ReadBufferSize  int           `mapstructure:"read_buffer_size"`
	WriteBufferSize int           `mapstructure:"write_buffer_size"`
	ReadDeadline    time.Duration `mapstructure:"read_deadline"`
	WriteDeadline   time.Duration `mapstructure:"write_deadline"`
	PingPeriod      time.Duration `mapstructure:"ping_period"`
}

type CORSConfig struct {
	AllowedOrigins []string `mapstructure:"allowed_origins"`
	AllowedMethods []string `mapstructure:"allowed_methods"`
	AllowedHeaders []string `mapstructure:"allowed_headers"`
}

type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
}

type App struct {
	Name        string `mapstructure:"name"`
	Version     string `mapstructure:"version"`
	Debug       bool   `mapstructure:"debug"`
	Environment string `mapstructure:"environment"`
}

type StorageConfig struct {
	Type string `mapstructure:"type"` // "minio" or "qiniu"
}

type MinioConfig struct {
	Endpoint   string `mapstructure:"endpoint"`
	AccessKey  string `mapstructure:"access_key"`
	SecretKey  string `mapstructure:"secret_key"`
	BucketName string `mapstructure:"bucket_name"`
	UseSSL     bool   `mapstructure:"use_ssl"`
	Region     string `mapstructure:"region"`
}

type QiniuConfig struct {
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	Bucket    string `mapstructure:"bucket"`
	Domain    string `mapstructure:"domain"`
	Region    string `mapstructure:"region"`
	UseHTTPS  bool   `mapstructure:"use_https"`
}

var GlobalConfig *Config

// LoadConfig loads configuration from config.yaml file
func LoadConfig(configPath ...string) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if len(configPath) > 0 {
		viper.AddConfigPath(configPath[0])
	} else {
		viper.AddConfigPath(".")
		viper.AddConfigPath("./config")
		viper.AddConfigPath("../")
	}

	// Set default values
	setDefaults()

	// Read environment variables
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	GlobalConfig = &config
	log.Printf("Configuration loaded successfully from: %s", viper.ConfigFileUsed())
	return &config, nil
}

// setDefaults sets default configuration values
func setDefaults() {
	viper.SetDefault("server.port", ":8080")
	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("server.read_timeout", "30s")
	viper.SetDefault("server.write_timeout", "30s")

	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("database.timezone", "UTC")
	viper.SetDefault("database.max_idle_conns", 10)
	viper.SetDefault("database.max_open_conns", 100)
	viper.SetDefault("database.conn_max_lifetime", "3600s")

	viper.SetDefault("jwt.secret", "default-secret-change-this")
	viper.SetDefault("jwt.expire_hours", 24)
	viper.SetDefault("jwt.issuer", "chatapp")

	viper.SetDefault("websocket.read_buffer_size", 1024)
	viper.SetDefault("websocket.write_buffer_size", 1024)
	viper.SetDefault("websocket.read_deadline", "60s")
	viper.SetDefault("websocket.write_deadline", "10s")
	viper.SetDefault("websocket.ping_period", "54s")

	viper.SetDefault("cors.allowed_origins", []string{"*"})
	viper.SetDefault("cors.allowed_methods", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	viper.SetDefault("cors.allowed_headers", []string{"Origin", "Content-Type", "Authorization"})

	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "json")
	viper.SetDefault("logging.output", "stdout")

	viper.SetDefault("app.name", "ChatApp")
	viper.SetDefault("app.version", "1.0.0")
	viper.SetDefault("app.debug", false)
	viper.SetDefault("app.environment", "development")

	viper.SetDefault("storage.type", "minio")

	viper.SetDefault("minio.endpoint", "127.0.0.1:9000")
	viper.SetDefault("minio.access_key", "minioadmin")
	viper.SetDefault("minio.secret_key", "minioadmin")
	viper.SetDefault("minio.bucket_name", "chatapp")
	viper.SetDefault("minio.use_ssl", false)

	viper.SetDefault("qiniu.access_key", "")
	viper.SetDefault("qiniu.secret_key", "")
	viper.SetDefault("qiniu.bucket", "")
	viper.SetDefault("qiniu.domain", "")
	viper.SetDefault("qiniu.region", "south-china")
	viper.SetDefault("qiniu.use_https", true)
}

// GetDatabaseDSN returns the database connection string
func (c *Config) GetDatabaseDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.DBName,
		c.Database.SSLMode,
		c.Database.Timezone,
	)
}
