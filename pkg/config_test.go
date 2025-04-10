package pkg

import (
	"runtime"
	"testing"
	"time"

	"github.com/spf13/viper"

	"github.com/stretchr/testify/assert"
)

func TestNewConfigPprof(t *testing.T) {
	viper.Set("PPROF_ENABLE", true)
	viper.Set("PPROF_BASICAUTH_USERNAME", "john")
	viper.Set("PPROF_BASICAUTH_PASSWORD", "test")

	c := NewConfigPprof()

	assert.Equal(t, c.Enable, true)
	assert.Equal(t, c.BasicAuthUsername, "john")
	assert.Equal(t, c.BasicAuthPassword, "test")
}

func TestNewConfigCORS(t *testing.T) {
	viper.Set("CORS_ALLOWED_ORIGINS", []string{"http://localhost:3000"})
	viper.Set("CORS_ALLOWED_METHODS", []string{"GET", "POST", "PUT", "DELETE"})
	viper.Set("CORS_ALLOWED_HEADERS", []string{"Accept", "Content-Type", "Authorization"})
	viper.Set("CORS_EXPOSED_HEADERS", []string{})
	viper.Set("CORS_ALLOW_CREDENTIALS", true)
	viper.Set("CORS_MAX_AGE", 10)

	c := NewConfigCORS()

	assert.Equal(t, c.AllowedOrigins, []string{"http://localhost:3000"})
	assert.Equal(t, c.AllowedMethods, []string{"GET", "POST", "PUT", "DELETE"})
	assert.Equal(t, c.AllowedHeaders, []string{"Accept", "Content-Type", "Authorization"})
	assert.Equal(t, c.ExposedHeaders, []string{})
	assert.Equal(t, c.AllowCredentials, true)
	assert.Equal(t, c.MaxAge, 10)
}

func TestNewConfigJWTWithCorrectParameters(t *testing.T) {
	// HS512
	viper.Set("JWT_ALGO", "HS512")
	viper.Set("JWT_SECRET", "mySecret")
	viper.Set("JWT_PRIVATE_KEY_PATH", "")
	viper.Set("JWT_PUBLIC_KEY_PATH", "")
	viper.Set("JWT_LIFETIME", 10)

	c, err := NewConfigJWT()

	assert.Nil(t, err)
	assert.Equal(t, c.Algorithm, "HS512")
	assert.Equal(t, c.SecretKey, "mySecret")
	assert.Equal(t, c.PrivateKeyPath, "")
	assert.Equal(t, c.PublicKeyPath, "")
	assert.Equal(t, c.Lifetime, 10*time.Hour)

	// ES384
	viper.Set("JWT_ALGO", "ES384")
	viper.Set("JWT_PRIVATE_KEY_PATH", "/path/to/private.key")
	viper.Set("JWT_PUBLIC_KEY_PATH", "/path/to/public.key")

	c, err = NewConfigJWT()

	assert.Nil(t, err)
	assert.Equal(t, c.Algorithm, "ES384")
	assert.Equal(t, c.SecretKey, "mySecret")
	assert.Equal(t, c.PrivateKeyPath, "/path/to/private.key")
	assert.Equal(t, c.PublicKeyPath, "/path/to/public.key")
	assert.Equal(t, c.Lifetime, 10*time.Hour)
}

func TestNewConfigJWTWithInvalidAlgo(t *testing.T) {
	viper.Set("JWT_ALGO", "HS256")
	viper.Set("JWT_SECRET", "mySecret")
	viper.Set("JWT_PRIVATE_KEY_PATH", "")
	viper.Set("JWT_PUBLIC_KEY_PATH", "")
	viper.Set("JWT_LIFETIME", 10)

	_, err := NewConfigJWT()

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "invalid JWT algorithm")
}

func TestNewConfigJWTWithEmptyHS512Secret(t *testing.T) {
	viper.Set("JWT_ALGO", "HS512")
	viper.Set("JWT_SECRET", "")
	viper.Set("JWT_PRIVATE_KEY_PATH", "")
	viper.Set("JWT_PUBLIC_KEY_PATH", "")
	viper.Set("JWT_LIFETIME", 10)

	_, err := NewConfigJWT()

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "missing JWT secret")
}

