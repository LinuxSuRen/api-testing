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

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ATestSpec defines the desired state of ATest
type ATestSpec struct {
	Image       string             `json:"image,omitempty"`
	Version     string             `json:"version,omitempty"`
	Replicas    *int32             `json:"replicas,omitempty"`
	Persistent  *Persistent        `json:"persistent,omitempty"`
	ServiceType corev1.ServiceType `json:"serviceType,omitempty"`
}

// Persistent defines the persistent volume claim
type Persistent struct {
	Enabled      bool    `json:"enabled,omitempty"`
	StorageClass *string `json:"storageClass,omitempty"`
}

// ATestStatus defines the observed state of ATest
type ATestStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ATest is the Schema for the atests API
type ATest struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ATestSpec   `json:"spec,omitempty"`
	Status ATestStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ATestList contains a list of ATest
type ATestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ATest `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ATest{}, &ATestList{})
}
