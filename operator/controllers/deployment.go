/*
MIT License

Copyright (c) 2023 API Testing Authors.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
package controllers

import (
	"fmt"
	"strings"

	corev1alpha1 "github.com/linuxsuren/api-testing/operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func newDeployment(atest *corev1alpha1.ATest) (deploy *appsv1.Deployment) {
	deploy = &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      atest.Name,
			Namespace: atest.Namespace,
			Labels:    labelSelector,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: atest.Spec.Replicas,
			Strategy: appsv1.DeploymentStrategy{
				RollingUpdate: &appsv1.RollingUpdateDeployment{
					MaxSurge: &intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: 1,
					},
					MaxUnavailable: &intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: 0,
					},
				},
			},
			Selector: &metav1.LabelSelector{
				MatchLabels: labelSelector,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labelSelector,
				},
				Spec: corev1.PodSpec{
					Affinity: &corev1.Affinity{
						PodAntiAffinity: &corev1.PodAntiAffinity{
							PreferredDuringSchedulingIgnoredDuringExecution: []corev1.WeightedPodAffinityTerm{{
								Weight: 5,
								PodAffinityTerm: corev1.PodAffinityTerm{
									LabelSelector: &metav1.LabelSelector{
										MatchLabels: labelSelector,
									},
									TopologyKey: "kubernetes.io/hostname",
								},
							}},
						},
					},
					Containers: []corev1.Container{{
						Name:  "server",
						Image: CombineImageTag(atest.Spec.Image, atest.Spec.Version),
						LivenessProbe: &corev1.Probe{
							ProbeHandler: corev1.ProbeHandler{
								HTTPGet: &corev1.HTTPGetAction{
									Path: "/healthz",
									Port: intstr.FromInt(8080),
								},
							},
							InitialDelaySeconds: 1,
							PeriodSeconds:       5,
						},
						ReadinessProbe: &corev1.Probe{
							ProbeHandler: corev1.ProbeHandler{
								HTTPGet: &corev1.HTTPGetAction{
									Path: "/healthz",
									Port: intstr.FromInt(8080),
								},
							},
							InitialDelaySeconds: 1,
							PeriodSeconds:       5,
						},
						Resources: corev1.ResourceRequirements{
							Limits: corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse("500m"),
								corev1.ResourceMemory: resource.MustParse("512Mi"),
							},
							Requests: corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse("100m"),
								corev1.ResourceMemory: resource.MustParse("100Mi"),
							},
						},
					}},
				},
			},
		},
	}
	deploy.Spec.Template.Spec.Volumes, deploy.Spec.Template.Spec.Containers[0].VolumeMounts = newVolumes(atest)
	return
}

// CombineImageTag will return the combined image and tag in the proper format for tags and digests.
func CombineImageTag(img string, tag string) string {
	if strings.Contains(tag, ":") {
		return fmt.Sprintf("%s@%s", img, tag) // Digest
	} else if len(tag) > 0 {
		return fmt.Sprintf("%s:%s", img, tag) // Tag
	}
	return img // No tag, use default
}

var labelSelector = map[string]string{
	"app.kubernetes.io/name": "atest",
}
