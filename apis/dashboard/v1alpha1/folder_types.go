/*
Copyright 2022 The Crossplane Authors.

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
	"reflect"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
)

// FolderParameters are the configurable fields of a Folder.
type FolderParameters struct {
	ConfigurableField string `json:"configurableField"`
}

// FolderObservation are the observable fields of a Folder.
type FolderObservation struct {
	ObservableField string `json:"observableField,omitempty"`
}

// A FolderSpec defines the desired state of a Folder.
type FolderSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       FolderParameters `json:"forProvider"`
}

// A FolderStatus represents the observed state of a Folder.
type FolderStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          FolderObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A Folder is an example API type.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,jenkins}
type Folder struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FolderSpec   `json:"spec"`
	Status FolderStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// FolderList contains a list of Folder
type FolderList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Folder `json:"items"`
}

// Folder type metadata.
var (
	FolderKind             = reflect.TypeOf(Folder{}).Name()
	FolderGroupKind        = schema.GroupKind{Group: Group, Kind: FolderKind}.String()
	FolderKindAPIVersion   = FolderKind + "." + SchemeGroupVersion.String()
	FolderGroupVersionKind = SchemeGroupVersion.WithKind(FolderKind)
)

func init() {
	SchemeBuilder.Register(&Folder{}, &FolderList{})
}
