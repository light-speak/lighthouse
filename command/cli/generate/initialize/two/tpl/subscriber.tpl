var (
	sub               		 *kafka.Subscriber
	EventTypeHeader          = "event_type"
	subscriberExecuteMapping = map[string]func(msg []byte) error{}
	Topics                   = []string{""} // edit by yourself
	GroupID                  = ""        // edit by yourself
)

func GetSubscriber() *kafka.Subscriber {
	return sub
}

func Start() {
	defer sub.Close()

	for _, topic := range Topics {
		go func(topic string) {
			for {
				msg, err := sub.Subscribe(bCtx.Background(), topic)
				if err != nil {
					log.Error().Err(err).Msg("failed to subscribe topic")
					return
				}
				for msg := range msg {
					eventType := msg.Metadata.Get(EventTypeHeader)
					fn := subscriberExecuteMapping[eventType]
					if fn == nil {
						log.Info().Msgf("no execute function for event type: %s", eventType)
						return
					}
					go func(eventType string) {
						err := fn(msg.Payload)
						if err != nil {
							log.Error().Err(err).Msg("failed to execute message")
						}
						msg.Ack()

					}(eventType)
				}
			}
		}(topic)
	}

	select {}
}

func init() {
	if !env.LighthouseConfig.Kafka.Enable {
		return
	}

	s, err := kafka.NewSubscriber(kafka.SubscriberConfig{
		Brokers:       env.LighthouseConfig.Kafka.Brokers,
		Unmarshaler:   kafka.DefaultMarshaler{},
		ConsumerGroup: GroupID,
	}, watermill.NewStdLogger(false, false))

	if err != nil {
		log.Fatal().Err(err).Msg("failed to create kafka subscriber")
	}
	sub = s
}
