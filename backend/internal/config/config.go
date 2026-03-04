package config

import (
	"strings"
	"time"

	"linux-iso-manager/internal/constants"

	"github.com/spf13/viper"
)

// Config holds all application configuration.
type Config struct {
	Log       LogConfig
	Server    ServerConfig
	Database  DatabaseConfig
	Download  DownloadConfig
	WebSocket WebSocketConfig
}

// ServerConfig holds HTTP server configuration.
type ServerConfig struct {
	Port            string
	CORSOrigins     []string
	RequireProxyAuth bool
	CreateISOAuthEnabled bool
	BasicAuthUsername string
	BasicAuthPassword string
	LDAPAuthEnabled  bool
	LDAPURL          string
	LDAPBindDN       string
	LDAPBindPassword string
	LDAPUsersBaseDN  string
	LDAPUserFilter   string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

// DatabaseConfig holds database configuration.
type DatabaseConfig struct {
	Path            string
	JournalMode     string
	BusyTimeout     time.Duration
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// DownloadConfig holds download manager configuration.
type DownloadConfig struct {
	DataDir                  string
	WorkerCount              int
	QueueBuffer              int
	MaxRetries               int
	RetryDelay               time.Duration
	BufferSize               int
	ProgressUpdateInterval   time.Duration
	ProgressPercentThreshold int
	CancellationWait         time.Duration
}

// WebSocketConfig holds WebSocket configuration.
type WebSocketConfig struct {
	BroadcastChannelSize int
}

// LogConfig holds logging configuration.
type LogConfig struct {
	Level  string // debug, info, warn, error
	Format string // json, text
}

// Load loads configuration from environment variables with defaults using Viper.
func Load() *Config {
	v := viper.New()

	// Set defaults for Server
	v.SetDefault("PORT", constants.DefaultPort)
	v.SetDefault("READ_TIMEOUT_SEC", constants.DefaultReadTimeoutSec)
	v.SetDefault("WRITE_TIMEOUT_SEC", constants.DefaultWriteTimeoutSec)
	v.SetDefault("IDLE_TIMEOUT_SEC", constants.DefaultIdleTimeoutSec)
	v.SetDefault("SHUTDOWN_TIMEOUT_SEC", constants.DefaultShutdownTimeoutSec)
	v.SetDefault("CORS_ORIGINS", "http://localhost:3000,http://localhost:5173,http://localhost:8080")
	v.SetDefault("REQUIRE_PROXY_AUTH", false)
	v.SetDefault("CREATE_ISO_AUTH_ENABLED", false)
	v.SetDefault("BASIC_AUTH_USERNAME", "")
	v.SetDefault("BASIC_AUTH_PASSWORD", "")
	v.SetDefault("LDAP_AUTH_ENABLED", false)
	v.SetDefault("LDAP_URL", "")
	v.SetDefault("LDAP_BIND_DN", "")
	v.SetDefault("LDAP_BIND_PASSWORD", "")
	v.SetDefault("LDAP_USERS_BASE_DN", "")
	v.SetDefault("LDAP_USER_FILTER", "(|(uid={user})(username={user})(mail={user}))")
	v.SetDefault("CLOUDRON_LDAP_URL", "")
	v.SetDefault("CLOUDRON_LDAP_BIND_DN", "")
	v.SetDefault("CLOUDRON_LDAP_BIND_PASSWORD", "")
	v.SetDefault("CLOUDRON_LDAP_USERS_BASE_DN", "")

	// Set defaults for Database
	v.SetDefault("DB_PATH", "")
	v.SetDefault("DB_BUSY_TIMEOUT_MS", constants.DefaultBusyTimeoutMs)
	v.SetDefault("DB_JOURNAL_MODE", constants.DefaultJournalMode)
	v.SetDefault("DB_MAX_OPEN_CONNS", constants.DefaultMaxOpenConns)
	v.SetDefault("DB_MAX_IDLE_CONNS", constants.DefaultMaxIdleConns)
	v.SetDefault("DB_CONN_MAX_LIFETIME_MIN", constants.DefaultConnMaxLifetimeMin)
	v.SetDefault("DB_CONN_MAX_IDLE_TIME_MIN", constants.DefaultConnMaxIdleTimeMin)

	// Set defaults for Download
	v.SetDefault("DATA_DIR", "./data")
	v.SetDefault("WORKER_COUNT", constants.DefaultWorkerCount)
	v.SetDefault("QUEUE_BUFFER", constants.DefaultQueueBuffer)
	v.SetDefault("MAX_RETRIES", constants.DefaultMaxRetries)
	v.SetDefault("RETRY_DELAY_MS", constants.DefaultRetryDelayMs)
	v.SetDefault("BUFFER_SIZE", constants.DefaultDownloadBufferSize)
	v.SetDefault("PROGRESS_UPDATE_INTERVAL_SEC", 1)
	v.SetDefault("PROGRESS_PERCENT_THRESHOLD", constants.DefaultProgressPercentThreshold)
	v.SetDefault("CANCELLATION_WAIT_MS", constants.DefaultCancellationWaitMs)

	// Set defaults for WebSocket
	v.SetDefault("WS_BROADCAST_SIZE", constants.DefaultBroadcastChannelSize)

	// Set defaults for Logging
	v.SetDefault("LOG_LEVEL", "info")
	v.SetDefault("LOG_FORMAT", "text")

	// Bind environment variables
	v.AutomaticEnv()

	// Parse CORS origins
	corsOriginsStr := v.GetString("CORS_ORIGINS")
	corsOrigins := strings.Split(corsOriginsStr, ",")

	return &Config{
		Server: ServerConfig{
			Port:            v.GetString("PORT"),
			RequireProxyAuth: v.GetBool("REQUIRE_PROXY_AUTH"),
			CreateISOAuthEnabled: v.GetBool("CREATE_ISO_AUTH_ENABLED"),
			BasicAuthUsername: v.GetString("BASIC_AUTH_USERNAME"),
			BasicAuthPassword: v.GetString("BASIC_AUTH_PASSWORD"),
			LDAPAuthEnabled:  v.GetBool("LDAP_AUTH_ENABLED") || v.GetString("CLOUDRON_LDAP_URL") != "",
			LDAPURL:          coalesce(v.GetString("CLOUDRON_LDAP_URL"), v.GetString("LDAP_URL")),
			LDAPBindDN:       coalesce(v.GetString("CLOUDRON_LDAP_BIND_DN"), v.GetString("LDAP_BIND_DN")),
			LDAPBindPassword: coalesce(v.GetString("CLOUDRON_LDAP_BIND_PASSWORD"), v.GetString("LDAP_BIND_PASSWORD")),
			LDAPUsersBaseDN:  coalesce(v.GetString("CLOUDRON_LDAP_USERS_BASE_DN"), v.GetString("LDAP_USERS_BASE_DN")),
			LDAPUserFilter:   v.GetString("LDAP_USER_FILTER"),
			ReadTimeout:     time.Duration(v.GetInt("READ_TIMEOUT_SEC")) * time.Second,
			WriteTimeout:    time.Duration(v.GetInt("WRITE_TIMEOUT_SEC")) * time.Second,
			IdleTimeout:     time.Duration(v.GetInt("IDLE_TIMEOUT_SEC")) * time.Second,
			ShutdownTimeout: time.Duration(v.GetInt("SHUTDOWN_TIMEOUT_SEC")) * time.Second,
			CORSOrigins:     corsOrigins,
		},
		Database: DatabaseConfig{
			Path:            v.GetString("DB_PATH"),
			BusyTimeout:     time.Duration(v.GetInt("DB_BUSY_TIMEOUT_MS")) * time.Millisecond,
			JournalMode:     v.GetString("DB_JOURNAL_MODE"),
			MaxOpenConns:    v.GetInt("DB_MAX_OPEN_CONNS"),
			MaxIdleConns:    v.GetInt("DB_MAX_IDLE_CONNS"),
			ConnMaxLifetime: time.Duration(v.GetInt("DB_CONN_MAX_LIFETIME_MIN")) * time.Minute,
			ConnMaxIdleTime: time.Duration(v.GetInt("DB_CONN_MAX_IDLE_TIME_MIN")) * time.Minute,
		},
		Download: DownloadConfig{
			DataDir:                  v.GetString("DATA_DIR"),
			WorkerCount:              v.GetInt("WORKER_COUNT"),
			QueueBuffer:              v.GetInt("QUEUE_BUFFER"),
			MaxRetries:               v.GetInt("MAX_RETRIES"),
			RetryDelay:               time.Duration(v.GetInt("RETRY_DELAY_MS")) * time.Millisecond,
			BufferSize:               v.GetInt("BUFFER_SIZE"),
			ProgressUpdateInterval:   time.Duration(v.GetInt("PROGRESS_UPDATE_INTERVAL_SEC")) * time.Second,
			ProgressPercentThreshold: v.GetInt("PROGRESS_PERCENT_THRESHOLD"),
			CancellationWait:         time.Duration(v.GetInt("CANCELLATION_WAIT_MS")) * time.Millisecond,
		},
		WebSocket: WebSocketConfig{
			BroadcastChannelSize: v.GetInt("WS_BROADCAST_SIZE"),
		},
		Log: LogConfig{
			Level:  v.GetString("LOG_LEVEL"),
			Format: v.GetString("LOG_FORMAT"),
		},
	}
}
