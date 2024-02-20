package database

import (
	"testing"

	"github.com/liteseed/bungo/database/schema"
	"gotest.tools/v3/assert"
)

func TestSqlite(t *testing.T) {
	db := NewSqliteDatabase("testSqlite")
	err := db.Migrate()
	assert.NilError(t, err)
	ord := &schema.Order{}
	err = db.DB.First(ord).Error
	assert.NilError(t, err)
	t.Log(ord)
}
