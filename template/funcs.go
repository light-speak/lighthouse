package template

import (
	"fmt"
)

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