func TestNewConfigJWTWithEmptyES384KeyPaths(t *testing.T) {
	// Two paths empty
	viper.Set("JWT_ALGO", "ES384")
	viper.Set("JWT_SECRET", "")
	viper.Set("JWT_PRIVATE_KEY_PATH", "")
	viper.Set("JWT_PUBLIC_KEY_PATH", "")
	viper.Set("JWT_LIFETIME", 10)

	_, err := NewConfigJWT()

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "missing JWT private or public key path")

	// Public key path empty
	viper.Set("JWT_PRIVATE_KEY_PATH", "/path/to/private.key")

	_, err = NewConfigJWT()

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "missing JWT private or public key path")

	// Private key path empty
	viper.Set("JWT_PRIVATE_KEY_PATH", "")
	viper.Set("JWT_PUBLIC_KEY_PATH", "/path/to/public.key")

	_, err = NewConfigJWT()

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "missing JWT private or public key path")
}

func TestNewConfigLogWithCorrectParameters(t *testing.T) {
	// With file output
	viper.Set("LOG_LEVEL", "debug")
	viper.Set("LOG_OUTPUTS", "file")
	viper.Set("LOG_PATH", "/path/to/log")
	viper.Set("LOG_ACCESS_ENABLE", true)

	c, err := NewConfigLog()

	assert.Nil(t, err)
	assert.Equal(t, c.Level, "debug")
	assert.Equal(t, c.Outputs, []string{"file"})
	assert.Equal(t, c.Path, "/path/to/log")
	assert.Equal(t, c.EnableAccessLog, true)

	// With stdout output
	viper.Set("LOG_OUTPUTS", "stdout")
	c, err = NewConfigLog()

	assert.Nil(t, err)
	assert.Equal(t, c.Level, "debug")
	assert.Equal(t, c.Outputs, []string{"stdout"})
	assert.Equal(t, c.Path, "/path/to/log")
	assert.Equal(t, c.EnableAccessLog, true)

	// With both outputs
	viper.Set("LOG_OUTPUTS", "file stdout")

	c, err = NewConfigLog()
	assert.Nil(t, err)
	assert.Equal(t, c.Outputs, []string{"file", "stdout"})
}

func TestNewConfigLogWithInvalidLevel(t *testing.T) {
	// With file output
	viper.Set("LOG_LEVEL", "test")
	viper.Set("LOG_OUTPUTS", "stdout")
	viper.Set("LOG_PATH", "")
	viper.Set("LOG_ACCESS_ENABLE", true)

	_, err := NewConfigLog()

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "invalid log level")
}

func TestNewConfigLogWithInvalidOutputs(t *testing.T) {
	// With one output
	viper.Set("LOG_LEVEL", "info")
	viper.Set("LOG_OUTPUTS", "stdin")
	viper.Set("LOG_PATH", "")
	viper.Set("LOG_ACCESS_ENABLE", true)

	_, err := NewConfigLog()

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "invalid log outputs")

	// With two outputs
	viper.Set("LOG_LEVEL", "info")
	viper.Set("LOG_OUTPUTS", "stdout stdin")

	_, err = NewConfigLog()

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "invalid log outputs")
}

func TestNewConfigLogWithInvalidFilePath(t *testing.T) {
	// With file output
	viper.Set("LOG_LEVEL", "error")
	viper.Set("LOG_OUTPUTS", "file")
	viper.Set("LOG_PATH", "")
	viper.Set("LOG_ACCESS_ENABLE", true)

	_, err := NewConfigLog()

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "missing log path")
}

