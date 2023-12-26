/*
Copyright 2023.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	StatusPending  = "PENDING"
	StatusRunning  = "RUNNING"
	StatusCleaning = "CLEANING"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type Task struct {
	Type                string `json:"Type,omitempty"`
	Cpu                 string `json:"cpu,omitempty"`
	Memory              string `json:"memory,omitempty"`
	Gpu                 string `json:"gpu,omitempty"`
	TaskName            string `json:"taskName,omitempty"`
	ContainerReplicas   int32  `json:"containerReplicas,omitempty"`
	ContainerRegistry   string `json:"containerRegistry,omitempty"`
	ContainerImage      string `json:"containerImage,omitempty"`
	ContainerTag        string `json:"containerTag,omitempty"`
	ContainerEntrypoint string `json:"containerEntrypoint,omitempty"`
}

// SlamlSpec defines the desired state of Slaml
type SlamlSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Slaml. Edit slaml_types.go to remove/update
	// Foo                 string `json:"foo,omitempty"`
	// ClientId    string `json:"clientId,omitempty"`
	IsSla       string `json:"IsSla,omitempty"`
	Name        string `json:"name,omitempty"`
	SlaTarget   int32  `json:"slaTarget,omitempty"`
	VolcanoKind string `json:"volcanoKind,omitempty"`
	Tasks       []Task `json:"tasks"`
}

// SlamlStatus defines the observed state of Slaml
type SlamlStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	ClientStatus string `json:"clientStatus,omitempty"`
	LastPodName  string `json:"lastPodName,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Slaml is the Schema for the slamls API
type Slaml struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SlamlSpec   `json:"spec,omitempty"`
	Status SlamlStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SlamlList contains a list of Slaml
type SlamlList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Slaml `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Slaml{}, &SlamlList{})
}
