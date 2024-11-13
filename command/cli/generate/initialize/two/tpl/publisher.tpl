var (
	pub *kafka.Publisher
)

const (
	EventTypeHeader = "event_type"
)

func GetPublisher() *kafka.Publisher {
	return pub
}

func init() {
	if !env.LighthouseConfig.Kafka.Enable {
		return
	}

	p, err := kafka.NewPublisher(kafka.PublisherConfig{
		Brokers:   env.LighthouseConfig.Kafka.Brokers,
		Marshaler: kafka.DefaultMarshaler{},
	}, watermill.NewStdLogger(false, false))

	if err != nil {
		log.Fatal().Err(err).Msg("failed to create kafka publisher")
	}
	pub = p
}