func TestConfigDatabaseWithCorrectParameters(t *testing.T) {
	viper.Set("DB_DRIVER", "mysql")
	viper.Set("DB_HOST", "localhost")
	viper.Set("DB_USERNAME", "root")
	viper.Set("DB_PASSWORD", "root")
	viper.Set("DB_PORT", 3306)
	viper.Set("DB_DATABASE", "test")
	viper.Set("DB_CHARSET", "utf8mb4")
	viper.Set("DB_COLLATION", "utf8mb4_general_ci")
	viper.Set("DB_LOCATION", "UTC")
	viper.Set("DB_MAX_IDLE_CONNS", 10)
	viper.Set("DB_MAX_OPEN_CONNS", 100)
	viper.Set("DB_CONN_MAX_LIFETIME", 1)
	viper.Set("DB_CONN_MAX_IDLE_TIME", 1)

	c, err := NewConfigDatabase()

	assert.Nil(t, err)
	assert.Equal(t, c.Driver, "mysql")
	assert.Equal(t, c.Host, "localhost")
	assert.Equal(t, c.Username, "root")
	assert.Equal(t, c.Password, "root")
	assert.Equal(t, c.Port, 3306)
	assert.Equal(t, c.Database, "test")
	assert.Equal(t, c.Charset, "utf8mb4")
	assert.Equal(t, c.Collation, "utf8mb4_general_ci")
	assert.Equal(t, c.Location, "UTC")
	assert.Equal(t, c.MaxIdleConns, 10)
	assert.Equal(t, c.MaxOpenConns, 100)
	assert.Equal(t, c.ConnMaxLifetime, 1*time.Hour)
	assert.Equal(t, c.ConnMaxIdleTime, 1*time.Hour)
}

func TestConfigDatabaseWithInvalidDriver(t *testing.T) {
	viper.Set("DB_DRIVER", "sqlite3")
	viper.Set("DB_HOST", "localhost")
	viper.Set("DB_USERNAME", "root")
	viper.Set("DB_PASSWORD", "root")
	viper.Set("DB_PORT", 3306)
	viper.Set("DB_DATABASE", "test")
	viper.Set("DB_CHARSET", "utf8mb4")
	viper.Set("DB_COLLATION", "utf8mb4_general_ci")
	viper.Set("DB_LOCATION", "UTC")
	viper.Set("DB_MAX_IDLE_CONNS", 10)
	viper.Set("DB_MAX_OPEN_CONNS", 100)
	viper.Set("DB_CONN_MAX_LIFETIME", 1)
	viper.Set("DB_CONN_MAX_IDLE_TIME", 1)

	_, err := NewConfigDatabase()

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "invalid database driver")
}

func TestConfigDatabaseWithInvalidLocation(t *testing.T) {
	viper.Set("DB_DRIVER", "mysql")
	viper.Set("DB_HOST", "localhost")
	viper.Set("DB_USERNAME", "root")
	viper.Set("DB_PASSWORD", "root")
	viper.Set("DB_PORT", 3306)
	viper.Set("DB_DATABASE", "test")
	viper.Set("DB_CHARSET", "utf8mb4")
	viper.Set("DB_COLLATION", "utf8mb4_general_ci")
	viper.Set("DB_LOCATION", "Europe/Paris")
	viper.Set("DB_MAX_IDLE_CONNS", 10)
	viper.Set("DB_MAX_OPEN_CONNS", 100)
	viper.Set("DB_CONN_MAX_LIFETIME", 1)
	viper.Set("DB_CONN_MAX_IDLE_TIME", 1)

	_, err := NewConfigDatabase()

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "invalid database location")
}

func TestConfigDatabaseWithEmptyDatabase(t *testing.T) {
	viper.Set("DB_DRIVER", "mysql")
	viper.Set("DB_HOST", "localhost")
	viper.Set("DB_USERNAME", "root")
	viper.Set("DB_PASSWORD", "root")
	viper.Set("DB_PORT", 3306)
	viper.Set("DB_DATABASE", "")
	viper.Set("DB_CHARSET", "utf8mb4")
	viper.Set("DB_COLLATION", "utf8mb4_general_ci")
	viper.Set("DB_LOCATION", "Local")
	viper.Set("DB_MAX_IDLE_CONNS", 10)
	viper.Set("DB_MAX_OPEN_CONNS", 100)
	viper.Set("DB_CONN_MAX_LIFETIME", 1)
	viper.Set("DB_CONN_MAX_IDLE_TIME", 1)

	_, err := NewConfigDatabase()

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "missing database name")
}

