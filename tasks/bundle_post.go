package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
)

// A list of task types.
const (
	TypeBundlePost = "bundle:post"
)

type BundlePostPayload struct {
	OrderID string
}

func NewDataPostTask(orderId string) (*asynq.Task, error) {
	payload, err := json.Marshal(BundlePostPayload{OrderID: orderId})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeBundlePost, payload), nil
}

func HandleDataPostTask(ctx context.Context, t *asynq.Task) error {
	var payload BundlePostPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	return nil
}

func (task *Task) ProcessTask(oid string) error {
	stores, err := task.db.GetStores(oid)
	if err != nil {
		log.Println(err)
		return err
	}
	for _, store := range *stores {
		data, err := task.store.Get(store.ID.String())
		if err != nil {
			log.Println(err)
		}
	}
}
