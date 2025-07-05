package templates

import (
	"fmt"
	"text/template"

	"github.com/light-speak/lighthouse/utils"
)

// addFunc adds template functions to the template
func (o *Options) addFunc() {
	if o.Funcs == nil {
		o.Funcs = template.FuncMap{}
	}
	o.Funcs["lc"] = utils.Lc
	o.Funcs["lcFirst"] = utils.LcFirst
	o.Funcs["ucFirst"] = utils.UcFirst
	o.Funcs["snakeCase"] = utils.SnakeCase
	o.Funcs["camelCase"] = utils.CamelCase
	o.Funcs["camelCaseWithSpecial"] = utils.CamelCaseWithSpecial
	o.Funcs["userCodeStart"] = UserCodeStart
	o.Funcs["userCodeEnd"] = UserCodeEnd
	o.Funcs["userCodeSection"] = UserCodeSection
}

func UserCodeStart(action string) string {
	return fmt.Sprintf("// Func:%s user code start. Do not remove this comment. ", action)
}

func UserCodeEnd(action string) string {
	return fmt.Sprintf("// Func:%s user code end. Do not remove this comment. ", action)
}

func UserCodeSection() string {
	return fmt.Sprint("// Section: user code section start. Do not remove this comment. \n",
		"// Section: user code section end. Do not remove this comment. \n")
}

func Add(a, b int) int { return a + b }
func Sub(a, b int) int { return a - b }
func Mul(a, b int) int { return a * b }
func Div(a, b int) int { return a / b }