func TestNewConfigServerWithCorrectParameters(t *testing.T) {
	viper.Set("SERVER_ADDR", "localhost")
	viper.Set("SERVER_PORT", 8080)
	viper.Set("SERVER_TIMEOUT", 10)
	viper.Set("SERVER_MAX_REQUEST_SIZE", 1)
	viper.Set("SERVER_BASICAUTH_USERNAME", "")
	viper.Set("SERVER_BASICAUTH_PASSWORD", "")
	viper.Set("SERVER_MAX_CPU", 0)

	c, err := NewConfigServer()

	assert.Nil(t, err)
	assert.Equal(t, c.Addr, "localhost")
	assert.Equal(t, c.Port, 8080)
	assert.Equal(t, c.Timeout, 10)
	assert.Equal(t, c.MaxRequestSize, int64(1))
	assert.Equal(t, c.BasicAuthUsername, "")
	assert.Equal(t, c.BasicAuthPassword, "")
	assert.Equal(t, c.MaxCPU, runtime.NumCPU())
}

func TestNewConfigServerWithEmptyAddress(t *testing.T) {
	viper.Set("SERVER_ADDR", "")
	viper.Set("SERVER_PORT", 8080)
	viper.Set("SERVER_TIMEOUT", 10)
	viper.Set("SERVER_BASICAUTH_USERNAME", "")
	viper.Set("SERVER_BASICAUTH_PASSWORD", "")

	_, err := NewConfigServer()

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "missing server address")
}

func TestNewConfigServerWithEmptyPort(t *testing.T) {
	viper.Set("SERVER_ADDR", "localhost")
	viper.Set("SERVER_PORT", 0)
	viper.Set("SERVER_TIMEOUT", 10)
	viper.Set("SERVER_BASICAUTH_USERNAME", "")
	viper.Set("SERVER_BASICAUTH_PASSWORD", "")

	_, err := NewConfigServer()

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "missing server port")
}

func TestNewConfigServerWithGreaterCPUNumber(t *testing.T) {
	viper.Set("SERVER_ADDR", "localhost")
	viper.Set("SERVER_PORT", 8080)
	viper.Set("SERVER_TIMEOUT", 10)
	viper.Set("SERVER_MAX_REQUEST_SIZE", 1)
	viper.Set("SERVER_BASICAUTH_USERNAME", "")
	viper.Set("SERVER_BASICAUTH_PASSWORD", "")
	viper.Set("SERVER_MAX_CPU", 1_000)

	c, err := NewConfigServer()

	assert.Nil(t, err)
	assert.Equal(t, c.Addr, "localhost")
	assert.Equal(t, c.Port, 8080)
	assert.Equal(t, c.Timeout, 10)
	assert.Equal(t, c.MaxRequestSize, int64(1))
	assert.Equal(t, c.BasicAuthUsername, "")
	assert.Equal(t, c.BasicAuthPassword, "")
	assert.Equal(t, c.MaxCPU, runtime.NumCPU())
}

func TestNewConfigServerWithCorrectCPUNumber(t *testing.T) {
	viper.Set("SERVER_ADDR", "localhost")
	viper.Set("SERVER_PORT", 8080)
	viper.Set("SERVER_TIMEOUT", 10)
	viper.Set("SERVER_MAX_REQUEST_SIZE", 1)
	viper.Set("SERVER_BASICAUTH_USERNAME", "")
	viper.Set("SERVER_BASICAUTH_PASSWORD", "")
	viper.Set("SERVER_MAX_CPU", 1)

	c, err := NewConfigServer()

	assert.Nil(t, err)
	assert.Equal(t, c.Addr, "localhost")
	assert.Equal(t, c.Port, 8080)
	assert.Equal(t, c.Timeout, 10)
	assert.Equal(t, c.MaxRequestSize, int64(1))
	assert.Equal(t, c.BasicAuthUsername, "")
	assert.Equal(t, c.BasicAuthPassword, "")
	assert.Equal(t, c.MaxCPU, 1)
}

