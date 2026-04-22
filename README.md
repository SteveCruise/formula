# formula

一个基于 `github.com/expr-lang/expr` 的 Go 表达式执行与变量提取示例项目。

## 功能

- 提供 API：执行符合 expr 规则的表达式并返回结果。
- 提供 API：解析表达式并提取符合 `fx/fxfx` 规则的变量。
  - `x` 为 1 位或多位数字。
  - 合法示例：`f1`、`f20`、`f1f2`、`f10f200`。
  - 非法示例：`f1f2f3`（即 `fxfxfx`，会报错）。

## 安装依赖

```bash
go mod tidy
```

## API 说明

### 1) 执行表达式

```go
engine := formula.NewEngine()
result, err := engine.Eval("f1 + f2 * 2", map[string]any{
    "f1": 3,
    "f2": 4,
})
```

### 2) 提取并校验变量（返回 Join）

```go
vars, err := formula.ExtractFXVariables("f1 + f2f3 + max(a, 10)")
// vars == []*formula.Join{{Id:1, AssocId:0}, {Id:2, AssocId:3}}
```

映射规则：
- `f1` -> `Join{Id:1, AssocId:0}`
- `f2f3` -> `Join{Id:2, AssocId:3}`

如果表达式里出现 `f1f2f3` 这种 `fxfxfx` 形式，会返回 `VariableRuleError`。

## 运行 demo

```bash
go run cmd/demo/main.go
```

## 运行测试

```bash
go test ./...
```

## 目录结构

```text
.
├── cmd/demo/main.go
├── engine.go
├── example_test.go
├── extract.go
├── formula_test.go
├── go.mod
└── README.md
```
