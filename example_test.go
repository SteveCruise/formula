package formula_test

import (
	"fmt"
	"log"

	"formula"
)

func Example() {
	engine := formula.NewEngine()

	result, err := engine.Eval("f1 + f2 * 10", map[string]any{
		"f1": 2,
		"f2": 3,
	})
	if err != nil {
		log.Fatal(err)
	}

	vars, err := formula.ExtractFXVariables("f1 + f2f3f4")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)
	fmt.Println(derefJoins(vars))
	// Output:
	// 32
	// [{1 0 0} {2 3 4}]
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
