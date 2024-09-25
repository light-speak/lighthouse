package test

import (
	"testing"

	"github.com/light-speak/lighthouse/log"
)

func TestLog(t *testing.T) {
	// 测试日志输出
	log.Info("这是一条信息日志")
	log.Warn("这是一条警告日志")
	log.Error("这是一条错误日志")

	// 测试带参数的日志输出
	log.Info("用户 %s 登录成功", "张三")
	log.Warn("数据库连接数达到 %d，接近上限", 80)
	log.Error("操作失败，错误码：%d，错误信息：%s", 500, "内部服务器错误")

	t.Log("测试日志输出完成")
	// 注意：这里我们没有实际验证日志输出，因为日志通常是写入文件或控制台
	// 在实际项目中，你可能需要使用 mock 或者捕获输出来进行更严格的测试
}
