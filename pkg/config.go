package pkg

import (
	"fmt"
	"go-clean-api/pkg/apperr"
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
		return nil, apperr.NewAppErr(fmt.Errorf("error in configuration"), "missing server address", nil, nil)
	}

	if port == 0 {
		return nil, apperr.NewAppErr(fmt.Errorf("error in configuration"), "missing server port", nil, nil)
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
		return nil, apperr.NewAppErr(fmt.Errorf("error in configuration"), "invalid database driver", nil, nil)
	}

	if location != "UTC" && location != "Local" {
		return nil, apperr.NewAppErr(fmt.Errorf("error in configuration"), "invalid database location", nil, nil)
	}

	if database == "" {
		return nil, apperr.NewAppErr(fmt.Errorf("error in configuration"), "invalid database name", nil, nil)
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

// DSN returns the DSN if the configuration is OK or an error in other case
func (c *ConfigDatabase) DSN() (dsn string, err error) {
	if c.Host == "" || c.Port == 0 || c.Username == "" || c.Password == "" {
		return dsn, apperr.NewAppErr(fmt.Errorf("error in configuration"), "invalid database configuration", nil, nil)
	}

	dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=True",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.Database)
	if c.Charset != "" {
		dsn += fmt.Sprintf("&charset=%s", c.Charset)
	}
	if c.Collation != "" {
		dsn += fmt.Sprintf("&collation=%s", c.Collation)
	}
	if c.Location != "" {
		dsn += fmt.Sprintf("&loc=%s", c.Location)
	}
	return
}

// ConfigGorm represents the configuration of gorm
type ConfigGorm struct {
	// Log level (silent | error | warn | info)
	LogLevel string

	// Log ouput (stdout | file)
	LogOutput string

	// Log file name
	LogFileName string

	// Slow threshold
	SlowThreshold time.Duration
}

// NewConfigGorm creates a new ConfigGorm instance
func NewConfigGorm() (*ConfigGorm, error) {
	level := viper.GetString("GORM_LOG_LEVEL")
	output := viper.GetString("GORM_LOG_OUTPUT")

	if level != "info" && level != "warn" && level != "error" && level != "silent" {
		return nil, apperr.NewAppErr(fmt.Errorf("error in configuration"), "invalid gorm log level", nil, nil)
	}

	if output != "stdout" && output != "file" {
		return nil, apperr.NewAppErr(fmt.Errorf("error in configuration"), "invalid log outputs", nil, nil)
	}

	return &ConfigGorm{
		LogLevel:      level,
		LogOutput:     output,
		LogFileName:   viper.GetString("GORM_LOG_FILE_NAME"),
		SlowThreshold: viper.GetDuration("GORM_SLOW_THRESHOLD"),
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
		return nil, apperr.NewAppErr(fmt.Errorf("error in configuration"), "invalid log level", nil, nil)
	}

	for _, output := range outputs {
		if output != "stdout" && output != "file" {
			return nil, apperr.NewAppErr(fmt.Errorf("error in configuration"), "invalid log outputs", nil, nil)
		}
	}

	if goutils.StringInSlice("file", outputs) && path == "" {
		return nil, apperr.NewAppErr(fmt.Errorf("error in configuration"), "missing log path", nil, nil)
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
		return nil, apperr.NewAppErr(fmt.Errorf("error in configuration"), "invalid JWT algorithm", nil, nil)
	}

	if algo == "HS512" && secret == "" {
		return nil, apperr.NewAppErr(fmt.Errorf("error in configuration"), "missing JWT secret", nil, nil)
	}

	if algo == "ES384" && (privateKeyPath == "" || publicKeyPath == "") {
		return nil, apperr.NewAppErr(fmt.Errorf("error in configuration"), "missing JWT private or public key path", nil, nil)
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

	// Gorm configuration
	Gorm ConfigGorm

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
		return nil, apperr.NewAppErr(err, "error when reading .env file", nil, nil)
	}

	jwtConfig, err := NewConfigJWT()
	if err != nil {
		return nil, apperr.NewAppErr(err, "error in JWT configuration", nil, nil)
	}

	logConfig, err := NewConfigLog()
	if err != nil {
		return nil, apperr.NewAppErr(err, "error in log configuration", nil, nil)
	}

	databaseConfig, err := NewConfigDatabase()
	if err != nil {
		return nil, apperr.NewAppErr(err, "error in database configuration", nil, nil)
	}

	gormConfig, err := NewConfigGorm()
	if err != nil {
		return nil, apperr.NewAppErr(err, "errorin GORM configuration", nil, nil)
	}

	serverConfig, err := NewConfigServer()
	if err != nil {
		return nil, apperr.NewAppErr(err, "error in server configuration", nil, nil)
	}

	return &Config{
		AppEnv:   viper.GetString("APP_ENV"),
		AppName:  viper.GetString("APP_NAME"),
		Server:   *serverConfig,
		Database: *databaseConfig,
		Gorm:     *gormConfig,
		Log:      *logConfig,
		JWT:      *jwtConfig,
		CORS:     *NewConfigCORS(),
		Pprof:    *NewConfigPprof(),
	}, nil
}
