var (
    {{ .Name }} = "{{ .Name | camelColon }}"
    {{ .Name }}Priority = 1
)

type {{ .Name }}TaskExecutor struct {}

type {{ .Name }}Payload struct{}

func {{ .Name }}Options() []asynq.Option {
	return []asynq.Option{
		asynq.Queue({{ .Name }}), // 指定 channel
		asynq.MaxRetry(15), // 重试次数
		asynq.ProcessIn(time.Second * 2), // 延迟执行时间
		asynq.Retention(24 * time.Hour), // 执行完成保留时间
	}
}

func New{{ .Name }}Task(payload {{ .Name }}Payload) (*asynq.Task, error) {
    data, err := sonic.Marshal(payload)
    if err != nil {
        return nil, err
    }
    return asynq.NewTask({{ .Name }}, data, {{ .Name }}Options()...), nil
}

func (ex *{{ .Name }}TaskExecutor) Execute(ctx context.Context, t *asynq.Task) error {
    var payload {{ .Name }}Payload
    if err := sonic.Unmarshal(t.Payload(), &payload); err != nil {
        return err
    }
    return nil
}

func init(){
    executor := &{{ .Name }}TaskExecutor{}
    lightqueue.JobConfigMap[{{ .Name }}] = lightqueue.JobConfig{
        Name: {{ .Name }},
        Priority: {{ .Name }}Priority,
        Executor: executor,
    }
}