package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func clearEnvVars() {
	envVars := []string{
		"SERVICE_NAME",
		"SERVER_HOST",
		"SERVER_PORT",
		"DB_USER",
		"DB_PASSWORD",
		"DB_HOST",
		"DB_PORT",
		"DB_NAME",
		"DB_SSLMODE",
		"DB_TIMEZONE",
		"DB_AUTO_MIGRATE",
	}

	for _, envVar := range envVars {
		os.Unsetenv(envVar)
	}
}

func TestNewConfig_WithDefaults(t *testing.T) {
	clearEnvVars()

	config := NewConfig()

	assert.NotNil(t, config)
	assert.Equal(t, "simple-gin-crud", config.ServiceName)
	assert.Equal(t, "0.0.0.0", config.Server.Host)
	assert.Equal(t, "8080", config.Server.Port)
	assert.Equal(t, "", config.Database.User)
	assert.Equal(t, "", config.Database.Password)
	assert.Equal(t, "", config.Database.Host)
	assert.Equal(t, "", config.Database.Port)
	assert.Equal(t, "", config.Database.DBName)
	assert.Equal(t, "", config.Database.SSLMode)
	assert.Equal(t, "", config.Database.TimeZone)
	assert.False(t, config.Database.AutoMigrate)
}

func TestNewConfig_WithEnvironmentVariables(t *testing.T) {
	os.Setenv("SERVICE_NAME", "test-service")
	os.Setenv("SERVER_HOST", "localhost")
	os.Setenv("SERVER_PORT", "9090")
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_PASSWORD", "testpass")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("DB_SSLMODE", "disable")
	os.Setenv("DB_TIMEZONE", "UTC")
	os.Setenv("DB_AUTO_MIGRATE", "true")

	defer clearEnvVars()

	config := NewConfig()

	assert.NotNil(t, config)
	assert.Equal(t, "test-service", config.ServiceName)
	assert.Equal(t, "localhost", config.Server.Host)
	assert.Equal(t, "9090", config.Server.Port)
	assert.Equal(t, "testuser", config.Database.User)
	assert.Equal(t, "testpass", config.Database.Password)
	assert.Equal(t, "localhost", config.Database.Host)
	assert.Equal(t, "5432", config.Database.Port)
	assert.Equal(t, "testdb", config.Database.DBName)
	assert.Equal(t, "disable", config.Database.SSLMode)
	assert.Equal(t, "UTC", config.Database.TimeZone)
	assert.True(t, config.Database.AutoMigrate)
}

func TestGetValue_WithEnvironmentVariable(t *testing.T) {
	os.Setenv("TEST_KEY", "test_value")
	defer os.Unsetenv("TEST_KEY")

	result := getValue("TEST_KEY", "default_value")
	assert.Equal(t, "test_value", result)
}

func TestGetValue_WithoutEnvironmentVariable(t *testing.T) {
	os.Unsetenv("NONEXISTENT_KEY")

	result := getValue("NONEXISTENT_KEY", "default_value")
	assert.Equal(t, "default_value", result)
}

func TestGetValue_WithEmptyEnvironmentVariable(t *testing.T) {
	os.Setenv("EMPTY_KEY", "")
	defer os.Unsetenv("EMPTY_KEY")

	result := getValue("EMPTY_KEY", "default_value")
	assert.Equal(t, "", result)
}

func TestConfig_FieldTypes(t *testing.T) {
	clearEnvVars()

	config := NewConfig()

	assert.IsType(t, "", config.ServiceName)
	assert.IsType(t, "", config.Database.User)
	assert.IsType(t, "", config.Database.Password)
	assert.IsType(t, "", config.Database.Host)
	assert.IsType(t, "", config.Database.Port)
	assert.IsType(t, "", config.Database.DBName)
	assert.IsType(t, "", config.Database.SSLMode)
	assert.IsType(t, "", config.Database.TimeZone)
	assert.IsType(t, false, config.Database.AutoMigrate)
	assert.IsType(t, "", config.Server.Host)
	assert.IsType(t, "", config.Server.Port)
}
