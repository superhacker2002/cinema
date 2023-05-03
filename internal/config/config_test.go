package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateConfig(t *testing.T) {
	config := Config{
		Port:      "8080",
		Db:        "localhost:3306/mydb",
		JWTSecret: "secret-key",
	}
	assert.NoError(t, config.Validate())

	config = Config{
		Port:      "",
		Db:        "localhost:3306/mydb",
		JWTSecret: "secret-key",
	}
	assert.EqualError(t, config.Validate(), ErrNoPort.Error())

	config = Config{
		Port:      "8080",
		Db:        "",
		JWTSecret: "secret-key",
	}
	assert.EqualError(t, config.Validate(), ErrNoDataBaseURL.Error())

	config = Config{
		Port:      "8080",
		Db:        "localhost:9999/mydb",
		JWTSecret: "",
	}
	assert.EqualError(t, config.Validate(), ErrNoJWTSecret.Error())
}
