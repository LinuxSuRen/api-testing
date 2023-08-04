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