func TestNewConfigGorm(t *testing.T) {
	viper.Set("GORM_LOG_LEVEL", "info")
	viper.Set("GORM_LOG_OUTPUT", "stdout")
	viper.Set("GORM_LOG_FILE_NAME", "")
	viper.Set("GORM_SLOW_THRESHOLD", "200ms")

	c, err := NewConfigGorm()

	assert.Nil(t, err)
	assert.Equal(t, c.LogLevel, "info")
	assert.Equal(t, c.LogOutput, "stdout")
	assert.Equal(t, c.LogFileName, "")
	assert.Equal(t, c.SlowThreshold, 200*time.Millisecond)
}

func TestNewConfigGormWithInvalidLevel(t *testing.T) {
	viper.Set("GORM_LOG_LEVEL", "fatal")
	viper.Set("GORM_LOG_OUTPUT", "stdout")
	viper.Set("GORM_LOG_FILE_NAME", "")
	viper.Set("GORM_SLOW_THRESHOLD", "200ms")

	_, err := NewConfigGorm()

	assert.NotNil(t, err)
}

func TestNewConfigGormWithInvalidOutput(t *testing.T) {
	viper.Set("GORM_LOG_LEVEL", "info")
	viper.Set("GORM_LOG_OUTPUT", "test")
	viper.Set("GORM_LOG_FILE_NAME", "")
	viper.Set("GORM_SLOW_THRESHOLD", "200ms")

	_, err := NewConfigGorm()

	assert.NotNil(t, err)
}

func TestConfigDatabaseDsn(t *testing.T) {
	viper.Set("DB_DRIVER", "mysql")
	viper.Set("DB_HOST", "localhost")
	viper.Set("DB_USERNAME", "root")
	viper.Set("DB_PASSWORD", "root")
	viper.Set("DB_PORT", 3306)
	viper.Set("DB_DATABASE", "test")
	viper.Set("DB_CHARSET", "utf8mb4")
	viper.Set("DB_COLLATION", "utf8mb4_general_ci")
	viper.Set("DB_LOCATION", "UTC")
	viper.Set("DB_MAX_IDLE_CONNS", 10)
	viper.Set("DB_MAX_OPEN_CONNS", 100)
	viper.Set("DB_CONN_MAX_LIFETIME", 1)
	viper.Set("DB_CONN_MAX_IDLE_TIME", 1)

	c, _ := NewConfigDatabase()
	dns, err := c.DSN()

	assert.Nil(t, err)
	assert.Equal(t, dns, "root:root@tcp(localhost:3306)/test?parseTime=True&charset=utf8mb4&collation=utf8mb4_general_ci&loc=UTC")
}

func TestConfigDatabaseDsnInvalid(t *testing.T) {
	viper.Set("DB_DRIVER", "mysql")
	viper.Set("DB_HOST", "localhost")
	viper.Set("DB_USERNAME", "")
	viper.Set("DB_PASSWORD", "root")
	viper.Set("DB_PORT", 3306)
	viper.Set("DB_DATABASE", "test")
	viper.Set("DB_CHARSET", "utf8mb4")
	viper.Set("DB_COLLATION", "utf8mb4_general_ci")
	viper.Set("DB_LOCATION", "UTC")
	viper.Set("DB_MAX_IDLE_CONNS", 10)
	viper.Set("DB_MAX_OPEN_CONNS", 100)
	viper.Set("DB_CONN_MAX_LIFETIME", 1)
	viper.Set("DB_CONN_MAX_IDLE_TIME", 1)

	c, _ := NewConfigDatabase()
	_, err := c.DSN()
	assert.NotNil(t, err)

	// Empty password
	viper.Set("DB_USERNAME", "root")
	viper.Set("DB_PASSWORD", "")

	c, _ = NewConfigDatabase()
	_, err = c.DSN()
	assert.NotNil(t, err)

	// Empty port
	viper.Set("DB_PASSWORD", "root")
	viper.Set("DB_PORT", 0)

	c, _ = NewConfigDatabase()
	_, err = c.DSN()
	assert.NotNil(t, err)

	// Empty host
	viper.Set("DB_PORT", 3306)
	viper.Set("DB_HOST", "")

	c, _ = NewConfigDatabase()
	_, err = c.DSN()
	assert.NotNil(t, err)
}
