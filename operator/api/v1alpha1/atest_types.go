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
