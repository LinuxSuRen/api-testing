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
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ConfigVolume  = "config"
	StorageVolume = "storage"
)

func newVolumes(atest *corev1alpha1.ATest) (vols []corev1.Volume, mounts []corev1.VolumeMount) {
	vols = []corev1.Volume{{
		Name: ConfigVolume,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: atest.Name,
				},
			},
		},
	}}
	mounts = []corev1.VolumeMount{{
		Name:      ConfigVolume,
		MountPath: "/root/.config/atest/",
	}}

	if atest.Spec.Persistent != nil && atest.Spec.Persistent.Enabled {
		vols = append(vols, corev1.Volume{
			Name: StorageVolume,
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: atest.Name,
				},
			},
		})
		mounts = append(mounts, corev1.VolumeMount{
			Name:      StorageVolume,
			MountPath: "/root/.config/storage/",
		})
	}
	return
}

func newPVC(atest *corev1alpha1.ATest) (pvc *corev1.PersistentVolumeClaim) {
	if atest.Spec.Persistent == nil || !atest.Spec.Persistent.Enabled {
		return
	}
	pvc = &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      atest.Name,
			Namespace: atest.Namespace,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteMany},
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse("1Gi"),
				},
			},
			StorageClassName: atest.Spec.Persistent.StorageClass,
		},
	}
	return
}
