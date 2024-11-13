var (
	{{.Name}} = "{{.Name | camelColon}}"
)

type {{.Name}}Value struct {}

func Publish{{.Name}}Message(value {{.Name}}Value) error {
	jsonVal, err := json.Marshal(value)
	if err != nil {
		return err
	}
	msg := message.NewMessage(uuid.New().String(), jsonVal)
	msg.Metadata.Set(EventTypeHeader, {{.Name}})
	return pub.Publish({{.Name}}, msg)
}
