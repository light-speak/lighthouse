package messaging

import (
	"context"
	"strings"
)

type Broker interface {
	Init(cfg Config) error
	Publish(topic string, payload []byte) error
	Subscribe(ctx context.Context, topic string, handler func(msg []byte) error, opts ...SubscriberOption) (func(), error)
	Close() error
}

func resolveSubscriberOption(topic string, opts ...SubscriberOption) SubscriberOption {
	if len(opts) > 0 {
		return opts[0]
	}
	return SubscriberOption{
		Mode:          ModeQueue,
		DurablePrefix: topic,
	}
}

func resolveSubscriptionID(topic, instanceID string, opts ...SubscriberOption) string {
	opt := &SubscriberOption{
		Mode:          ModeBroadcast,
		DurablePrefix: "",
	}
	if len(opts) > 0 {
		opt = &opts[0]
	}

	if opt.Mode == ModeQueue {
		if opt.DurablePrefix != "" {
			return opt.DurablePrefix
		}
		return topic + "-queue"
	}
	prefix := opt.DurablePrefix
	if prefix == "" {
		prefix = topic
	}
	id := prefix + "-" + instanceID
	safe := strings.ReplaceAll(id, ".", "-")
	return safe
}
