package queue

import (
	"embed"
	"os"
	"path/filepath"

	"github.com/light-speak/lighthouse/templates"
	"github.com/light-speak/lighthouse/utils"
)

//go:embed tpl
var tpl embed.FS

func InitQueue() error {
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	queuestartTpl, err := tpl.ReadFile("tpl/queuestart.tpl")
	if err != nil {
		return err
	}
	templates.AddImportRegex("cmd", "github.com/light-speak/lighthouse/cmd", "")
	templates.AddImportRegex("queue", "github.com/light-speak/lighthouse/queue", "")
	options := &templates.Options{
		Path:         filepath.Join(currentDir, "commands"),
		Template:     string(queuestartTpl),
		FileName:     "queue-start",
		FileExt:      "go",
		Editable:     false,
		SkipIfExists: false,
	}
	return templates.Render(options)
}

func GenTask(name string) error {
	taskTemplate, err := tpl.ReadFile("tpl/task.tpl")
	if err != nil {
		return err
	}
	curPath, err := os.Getwd()
	if err != nil {
		return err
	}
	templates.AddImportRegex("asynq", "github.com/hibiken/asynq", "")
	templates.AddImportRegex("sonic", "github.com/bytedance/sonic", "")
	templates.AddImportRegex("context", "context", "")
	templates.AddImportRegex("queue", "github.com/light-speak/lighthouse/queue", "")
	templates.AddImportRegex("time", "time", "")
	options := &templates.Options{
		Path:         filepath.Join(curPath, "queue"),
		Template:     string(taskTemplate),
		FileName:     utils.SnakeCase(name),
		FileExt:      "go",
		Package:      "queue",
		Editable:     true,
		SkipIfExists: false,
		Funcs: map[string]any{
			"camelColon": utils.CamelColon,
		},
		Data: map[string]interface{}{
			"Name": name,
		},
	}
	return templates.Render(options)
}
