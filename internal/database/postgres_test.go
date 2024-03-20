package database

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostgres(t *testing.T) {
	url := os.Getenv("DATABASE_URL")
	db := Postgres(url)

	err := db.Migrate()
	assert.NoError(t, err)
}
