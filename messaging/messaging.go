package messaging

import "log"

var broker Broker

func init() {
	var b Broker

	switch cfg.Driver {
	case DriverNats:
		b = &NatsBroker{}
	case DriverKafka:
		log.Fatal("Kafka is not implemented")
	}
	if err := b.Init(cfg); err != nil {
		log.Fatal("failed to initialize broker", err)
	}
	broker = b
}
