package pkg

import (
	"fmt"
	"runtime"
	"time"

	"github.com/fabienbellanger/goutils"
	"github.com/spf13/viper"
)

// ConfigServer represents the configuration of the HTTP server
type ConfigServer struct {
	// Address
	Addr string

	// Port
	Port int

	// Timeout
	Timeout int

	// Max request size in KB (0 = unlimited)
	MaxRequestSize int64

	// Basic Auth username
	BasicAuthUsername string

	// Basic Auth password
	BasicAuthPassword string

	// Maximal number of CPUs (Mst be lower than the number of CPUs of the machine)
	MaxCPU int
}

// NewConfigServer creates a new ConfigServer instance
func NewConfigServer() (*ConfigServer, error) {
	addr := viper.GetString("SERVER_ADDR")
	port := viper.GetInt("SERVER_PORT")
	maxCPU := viper.GetInt("SERVER_MAX_CPU")

	if addr == "" {
		return nil, fmt.Errorf("missing server address")
	}

	if port == 0 {
		return nil, fmt.Errorf("missing server port")
	}

	defaultCPU := runtime.NumCPU()
	if maxCPU > defaultCPU || maxCPU < 1 {
		maxCPU = defaultCPU
	}

	return &ConfigServer{
		Addr:              addr,
		Port:              port,
		Timeout:           viper.GetInt("SERVER_TIMEOUT"),
		MaxRequestSize:    viper.GetInt64("SERVER_MAX_REQUEST_SIZE"),
		BasicAuthUsername: viper.GetString("SERVER_BASICAUTH_USERNAME"),
		BasicAuthPassword: viper.GetString("SERVER_BASICAUTH_PASSWORD"),
		MaxCPU:            maxCPU,
	}, nil
}

// ConfigDatabase represents the configuration of the database
type ConfigDatabase struct {
	// Driver
	Driver string

	// Host
	Host string

	// Username
	Username string

	// Password
	Password string

	// Port
	Port int

	// Database
	Database string

	// Charset
	Charset string

	// Collation
	Collation string

	// Location (UTC | Local)
	Location string

	// Max idle connections
	MaxIdleConns int

	// Max open connections
	MaxOpenConns int

	// Connection max lifetime
	ConnMaxLifetime time.Duration

	// Connection max idle time
	ConnMaxIdleTime time.Duration
}

// NewConfigDatabase creates a new ConfigDatabase instance
func NewConfigDatabase() (*ConfigDatabase, error) {
	driver := viper.GetString("DB_DRIVER")
	location := viper.GetString("DB_LOCATION")
	database := viper.GetString("DB_DATABASE")

	if driver != "mysql" {
		return nil, fmt.Errorf("invalid database driver")
	}

	if location != "UTC" && location != "Local" {
		return nil, fmt.Errorf("invalid database location")
	}

	if database == "" {
		return nil, fmt.Errorf("missing database name")
	}

	return &ConfigDatabase{
		Driver:          driver,
		Host:            viper.GetString("DB_HOST"),
		Username:        viper.GetString("DB_USERNAME"),
		Password:        viper.GetString("DB_PASSWORD"),
		Port:            viper.GetInt("DB_PORT"),
		Database:        database,
		Charset:         viper.GetString("DB_CHARSET"),
		Collation:       viper.GetString("DB_COLLATION"),
		Location:        location,
		MaxIdleConns:    viper.GetInt("DB_MAX_IDLE_CONNS"),
		MaxOpenConns:    viper.GetInt("DB_MAX_OPEN_CONNS"),
		ConnMaxLifetime: viper.GetDuration("DB_CONN_MAX_LIFETIME") * time.Hour,
		ConnMaxIdleTime: viper.GetDuration("DB_CONN_MAX_IDLE_TIME") * time.Hour,
	}, nil
}

// ConfigLog represents the configuration of the logs
type ConfigLog struct {
	// Path
	Path string

	// Outputs (stdout | file)
	Outputs []string

	// Level (debug | info | warn | error | fatal | panic)
	Level string

	// Enable access log
	EnableAccessLog bool
}

// NewConfigLog creates a new ConfigLog instance
func NewConfigLog() (*ConfigLog, error) {
	level := viper.GetString("LOG_LEVEL")
	outputs := viper.GetStringSlice("LOG_OUTPUTS")
	path := viper.GetString("LOG_PATH")

	if level != "debug" && level != "info" && level != "warn" && level != "error" && level != "fatal" && level != "panic" {
		return nil, fmt.Errorf("invalid log level")
	}

	for _, output := range outputs {
		if output != "stdout" && output != "file" {
			return nil, fmt.Errorf("invalid log outputs")
		}
	}

	if goutils.StringInSlice("file", outputs) && path == "" {
		return nil, fmt.Errorf("missing log path")
	}

	return &ConfigLog{
		Path:            path,
		Outputs:         outputs,
		Level:           level,
		EnableAccessLog: viper.GetBool("LOG_ACCESS_ENABLE"),
	}, nil
}

