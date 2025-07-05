package messaging

type SubscribeMode string

const (
	ModeBroadcast SubscribeMode = "broadcast"
	ModeQueue     SubscribeMode = "queue"
)

type SubscriberOption struct {
	Mode          SubscribeMode
	DurablePrefix string
}
