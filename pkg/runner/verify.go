package runner

import (
	"fmt"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"github.com/linuxsuren/api-testing/pkg/runner/kubernetes"
	"github.com/linuxsuren/api-testing/pkg/testing"
)

// Verify if the data satisfies the expression.
func Verify(expect testing.Response, data map[string]any) (err error) {
	for _, verify := range expect.Verify {
		var program *vm.Program
		if program, err = expr.Compile(verify, expr.Env(data),
			expr.AsBool(), kubernetes.PodValidatorFunc(),
			kubernetes.KubernetesValidatorFunc()); err != nil {
			return err
		}

		var result interface{}
		if result, err = expr.Run(program, data); err != nil {
			return err
		}

		if !result.(bool) {
			err = fmt.Errorf("failed to verify: %s", verify)
			fmt.Println(err)
			break
		}
	}
	return
}
