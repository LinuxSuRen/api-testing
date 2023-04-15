package kubernetes

import "github.com/antonmedv/expr"

// PodValidatorFunc returns a expr for checking pod existing
func PodValidatorFunc() expr.Option {
	return expr.Function("pod", podValidator, new(func(...string) ResourceValidator))
}

func KubernetesValidatorFunc() expr.Option {
	return expr.Function("k8s", resourceValidator, new(func(interface{}, ...string) ResourceValidator))
}
