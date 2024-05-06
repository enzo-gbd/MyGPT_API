// Package configs handles the loading and management of configuration settings from environment variables and configuration files.
package configs

import (
	"path/filepath"
	"runtime"
	"time"

	"github.com/spf13/viper"
)

// env defaults to "local", indicating the environment settings to load. This can be overridden by setting a different value in the environment variables.
var env = "local"

// Config defines the structure for configuration settings, supporting loading from environment variables and config files. Fields are mapped to specific environment variables.
type Config struct {
	DBHost         string `mapstructure:"POSTGRES_HOST"`     // DBHost represents the database server address.
	DBUserName     string `mapstructure:"POSTGRES_USER"`     // DBUserName represents the database user name.
	DBUserPassword string `mapstructure:"POSTGRES_PASSWORD"` // DBUserPassword represents the database user password.
	DBName         string `mapstructure:"POSTGRES_DATABASE"` // DBName represents the name of the database.
	DBPort         string `mapstructure:"POSTGRES_PORT"`     // DBPort represents the port on which the database server is running.
	DBSSLMode      string `mapstructure:"POSTGRES_SSLMODE"`  // DBSSLMode represents the SSL mode for database connections.

	ClientOrigin string `mapstructure:"CLIENT_ORIGIN"` // ClientOrigin specifies the CORS origin for client requests.

	AccessTokenPrivateKey  string        `mapstructure:"ACCESS_TOKEN_PRIVATE_KEY"`  // AccessTokenPrivateKey holds the private key for generating access tokens.
	AccessTokenPublicKey   string        `mapstructure:"ACCESS_TOKEN_PUBLIC_KEY"`   // AccessTokenPublicKey holds the public key for validating access tokens.
	RefreshTokenPrivateKey string        `mapstructure:"REFRESH_TOKEN_PRIVATE_KEY"` // RefreshTokenPrivateKey holds the private key for generating refresh tokens.
	RefreshTokenPublicKey  string        `mapstructure:"REFRESH_TOKEN_PUBLIC_KEY"`  // RefreshTokenPublicKey holds the public key for validating refresh tokens.
	AccessTokenExpiresIn   time.Duration `mapstructure:"ACCESS_TOKEN_EXPIRED_IN"`   // AccessTokenExpiresIn specifies the duration after which access tokens expire.
	RefreshTokenExpiresIn  time.Duration `mapstructure:"REFRESH_TOKEN_EXPIRED_IN"`  // RefreshTokenExpiresIn specifies the duration after which refresh tokens expire.
	AccessTokenMaxAge      int           `mapstructure:"ACCESS_TOKEN_MAXAGE"`       // AccessTokenMaxAge specifies the maximum age in seconds for access tokens.
	RefreshTokenMaxAge     int           `mapstructure:"REFRESH_TOKEN_MAXAGE"`      // RefreshTokenMaxAge specifies the maximum age in seconds for refresh tokens.
}

// getAbsoluteRootPath computes and returns the absolute path to the root directory of the project by examining the caller's location in the filesystem.
func getAbsoluteRootPath() string {
	_, b, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(b), "..")
}

// LoadConfig loads configuration settings based on the environment specified by `env`. It utilizes the Viper library to read and parse environment variables and configuration files.
// Returns a populated Config struct and any error encountered during the loading process.
func LoadConfig() (config Config, err error) {
	viper.AddConfigPath(getAbsoluteRootPath()) // Set the path for reading the configuration files.
	viper.SetConfigType("env")                 // Set the type of the configuration files (environment variables).
	viper.SetConfigName(env)                   // Set the name of the configuration file based on the current environment.

	viper.AutomaticEnv() // Automatically override configuration file values with environment variables that match.

	err = viper.ReadInConfig() // Read the configuration file.
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config) // Unmarshal the configuration file into the Config struct.
	return
}
