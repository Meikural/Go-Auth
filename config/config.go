package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	DBDriver                  string
	DBSource                  string
	JWTSecret                 string
	ServerPort                string
	Roles                     []string
	DefaultRegistrationRole   string
	SuperAdminEmail           string
	SuperAdminPassword        string
}

// Load reads configuration from environment variables
func Load() *Config {
	config := &Config{
		DBDriver:                getEnv("DB_DRIVER", "postgres"),
		DBSource:                getEnv("DB_SOURCE", ""),
		JWTSecret:               getEnv("JWT_SECRET", ""),
		ServerPort:              getEnv("SERVER_PORT", "8080"),
		DefaultRegistrationRole: getEnv("DEFAULT_REGISTRATION_ROLE", "User"),
		SuperAdminEmail:         getEnv("SUPER_ADMIN_EMAIL", ""),
		SuperAdminPassword:      getEnv("SUPER_ADMIN_PASSWORD", ""),
	}

	// Parse ROLES from env (JSON array)
	rolesEnv := getEnv("ROLES", `["Super Admin", "User"]`)
	if err := json.Unmarshal([]byte(rolesEnv), &config.Roles); err != nil {
		panic(fmt.Sprintf("Failed to parse ROLES environment variable: %v", err))
	}

	// Validate required fields
	if config.DBSource == "" {
		panic("DB_SOURCE environment variable is required")
	}
	if config.JWTSecret == "" {
		panic("JWT_SECRET environment variable is required")
	}
	if config.SuperAdminEmail == "" {
		panic("SUPER_ADMIN_EMAIL environment variable is required")
	}
	if config.SuperAdminPassword == "" {
		panic("SUPER_ADMIN_PASSWORD environment variable is required")
	}

	// Validate roles
	if len(config.Roles) == 0 {
		panic("ROLES must contain at least one role")
	}

	// Validate default registration role exists in roles
	roleExists := false
	for _, role := range config.Roles {
		if strings.EqualFold(role, config.DefaultRegistrationRole) {
			roleExists = true
			break
		}
	}
	if !roleExists {
		panic(fmt.Sprintf("DEFAULT_REGISTRATION_ROLE '%s' not found in ROLES", config.DefaultRegistrationRole))
	}

	return config
}

// getEnv retrieves an environment variable with a fallback default
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// String returns a formatted string representation of the config
func (c *Config) String() string {
	return fmt.Sprintf("Config{Driver: %s, Port: %s, Roles: %v, DefaultRole: %s}", 
		c.DBDriver, c.ServerPort, c.Roles, c.DefaultRegistrationRole)
}