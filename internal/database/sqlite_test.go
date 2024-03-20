package database

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSqlite(t *testing.T) {
	defer os.Remove("./temp/testSQLITE")

	db := Sqlite("./temp/testSQLITE")

	err := db.Migrate()
	assert.NoError(t, err)
}
