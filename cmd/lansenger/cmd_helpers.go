package main

import "strings"

// 通用命令行参数辅助函数。
//
// 这几个函数与具体业务域无关，集中放在此处供各 cmd_*.go 共用
// （原先就近定义在部分业务域命令文件中，那些文件已随应用域一并移除）。
//
// 依赖的 parseJSONRaw / checkError 定义在 main.go 中。

// intPtrOrNil 把可选的 int 参数转成 *int，负数视为「未提供」。
func intPtrOrNil(v int) *int {
	if v < 0 {
		return nil
	}
	val := v
	return &val
}

// int64PtrOrNil 把可选的 int64 参数转成 *int64，负数视为「未提供」。
func int64PtrOrNil(v int64) *int64 {
	if v < 0 {
		return nil
	}
	val := v
	return &val
}

// splitCommaList 按逗号切分命令行传入的列表参数。
func splitCommaList(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, ",")
}

// parseJSONListFlag 解析形如 `[{"a":1},{"b":2}]` 的 JSON 数组参数。
// flagName 保留用于与调用处签名一致。
func parseJSONListFlag(raw, flagName string) []map[string]interface{} {
	if raw == "" {
		return nil
	}
	decoded, err := parseJSONRaw(raw)
	checkError(err)
	list, ok := decoded.([]interface{})
	if !ok {
		return nil
	}
	out := make([]map[string]interface{}, 0, len(list))
	for _, item := range list {
		if m, ok := item.(map[string]interface{}); ok {
			out = append(out, m)
		}
	}
	return out
}
