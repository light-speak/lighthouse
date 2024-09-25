package utils

import (
	"fmt"
	"os"
	"strings"
)

func AddImport(file, pkg string) error {
	content, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("读取文件失败: %v", err)
	}

	if strings.Contains(string(content), fmt.Sprintf(`"%s"`, pkg)) {
		return nil
	}

	newContent := strings.Replace(string(content), "import (", fmt.Sprintf("import (\n\t\"%s\"", pkg), 1)
	err = os.WriteFile(file, []byte(newContent), 0644)
	if err != nil {
		return fmt.Errorf("写入文件失败: %v", err)
	}

	return nil
}
