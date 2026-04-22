package formula

import (
	"fmt"

	"github.com/expr-lang/expr"
)

// Engine 是表达式执行引擎。
//
// 该结构体对 github.com/expr-lang/expr 进行轻量封装，
// 目的是向上层业务提供稳定、清晰的 API：
// 1) 接收一段符合 expr 语法规则的表达式；
// 2) 接收表达式执行所需的变量环境（env）；
// 3) 返回表达式计算结果或可读的错误信息。
//
// 说明：
// - Engine 当前不维护内部状态，多个 goroutine 可安全复用同一个实例；
// - 保留结构体是为了未来扩展（如缓存编译结果、注入函数白名单等）。
type Engine struct{}

// NewEngine 创建一个新的表达式执行引擎实例。
//
// 返回值：
// - *Engine：可直接用于调用 Eval 执行表达式。
func NewEngine() *Engine {
	return &Engine{}
}

// Eval 执行一段符合 expr 规则的表达式，并返回计算结果。
//
// 参数：
// - expression：待执行的表达式字符串，例如 "f1 + f2 * 10"；
// - env：表达式变量环境，key 为变量名，value 为变量值。
//
// 返回值：
// - any：表达式执行结果（具体类型由表达式运行时决定）；
// - error：当表达式编译失败、执行失败或参数非法时返回错误。
//
// 执行流程：
// 1) 使用 expr.Compile 进行语法与语义编译；
// 2) 使用 expr.Run 在给定 env 下执行程序；
// 3) 将底层错误包装为带上下文的信息，便于调用方排查。
func (e *Engine) Eval(expression string, env map[string]any) (any, error) {
	if expression == "" {
		return nil, fmt.Errorf("expression 不能为空")
	}

	program, err := expr.Compile(expression, expr.Env(env))
	if err != nil {
		return nil, fmt.Errorf("编译表达式失败: %w", err)
	}

	result, err := expr.Run(program, env)
	if err != nil {
		return nil, fmt.Errorf("执行表达式失败: %w", err)
	}

	return result, nil
}
