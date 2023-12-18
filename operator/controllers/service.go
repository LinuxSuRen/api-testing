/*
Copyright 2023 API Testing Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	corev1alpha1 "github.com/linuxsuren/api-testing/operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func newService(atest *corev1alpha1.ATest) (service *corev1.Service) {
	service = &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      atest.Name,
			Namespace: atest.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: labelSelector,
			Type:     atest.Spec.ServiceType,
			Ports: []corev1.ServicePort{{
				Name:       "web",
				Port:       8080,
				TargetPort: intstr.FromInt(8080),
				Protocol:   corev1.ProtocolTCP,
			}},
		},
	}
	return
}
