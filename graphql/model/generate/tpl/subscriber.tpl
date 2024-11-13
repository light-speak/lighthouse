var (
	{{.Name}} = "{{.Name | camelColon}}"
)

type {{.Name}}Payload struct {}


func Execute{{.Name}}(msg []byte) error {
	var payload {{.Name}}Payload
	if err := json.Unmarshal(msg, &payload); err != nil {
		return err
	}
	// start code here
	return nil
}

func init() {
	subscriberExecuteMapping[{{.Name}}] = Execute{{.Name}}
}
