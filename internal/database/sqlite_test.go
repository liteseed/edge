package database

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSqlite(t *testing.T) {
	defer os.RemoveAll("./temp")
	_ = os.Mkdir("./temp", os.ModePerm)

	db := Sqlite("./temp/data")

	err := db.Migrate()
	assert.NoError(t, err)
}
