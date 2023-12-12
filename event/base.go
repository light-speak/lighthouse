package event

type LighthouseEvent interface {
	Handle(args interface{}) error
}