// ConfigJWT represents the configuration of the JWT
type ConfigJWT struct {
	// Algorithm (HS512 | ES384)
	Algorithm string

	// Lifetime (in hour)
	Lifetime time.Duration

	// Secret key
	SecretKey string

	// Private key path
	PrivateKeyPath string

	// Public key path
	PublicKeyPath string
}

// NewConfigJWT creates a new ConfigJWT instance
func NewConfigJWT() (*ConfigJWT, error) {
	algo := viper.GetString("JWT_ALGO")
	secret := viper.GetString("JWT_SECRET")
	privateKeyPath := viper.GetString("JWT_PRIVATE_KEY_PATH")
	publicKeyPath := viper.GetString("JWT_PUBLIC_KEY_PATH")

	if algo != "HS512" && algo != "ES384" {
		return nil, fmt.Errorf("invalid JWT algorithm")
	}

	if algo == "HS512" && secret == "" {
		return nil, fmt.Errorf("missing JWT secret")
	}

	if algo == "ES384" && (privateKeyPath == "" || publicKeyPath == "") {
		return nil, fmt.Errorf("missing JWT private or public key path")
	}

	return &ConfigJWT{
		Algorithm:      algo,
		Lifetime:       viper.GetDuration("JWT_LIFETIME") * time.Hour,
		SecretKey:      secret,
		PrivateKeyPath: privateKeyPath,
		PublicKeyPath:  publicKeyPath,
	}, nil
}

// ConfigCORS represents the configuration of the CORS
type ConfigCORS struct {
	// Allowed origins
	AllowedOrigins []string

	// Allowed methods
	AllowedMethods []string

	// Allowed headers
	AllowedHeaders []string

	// Exposed headers
	ExposedHeaders []string

	// Allow credentials
	AllowCredentials bool

	// Max age
	MaxAge int
}

// NewConfigCORS creates a new ConfigCORS instance
func NewConfigCORS() *ConfigCORS {
	return &ConfigCORS{
		AllowedOrigins:   viper.GetStringSlice("CORS_ALLOWED_ORIGINS"),
		AllowedMethods:   viper.GetStringSlice("CORS_ALLOWED_METHODS"),
		AllowedHeaders:   viper.GetStringSlice("CORS_ALLOWED_HEADERS"),
		ExposedHeaders:   viper.GetStringSlice("CORS_EXPOSED_HEADERS"),
		AllowCredentials: viper.GetBool("CORS_ALLOW_CREDENTIALS"),
		MaxAge:           viper.GetInt("CORS_MAX_AGE"),
	}
}

// ConfigPprof represents the configuration of the pprof
type ConfigPprof struct {
	// Enable pprof
	Enable bool

	// Basic Auth username
	BasicAuthUsername string

	// Basic Auth password
	BasicAuthPassword string
}

// NewConfigServer creates a new ConfigServer instance
func NewConfigPprof() *ConfigPprof {
	return &ConfigPprof{
		Enable:            viper.GetBool("PPROF_ENABLE"),
		BasicAuthUsername: viper.GetString("PPROF_BASICAUTH_USERNAME"),
		BasicAuthPassword: viper.GetString("PPROF_BASICAUTH_PASSWORD"),
	}
}

// Config represents the configuration of the application from the .env file
type Config struct {
	// Application environment (development, production or test)
	AppEnv string

	// Application name
	AppName string

	// Server configuration
	Server ConfigServer

	// Database configuration
	Database ConfigDatabase

	// Log configuration
	Log ConfigLog

	// JWT configuration
	JWT ConfigJWT

	// CORS configuration
	CORS ConfigCORS

	// Pprof configuration
	Pprof ConfigPprof
}

// NewConfig creates a new Config instance
func NewConfig(file string) (*Config, error) {
	// Read .env file
	viper.SetConfigFile(file)
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	jwtConfig, err := NewConfigJWT()
	if err != nil {
		return nil, err
	}

	logConfig, err := NewConfigLog()
	if err != nil {
		return nil, err
	}

	databaseConfig, err := NewConfigDatabase()
	if err != nil {
		return nil, err
	}

	serverConfig, err := NewConfigServer()
	if err != nil {
		return nil, err
	}

	return &Config{
		AppEnv:   viper.GetString("APP_ENV"),
		AppName:  viper.GetString("APP_NAME"),
		Server:   *serverConfig,
		Database: *databaseConfig,
		Log:      *logConfig,
		JWT:      *jwtConfig,
		CORS:     *NewConfigCORS(),
		Pprof:    *NewConfigPprof(),
	}, nil
}
