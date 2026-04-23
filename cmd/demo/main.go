package main

import (
	"fmt"
	"log"

	"formula"
)

func main() {
	engine := formula.NewEngine()

	expression := "f1 + f2 * 3"
	env := map[string]any{
		"f1": 5,
		"f2": 6,
	}

	result, err := engine.Eval(expression, env)
	if err != nil {
		log.Fatalf("执行失败: %v", err)
	}

	vars, err := formula.ExtractFXVariables("f1 + f2f3f4")
	if err != nil {
		log.Fatalf("变量提取失败: %v", err)
	}

	fmt.Printf("表达式: %s\n", expression)
	fmt.Printf("结果: %v\n", result)
	fmt.Printf("提取变量: %v\n", derefJoins(vars))
}

func derefJoins(in []*formula.Join) []formula.Join {
	out := make([]formula.Join, 0, len(in))
	for _, item := range in {
		if item == nil {
			continue
		}
		out = append(out, *item)
	}
	return out
}
