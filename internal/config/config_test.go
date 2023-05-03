package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateConfig(t *testing.T) {
	config := Config{
		Port: "8080",
		Db:   "localhost:3306/mydb",
	}
	assert.NoError(t, config.Validate())

	config = Config{
		Port: "",
		Db:   "localhost:3306/mydb",
	}
	assert.EqualError(t, config.Validate(), ErrNoPort.Error())

	config = Config{
		Port: "8080",
		Db:   "",
	}
	assert.EqualError(t, config.Validate(), ErrNoDataBaseURL.Error())
}
