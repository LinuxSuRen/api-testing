package kubernetes

import "github.com/antonmedv/expr"

// PodValidatorFunc returns a expr for checking pod existing
func PodValidatorFunc() expr.Option {
	return expr.Function("pod", podValidator, new(func(...string) ResourceValidator))
}

// KubernetesValidatorFunc returns a expr for checking the generic Kubernetes resources
func KubernetesValidatorFunc() expr.Option {
	return expr.Function("k8s", resourceValidator, new(func(interface{}, ...string) ResourceValidator))
}
