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

// NodeParameters are the configurable fields of a Node.
type NodeParameters struct {
	Name         string `json:"name"`
	numExecutors int    `json:"numExecutors"`
	description  string `json:"description"`
	remoteFS     string `json:"remoteFS"`
	label        string `json:"label"`
}

// NodeObservation are the observable fields of a Node.
type NodeObservation struct {
	Name string `json:"name"`
}

// A NodeSpec defines the desired state of a Node.
type NodeSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       NodeParameters `json:"forProvider"`
}

// A NodeStatus represents the observed state of a Node.
type NodeStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          NodeObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A Node is an example API type.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,provider-jenkins}
type Node struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NodeSpec   `json:"spec"`
	Status NodeStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// NodeList contains a list of Node
type NodeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Node `json:"items"`
}

// Node type metadata.
var (
	NodeKind             = reflect.TypeOf(Node{}).Name()
	NodeGroupKind        = schema.GroupKind{Group: Group, Kind: NodeKind}.String()
	NodeKindAPIVersion   = NodeKind + "." + SchemeGroupVersion.String()
	NodeGroupVersionKind = SchemeGroupVersion.WithKind(NodeKind)
)

func init() {
	SchemeBuilder.Register(&Node{}, &NodeList{})
}
