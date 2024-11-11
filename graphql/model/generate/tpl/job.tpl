var (
    {{ .Name }} = "{{ .Name | camelColon }}"
    {{ .Name }}Priority = 1
)

type {{ .Name }}Payload struct{}

func New{{ .Name }}Task(payload {{ .Name }}Payload) (*asynq.Task, error) {
    data, err := json.Marshal(payload)
    if err != nil {
        return nil, err
    }
    return asynq.NewTask({{ .Name }}, data, asynq.Queue({{ .Name }})), nil
}

func Execute{{ .Name }}Task(ctx bCtx.Context, t *asynq.Task) error {
    var payload {{ .Name }}Payload
    if err := json.Unmarshal(t.Payload(), &payload); err != nil {
        return err
    }
    return nil
}

func init(){
    JobConfigMap[{{ .Name }}] = JobConfig{
        Name: {{ .Name }},
        Priority: {{ .Name }}Priority,
        Handler: Execute{{ .Name }}Task,
    }
}