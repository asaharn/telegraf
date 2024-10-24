package microsoft_fabric

import (
	"context"

	eventhub "github.com/Azure/azure-event-hubs-go/v3"
)

type eventHub struct {
	hub *eventhub.Hub
}

func (e *eventHub) GetHub(connectionString string) error {
	hub, err := eventhub.NewHubFromConnectionString(connectionString)
	if err != nil {
		return err
	}
	e.hub = hub
	return nil
}

func (e *eventHub) Close(ctx context.Context) error {
	return e.hub.Close(ctx)
}

func (e *eventHub) SendBatch(ctx context.Context, iterator eventhub.BatchIterator, opts ...eventhub.BatchOption) error {
	return e.hub.SendBatch(ctx, iterator, opts...)
}
