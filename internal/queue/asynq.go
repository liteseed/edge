package queue

import (
	"context"
	"encoding/json"
	"fmt"

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

type AsynqClient struct {
	client *asynq.Client
}

func NewAsynqClient(address string) *Queue {
	client := &AsynqClient{client: asynq.NewClient(asynq.RedisClientOpt{Addr: address})}
	return &Queue{Queue: client}
}

func (q *AsynqClient) Process() {

}

func (q *AsynqClient) Schedule() error {
	return nil
}

func (q *AsynqClient) Close() {

}
