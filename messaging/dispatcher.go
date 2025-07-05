package messaging

import (
	"context"

	"github.com/bytedance/sonic"

	"github.com/light-speak/lighthouse/errors"
	"github.com/light-speak/lighthouse/logs"
)

func SubscribeTyped[T any](ctx context.Context, topic string, handler func(T) error, opts ...SubscriberOption) error {
	return subscribe(ctx, topic, func(data []byte) error {
		var msg T
		if err := sonic.Unmarshal(data, &msg); err != nil {
			logs.Error().Err(err).Msg("failed to unmarshal message")
			return errors.NewBadRequestError("failed to unmarshal message", err)
		}
		if err := handler(msg); err != nil {
			logs.Error().Err(err).Msg("failed to handle message")
			return errors.NewInternalError("failed to handle message", err)
		}
		return nil
	}, resolveSubscriberOption(topic, opts...))
}

func subscribe(ctx context.Context, topic string, handler func([]byte) error, opt SubscriberOption) error {
	if broker == nil {
		return errors.NewServiceUnavailableError("broker not initialized")
	}
	_, err := broker.Subscribe(ctx, topic, handler, opt)
	return err
}

func PublishTyped[T any](ctx context.Context, topic string, msg T) error {
	if broker == nil {
		return errors.NewServiceUnavailableError("broker not initialized")
	}
	raw, err := sonic.Marshal(msg)
	if err != nil {
		logs.Error().Err(err).Msg("failed to marshal message")
		return errors.NewInternalError("failed to marshal message", err)
	}
	return broker.Publish(topic, raw)
}
