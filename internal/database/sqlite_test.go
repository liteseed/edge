package database

import (
	"os"
	"testing"

	"gotest.tools/v3/assert"
)

func TestSqlite(t *testing.T) {
	defer os.Remove("./data/testSQLITE")

	db := NewSqliteDatabase("./data/testSQLITE")

	err := db.Migrate()
	assert.NilError(t, err)
}
