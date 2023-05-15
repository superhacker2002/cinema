package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateConfig(t *testing.T) {
	t.Run("successful validation", func(t *testing.T) {
		config := Config{
			Port: "8080",
			Db:   "localhost:3306/mydb",
		}
		assert.NoError(t, config.Validate())
	})

	t.Run("missing server port", func(t *testing.T) {
		config := Config{
			Port: "",
			Db:   "localhost:3306/mydb",
		}
		assert.EqualError(t, config.Validate(), ErrNoPort.Error())
	})

	t.Run("missing database URL", func(t *testing.T) {
		config := Config{
			Port: "8080",
			Db:   "",
		}
		assert.EqualError(t, config.Validate(), ErrNoDataBaseURL.Error())
	})
}
