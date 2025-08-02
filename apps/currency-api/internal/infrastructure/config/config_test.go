package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	originalEnv := make(map[string]string)
	envVars := []string{
		"PORT", "GIN_MODE", "LOG_LEVEL", "OPEN_EXCHANGE_API_KEY",
		"OPEN_EXCHANGE_BASE_URL", "REDIS_URL", "ENV",
	}

	for _, env := range envVars {
		originalEnv[env] = os.Getenv(env)
	}

	defer func() {
		for _, env := range envVars {
			if val, exists := originalEnv[env]; exists {
				os.Setenv(env, val)
			} else {
				os.Unsetenv(env)
			}
		}
	}()

	tests := []struct {
		name     string
		envVars  map[string]string
		expected *Config
		hasError bool
	}{
		{
			name: "default configuration",
			envVars: map[string]string{
				"PORT":                   "",
				"GIN_MODE":               "",
				"LOG_LEVEL":              "",
				"OPEN_EXCHANGE_API_KEY":  "",
				"OPEN_EXCHANGE_BASE_URL": "",
				"REDIS_URL":              "",
				"ENV":                    "",
			},
			expected: &Config{
				Port:                "8080",
				GinMode:             "debug",
				LogLevel:            "info",
				OpenExchangeAPIKey:  "",
				OpenExchangeBaseURL: "https://openexchangerates.org/api",
				RedisURL:            "redis://localhost:6379",
				Environment:         "development",
			},
		},
		{
			name: "custom configuration",
			envVars: map[string]string{
				"PORT":                   "3000",
				"GIN_MODE":               "release",
				"LOG_LEVEL":              "debug",
				"OPEN_EXCHANGE_API_KEY":  "test-api-key",
				"OPEN_EXCHANGE_BASE_URL": "https://custom-api.com",
				"REDIS_URL":              "redis://custom:6380",
				"ENV":                    "production",
			},
			expected: &Config{
				Port:                "3000",
				GinMode:             "release",
				LogLevel:            "debug",
				OpenExchangeAPIKey:  "test-api-key",
				OpenExchangeBaseURL: "https://custom-api.com",
				RedisURL:            "redis://custom:6380",
				Environment:         "production",
			},
		},
		{
			name: "test mode configuration",
			envVars: map[string]string{
				"PORT":                   "8081",
				"GIN_MODE":               "test",
				"LOG_LEVEL":              "error",
				"ENV":                    "test",
				"OPEN_EXCHANGE_API_KEY":  "",
				"OPEN_EXCHANGE_BASE_URL": "",
				"REDIS_URL":              "",
			},
			expected: &Config{
				Port:                "8081",
				GinMode:             "test",
				LogLevel:            "error",
				OpenExchangeAPIKey:  "",
				OpenExchangeBaseURL: "https://openexchangerates.org/api",
				RedisURL:            "redis://localhost:6379",
				Environment:         "test",
			},
		},
		{
			name: "invalid port - non-numeric",
			envVars: map[string]string{
				"PORT": "invalid-port",
			},
			hasError: true,
		},
		{
			name: "invalid gin mode",
			envVars: map[string]string{
				"GIN_MODE": "invalid-mode",
			},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for key, value := range tt.envVars {
				if value == "" {
					os.Unsetenv(key)
				} else {
					os.Setenv(key, value)
				}
			}

			config, err := Load()

			if tt.hasError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, config)

			assert.Equal(t, tt.expected.Port, config.Port)
			assert.Equal(t, tt.expected.GinMode, config.GinMode)
			assert.Equal(t, tt.expected.LogLevel, config.LogLevel)
			assert.Equal(t, tt.expected.OpenExchangeAPIKey, config.OpenExchangeAPIKey)
			assert.Equal(t, tt.expected.OpenExchangeBaseURL, config.OpenExchangeBaseURL)
			assert.Equal(t, tt.expected.RedisURL, config.RedisURL)
			assert.Equal(t, tt.expected.Environment, config.Environment)
		})
	}
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name          string
		config        *Config
		expectedError string
	}{
		{
			name: "valid configuration",
			config: &Config{
				Port:                "8080",
				GinMode:             "debug",
				LogLevel:            "info",
				OpenExchangeAPIKey:  "test-key",
				OpenExchangeBaseURL: "https://api.example.com",
				RedisURL:            "redis://localhost:6379",
				Environment:         "development",
			},
		},
		{
			name: "empty port",
			config: &Config{
				Port:     "",
				GinMode:  "debug",
				LogLevel: "info",
			},
			expectedError: "PORT cannot be empty",
		},
		{
			name: "invalid gin mode",
			config: &Config{
				Port:     "8080",
				GinMode:  "invalid",
				LogLevel: "info",
			},
			expectedError: "GIN_MODE must be one of: debug, release, test",
		},
		{
			name: "empty log level",
			config: &Config{
				Port:     "8080",
				GinMode:  "debug",
				LogLevel: "",
			},
			expectedError: "LOG_LEVEL cannot be empty",
		},
		{
			name: "non-numeric port",
			config: &Config{
				Port:     "not-a-number",
				GinMode:  "debug",
				LogLevel: "info",
			},
			expectedError: "PORT must be a valid number",
		},
		{
			name: "negative port should still validate",
			config: &Config{
				Port:     "-1",
				GinMode:  "debug",
				LogLevel: "info",
			},
		},
		{
			name: "zero port should still validate",
			config: &Config{
				Port:     "0",
				GinMode:  "debug",
				LogLevel: "info",
			},
		},
		{
			name: "high port number should validate",
			config: &Config{
				Port:     "65535",
				GinMode:  "debug",
				LogLevel: "info",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestConfig_IsProduction(t *testing.T) {
	tests := []struct {
		name        string
		environment string
		ginMode     string
		expected    bool
	}{
		{
			name:        "production environment",
			environment: "production",
			ginMode:     "debug",
			expected:    true,
		},
		{
			name:        "release gin mode",
			environment: "development",
			ginMode:     "release",
			expected:    true,
		},
		{
			name:        "both production",
			environment: "production",
			ginMode:     "release",
			expected:    true,
		},
		{
			name:        "development environment",
			environment: "development",
			ginMode:     "debug",
			expected:    false,
		},
		{
			name:        "test environment",
			environment: "test",
			ginMode:     "test",
			expected:    false,
		},
		{
			name:        "staging environment",
			environment: "staging",
			ginMode:     "debug",
			expected:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{
				Environment: tt.environment,
				GinMode:     tt.ginMode,
			}

			result := config.IsProduction()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetEnv(t *testing.T) {
	originalValue := os.Getenv("TEST_ENV_VAR")
	defer func() {
		if originalValue != "" {
			os.Setenv("TEST_ENV_VAR", originalValue)
		} else {
			os.Unsetenv("TEST_ENV_VAR")
		}
	}()

	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		expected     string
	}{
		{
			name:         "environment variable exists",
			key:          "TEST_ENV_VAR",
			defaultValue: "default",
			envValue:     "custom-value",
			expected:     "custom-value",
		},
		{
			name:         "environment variable does not exist",
			key:          "TEST_ENV_VAR",
			defaultValue: "default",
			envValue:     "",
			expected:     "default",
		},
		{
			name:         "empty default value",
			key:          "TEST_ENV_VAR",
			defaultValue: "",
			envValue:     "",
			expected:     "",
		},
		{
			name:         "empty environment value should use default",
			key:          "TEST_ENV_VAR",
			defaultValue: "default",
			envValue:     "",
			expected:     "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
			} else {
				os.Unsetenv(tt.key)
			}

			result := getEnv(tt.key, tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConfig_LoadAndValidate_Integration(t *testing.T) {
	originalEnv := make(map[string]string)
	envVars := []string{"PORT", "GIN_MODE", "LOG_LEVEL", "ENV"}

	for _, env := range envVars {
		originalEnv[env] = os.Getenv(env)
	}

	defer func() {
		for _, env := range envVars {
			if val, exists := originalEnv[env]; exists {
				os.Setenv(env, val)
			} else {
				os.Unsetenv(env)
			}
		}
	}()

	t.Run("valid production config", func(t *testing.T) {
		os.Setenv("PORT", "8080")
		os.Setenv("GIN_MODE", "release")
		os.Setenv("LOG_LEVEL", "info")
		os.Setenv("ENV", "production")

		config, err := Load()

		require.NoError(t, err)
		assert.Equal(t, "8080", config.Port)
		assert.Equal(t, "release", config.GinMode)
		assert.Equal(t, "info", config.LogLevel)
		assert.Equal(t, "production", config.Environment)
		assert.True(t, config.IsProduction())
	})

	t.Run("valid development config", func(t *testing.T) {
		os.Setenv("PORT", "3000")
		os.Setenv("GIN_MODE", "debug")
		os.Setenv("LOG_LEVEL", "debug")
		os.Setenv("ENV", "development")

		config, err := Load()

		require.NoError(t, err)
		assert.Equal(t, "3000", config.Port)
		assert.Equal(t, "debug", config.GinMode)
		assert.Equal(t, "debug", config.LogLevel)
		assert.Equal(t, "development", config.Environment)
		assert.False(t, config.IsProduction())
	})

	t.Run("invalid config should fail validation", func(t *testing.T) {
		os.Setenv("PORT", "invalid-port")
		os.Setenv("GIN_MODE", "debug")
		os.Setenv("LOG_LEVEL", "info")

		_, err := Load()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "config validation failed")
	})
}

func TestConfig_EnvironmentSpecificBehavior(t *testing.T) {
	originalEnv := os.Getenv("ENV")
	defer func() {
		if originalEnv != "" {
			os.Setenv("ENV", originalEnv)
		} else {
			os.Unsetenv("ENV")
		}
	}()

	environments := []struct {
		env             string
		expectedDefault string
	}{
		{"development", "development"},
		{"test", "test"},
		{"staging", "staging"},
		{"production", "production"},
	}

	for _, envTest := range environments {
		t.Run("environment_"+envTest.env, func(t *testing.T) {
			os.Unsetenv("GIN_MODE")
			os.Unsetenv("LOG_LEVEL")
			os.Unsetenv("PORT")
			os.Setenv("ENV", envTest.env)

			config, err := Load()

			require.NoError(t, err)
			assert.Equal(t, envTest.expectedDefault, config.Environment)
			assert.Equal(t, "8080", config.Port)
			assert.Equal(t, "debug", config.GinMode)
			assert.Equal(t, "info", config.LogLevel)
		})
	}
}

func TestConfig_AllFieldsLoaded(t *testing.T) {
	envVars := map[string]string{
		"PORT":                   "9000",
		"GIN_MODE":               "release",
		"LOG_LEVEL":              "warn",
		"OPEN_EXCHANGE_API_KEY":  "secret-key-123",
		"OPEN_EXCHANGE_BASE_URL": "https://custom-exchange-api.com/v2",
		"REDIS_URL":              "redis://redis-server:6380/1",
		"ENV":                    "staging",
	}

	originalEnv := make(map[string]string)
	for key := range envVars {
		originalEnv[key] = os.Getenv(key)
	}

	defer func() {
		for key := range envVars {
			if val, exists := originalEnv[key]; exists {
				os.Setenv(key, val)
			} else {
				os.Unsetenv(key)
			}
		}
	}()

	for key, value := range envVars {
		os.Setenv(key, value)
	}

	config, err := Load()

	require.NoError(t, err)
	assert.Equal(t, "9000", config.Port)
	assert.Equal(t, "release", config.GinMode)
	assert.Equal(t, "warn", config.LogLevel)
	assert.Equal(t, "secret-key-123", config.OpenExchangeAPIKey)
	assert.Equal(t, "https://custom-exchange-api.com/v2", config.OpenExchangeBaseURL)
	assert.Equal(t, "redis://redis-server:6380/1", config.RedisURL)
	assert.Equal(t, "staging", config.Environment)
}
