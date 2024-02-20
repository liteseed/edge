package tasks

import (
	"github.com/liteseed/bungo/database"
	"github.com/liteseed/bungo/store"
)

type Task struct {
	db    *database.Database
	store *store.Store
}

func New(
	db *database.Database,
	store *store.Store,
) *Task {
	return &Task{db: db, store: store}
}
