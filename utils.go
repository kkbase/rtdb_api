package rtdb_api

import "C"
import (
	"regexp"
	"strings"
)

// CFunc 把原始的C的函数转换成解析后的C函数
func CFunc(code string) string {
	// 1. 提取注释块
	commentRe := regexp.MustCompile(`(?s)/\*\*.*?\*/`)
	comment := commentRe.FindString(code)

	// 2. 去掉注释，只保留函数声明部分
	codeNoComment := commentRe.ReplaceAllString(code, "")

	// 3. 把多行压成一行，便于解析
	normalized := strings.Join(strings.Fields(codeNoComment), " ")

	// 示例结果：
	// RTDBAPI rtdb_error RTDBAPI_CALLRULE rtdb_get_own_connection( rtdb_int32 handle, rtdb_int32 node_number, rtdb_int32* socket );

	// 4. 解析函数声明
	funcRe := regexp.MustCompile(
		`RTDBAPI\s+(\w+)\s+RTDBAPI_CALLRULE\s+(\w+)\s*\((.*?)\)\s*;`,
	)
	matches := funcRe.FindStringSubmatch(normalized)
	if len(matches) != 4 {
		return "// parse error: invalid function signature"
	}

	retType := matches[1]
	funcName := matches[2]
	params := strings.TrimSpace(matches[3])

	// 5. wrapper 名称
	wrapName := funcName + "_warp"
	typedefName := funcName + "_fn"

	// 6. 生成 wrapper
	var sb strings.Builder

	if comment != "" {
		sb.WriteString(comment)
		sb.WriteString("\n")
	}

	sb.WriteString(retType)
	sb.WriteString(" RTDBAPI_CALLRULE ")
	sb.WriteString(wrapName)
	sb.WriteString("(")
	sb.WriteString(params)
	sb.WriteString(")\n{\n")

	sb.WriteString("    typedef ")
	sb.WriteString(retType)
	sb.WriteString(" (RTDBAPI_CALLRULE *")
	sb.WriteString(typedefName)
	sb.WriteString(")(")
	sb.WriteString(params)
	sb.WriteString(");\n")

	sb.WriteString("    ")
	sb.WriteString(typedefName)
	sb.WriteString(" fn = (")
	sb.WriteString(typedefName)
	sb.WriteString(")get_function(\"")
	sb.WriteString(funcName)
	sb.WriteString("\");\n")

	sb.WriteString("    return fn(")
	sb.WriteString(extractParamNames(params))
	sb.WriteString(");\n")

	sb.WriteString("}\n")

	return sb.String()
}

func extractParamNames(params string) string {
	if strings.TrimSpace(params) == "" {
		return ""
	}

	parts := strings.Split(params, ",")
	names := make([]string, 0, len(parts))

	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "void" {
			continue
		}
		fields := strings.Fields(p)
		if len(fields) == 0 {
			continue
		}
		// 参数名永远在最后（支持 rtdb_int32* socket 这种）
		names = append(names, fields[len(fields)-1])
	}

	return strings.Join(names, ", ")
}

func sp(s string) []string {
	// 统一换行符（兼容 Windows / Unix）
	s = strings.ReplaceAll(s, "\r\n", "\n")

	// 正则：一个或多个“空白行”
	re := regexp.MustCompile(`\n\s*\n+`)

	parts := re.Split(strings.TrimSpace(s), -1)

	// 可选：过滤掉意外的空段
	res := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			res = append(res, p)
		}
	}
	return res
}

// XXX 将 C 接口注释 + 函数定义，转换为 Go 的空函数定义
func XXX(s string) string {
	comment := extractComment(s)
	cFuncName := extractCFuncName(s)
	goFuncName := "Raw" + toGoFuncName(cFuncName)

	var b strings.Builder
	if comment != "" {
		b.WriteString(comment)
		b.WriteString("\n")
	}
	b.WriteString("func ")
	b.WriteString(goFuncName)
	b.WriteString("() {}\n")

	return b.String()
}

// 提取 /** ... */ 注释
func extractComment(s string) string {
	re := regexp.MustCompile(`(?s)/\*\*.*?\*/`)
	return re.FindString(s)
}

// 提取 C 函数名
func extractCFuncName(s string) string {
	re := regexp.MustCompile(`\b([a-zA-Z_][a-zA-Z0-9_]*)\s*\(`)
	m := re.FindStringSubmatch(s)
	if len(m) > 1 {
		return m[1]
	}
	return ""
}

// rtdbh_update_value64_warp -> RawRtdbhUpdateValue64Warp
func toGoFuncName(cName string) string {
	parts := strings.Split(cName, "_")
	for i, p := range parts {
		if p == "" {
			continue
		}
		parts[i] = strings.ToUpper(p[:1]) + p[1:]
	}
	return strings.Join(parts, "")
}
