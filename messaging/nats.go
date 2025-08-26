package messaging

import (
	"context"
	"time"

	"github.com/light-speak/lighthouse/lighterr"
	"github.com/light-speak/lighthouse/logs"
	"github.com/nats-io/nats.go"
)

type NatsBroker struct {
	conn       *nats.Conn
	js         nats.JetStream
	instanceID string
}

func (n *NatsBroker) Init(cfg Config) error {
	nc, err := nats.Connect(cfg.URL,
		nats.MaxReconnects(-1), nats.ReconnectWait(2*time.Second))
	if err != nil {
		return lighterr.NewServiceUnavailableError("failed to connect to nats", err)
	}
	js, err := nc.JetStream()
	if err != nil {
		return lighterr.NewServiceUnavailableError("failed to create jetstream client", err)
	}
	streamName := "messaging"
	streamSubjects := []string{"messaging.>"}
	if _, err := js.StreamInfo(streamName); err == nats.ErrStreamNotFound {
		if _, err := js.AddStream(&nats.StreamConfig{
			Name:      streamName,
			Subjects:  streamSubjects,
			Storage:   nats.FileStorage,
			Retention: nats.LimitsPolicy,
		}); err != nil {
			return lighterr.NewServiceUnavailableError("failed to add stream", err)
		}
	}
	n.conn = nc
	n.js = js
	n.instanceID = cfg.InstanceID
	return nil
}

func fullSubject(topic string) string {
	return "messaging." + topic
}

func (n *NatsBroker) Publish(topic string, payload []byte) error {
	_, err := n.js.Publish(fullSubject(topic), payload)
	return err
}

func (n *NatsBroker) Subscribe(ctx context.Context, topic string, handler func(msg []byte) error, opts ...SubscriberOption) (func(), error) {
	subId := resolveSubscriptionID(fullSubject(topic), n.instanceID, opts...)
	var sub *nats.Subscription
	var err error
	opt := &SubscriberOption{
		Mode:          ModeBroadcast,
		DurablePrefix: "",
	}
	if len(opts) > 0 {
		opt = &opts[0]
	}

	// Create a goroutine that monitors the context cancellation
	unsubscribeCh := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			// Context was cancelled, unsubscribe from the topic
			if sub != nil {
				_ = sub.Unsubscribe()
			}
			close(unsubscribeCh)
		case <-unsubscribeCh:
			// Unsubscribe was called manually, just return
			return
		}
	}()
	switch opt.Mode {
	case ModeQueue:
		sub, err = n.js.QueueSubscribe(
			fullSubject(topic),
			subId,
			wrapHandler(ctx, handler),
			nats.Durable(subId),
			nats.ManualAck(),
		)
	case ModeBroadcast:
		sub, err = n.js.Subscribe(
			fullSubject(topic),
			wrapHandler(ctx, handler),
			nats.DeliverNew(),
			nats.Durable(subId),
			nats.ManualAck(),
		)
	}
	if err != nil {
		logs.Error().Err(err).Msg("failed to subscribe to topic")
		close(unsubscribeCh)
		return nil, lighterr.NewServiceUnavailableError("failed to subscribe to topic", err)
	}

	return func() {
		logs.Debug().Msg("unsubscribing from topic")
		_ = sub.Unsubscribe()
		close(unsubscribeCh)
	}, nil
}

func wrapHandler(ctx context.Context, handler func([]byte) error) nats.MsgHandler {
	return func(m *nats.Msg) {
		if ctx.Err() != nil {
			logs.Debug().Msg("context cancelled")
			return
		}
		defer func() {
			if r := recover(); r != nil {
				logs.Error().Err(lighterr.NewInternalError("panic in message handler")).Msg("panic in message handler")
			}
		}()
		err := handler(m.Data)
		if err != nil {
			logs.Error().Err(err).Msg("failed to handle message")
			return
		}
		m.Ack()
	}
}

func (n *NatsBroker) Close() error {
	if n.conn != nil && !n.conn.IsClosed() {
		n.conn.Close()
	}
	return nil
}
