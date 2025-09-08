package app

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestLoadEnv(t *testing.T) {
	err := godotenv.Load("../../.env")
	assert.NoError(t, err)
	assert.NotEmpty(t, os.Getenv("DB_HOST"))

}
