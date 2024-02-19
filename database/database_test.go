package database

import (
	"testing"

	"github.com/liteseed/bungo/database/schema"
	"github.com/stretchr/testify/assert"
)

func TestSqlite(t *testing.T) {
	db := NewSqliteDatabase("testSqlite")
	err := db.Migrate()
	assert.NoError(t, err)
	ord := &schema.Order{}
	err = db.DB.First(ord).Error
	assert.NoError(t, err)
	t.Log(ord)
}
