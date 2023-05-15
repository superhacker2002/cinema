package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateConfig(t *testing.T) {
	t.Run("successful validation", func(t *testing.T) {
		config := Config{
			Port:      "8080",
			Db:        "localhost:3306/mydb",
			JWTSecret: "secret-key",
			TokenExp:  "24",
		}
		assert.NoError(t, config.Validate())
	})

	t.Run("missing server port", func(t *testing.T) {
		config := Config{
			Port:      "",
			Db:        "localhost:3306/mydb",
			JWTSecret: "secret-key",
			TokenExp:  "24",
		}
		assert.EqualError(t, config.Validate(), ErrNoPort.Error())
	})

	t.Run("missing database URL", func(t *testing.T) {
		config := Config{
			Port:      "8080",
			Db:        "",
			JWTSecret: "secret-key",
			TokenExp:  "24",
		}
		assert.EqualError(t, config.Validate(), ErrNoDataBaseURL.Error())
	})

	t.Run("missing JWT secret key", func(t *testing.T) {
		config := Config{
			Port:      "8080",
			Db:        "localhost:9999/mydb",
			JWTSecret: "",
			TokenExp:  "24",
		}
		assert.EqualError(t, config.Validate(), ErrNoJWTSecret.Error())
	})

	t.Run("missing JWT expiration time", func(t *testing.T) {
		config := Config{
			Port:      "8080",
			Db:        "localhost:9999/mydb",
			JWTSecret: "secret-key",
			TokenExp:  "",
		}
		assert.EqualError(t, config.Validate(), ErrNoTokenExpiration.Error())
	})
}
