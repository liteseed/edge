package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostgres(t *testing.T) {

	db := Postgres("postgresql://shaurya.saklani130:1yvXLlTd9gSz@ep-dry-pond-a5alprmm.us-east-2.aws.neon.tech/arc?sslmode=require")

	err := db.Migrate()
	assert.NoError(t, err)
}
