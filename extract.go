package formula

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"

	"github.com/expr-lang/expr/ast"
	"github.com/expr-lang/expr/parser"
)

var (
	fxOrFxfxPattern = regexp.MustCompile(`^f(\d+)(?:f(\d+))?$`)
)

// VariableRuleError 表示表达式中出现了不符合变量命名规则的变量。
//
// 规则说明：
// - 合法变量仅允许两种形态：
//  1. fx   -> 例如 f1、f20；
//  2. fxfx -> 例如 f1f2、f10f200。
//
// - x 表示 1 位或多位数字；
// - 若出现任何不符合以上两种形态的变量（如 f1f2f3、abc、x1 等），将返回该错误。
type VariableRuleError struct {
	// Name 为不符合规则的变量名。
	Name string
}

// Error 实现 error 接口，返回可直接用于日志与接口响应的错误文本。
func (e *VariableRuleError) Error() string {
	return fmt.Sprintf("变量 %q 不符合规则：仅允许 fx 或 fxfx（x 为 1 位或多位数字）", e.Name)
}

// ExtractFXVariables 解析 expr 表达式并提取所有符合 fx/fxfx 规则的变量到 Join 结构体中。
//
// 参数：
// - expression：待解析表达式。
//
// 返回值：
// - []*Join：去重后、按 Id/AssocId 升序排序的 Join 指针列表；
// - error：当表达式语法错误，或存在不符合规则的 f... 变量时返回错误。
//
// 处理细节：
// 1) 使用 parser.Parse 对表达式做语法解析，确保“先解析后提取”；
// 2) 遍历 AST，提取所有标识符；
// 3) 对每个标识符进行规则校验：
//   - 命中 ^f(\d+)(f(\d+))?$ 则提取数字并构建 Join；
//   - 不命中的任何标识符（如 f1f2f3、abc、max 等）均视为非法，返回 VariableRuleError。
func ExtractFXVariables(expression string) ([]*Join, error) {
	tree, err := parser.Parse(expression)
	if err != nil {
		return nil, fmt.Errorf("解析表达式失败: %w", err)
	}

	collector := &identifierCollector{set: map[string]struct{}{}}
	ast.Walk(&tree.Node, collector)

	joins := make([]*Join, 0, len(collector.set))
	seen := make(map[Join]struct{}, len(collector.set))
	for name := range collector.set {
		match := fxOrFxfxPattern.FindStringSubmatch(name)
		if match != nil {
			id, convErr := strconv.ParseInt(match[1], 10, 64)
			if convErr != nil {
				return nil, fmt.Errorf("解析变量 %q 的 id 失败: %w", name, convErr)
			}

			assocID := int64(0)
			if len(match) > 2 && match[2] != "" {
				parsedAssocID, assocErr := strconv.ParseInt(match[2], 10, 64)
				if assocErr != nil {
					return nil, fmt.Errorf("解析变量 %q 的 assocId 失败: %w", name, assocErr)
				}
				assocID = parsedAssocID
			}

			key := Join{Id: id, AssocId: assocID}
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}

			joins = append(joins, &Join{Id: id, AssocId: assocID})
			continue
		}

		return nil, &VariableRuleError{Name: name}
	}

	sort.Slice(joins, func(i, j int) bool {
		if joins[i].Id == joins[j].Id {
			return joins[i].AssocId < joins[j].AssocId
		}
		return joins[i].Id < joins[j].Id
	})

	return joins, nil
}

// identifierCollector 用于在 AST 遍历过程中收集标识符。
//
// 仅收集 IdentifierNode，避免把成员访问中的字段名误判为顶层变量。
// 例如 "obj.f1" 中，f1 是成员名，不应视为变量标识符。
type identifierCollector struct {
	set map[string]struct{}
}

// Visit 是 AST 访问回调，用于把命中的标识符写入 set 去重集合。
func (c *identifierCollector) Visit(node *ast.Node) {
	if node == nil || *node == nil {
		return
	}

	if ident, ok := (*node).(*ast.IdentifierNode); ok {
		c.set[ident.Value] = struct{}{}
	}
}

type Join struct {
	Id      int64
	AssocId int64
}
