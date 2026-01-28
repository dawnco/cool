package fn

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
)

// isSkipFrame 判断是否跳过该栈帧
func isSkipFrame(file string) bool {
	skipPatterns := []string{
		"/runtime/",
		"/src/runtime/",
		"/go/src/",           // Go 标准库（Linux/macOS 默认安装路径）
		"/usr/local/go/src/", // macOS / Linux 常见 Go 安装路径
		"/net/http/",
		"/internal/fn/stack_trace.go",
		"/middleware/panic.go",
		"_test.go",
		"/testing/",
		// 可根据项目结构调整，例如：
		// "/vendor/",
	}
	for _, p := range skipPatterns {
		if strings.Contains(file, p) {
			return true
		}
	}
	return false
}

// StackTrace 获取过滤后的调用栈：仅保留业务代码的文件+行号
func StackTrace() string {
	var buf bytes.Buffer
	const maxDepth = 50 // 防止无限循环

	for i := 1; i < maxDepth; i++ { // 从 1 开始跳过 recover 自身
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		if isSkipFrame(file) {
			continue
		}

		// 只保留相对路径或简化路径（可选）
		// 例如：去掉 GOPATH 前缀，或只保留项目内路径
		// 这里直接记录完整路径，你也可以用 filepath.Base(file) 只留文件名

		buf.WriteString(fmt.Sprintf("%s:%d\n", file, line))
	}

	stack := buf.String()
	if stack == "" {
		return "no business stack trace found"
	}
	return strings.TrimSuffix(stack, " ")
}
