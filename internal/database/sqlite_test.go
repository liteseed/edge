package database

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSqlite(t *testing.T) {
	defer os.Remove("./data/testSQLITE")

	db := NewSqliteDatabase("./data/testSQLITE")

	err := db.Migrate()
	quy.NoError(t, err)
}
