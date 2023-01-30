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

// JenkinsNodeParameters are the configurable fields of a JenkinsNode.
type JenkinsNodeParameters struct {
	Name         string `json:"name"`
	NumExecutors int64  `json:"numExecutors"`
	Description  string `json:"description"`
	RemoteFS     string `json:"remoteFS"`
	Label        string `json:"label"`
}

// JenkinsNodeObservation are the observable fields of a JenkinsNode.
type JenkinsNodeObservation struct {
	Name         string `json:"name"`
	NumExecutors int64  `json:"numExecutors"`
	Description  string `json:"description"`
	RemoteFS     string `json:"remoteFS"`
	Label        string `json:"label"`
}

// A JenkinsNodeSpec defines the desired state of a JenkinsNode.
type JenkinsNodeSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       JenkinsNodeParameters `json:"forProvider"`
}

// A JenkinsNodeStatus represents the observed state of a JenkinsNode.
type JenkinsNodeStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          JenkinsNodeObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A JenkinsNode is an example API type.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,jenkins}
type JenkinsNode struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   JenkinsNodeSpec   `json:"spec"`
	Status JenkinsNodeStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// JenkinsNodeList contains a list of JenkinsNode
type JenkinsNodeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []JenkinsNode `json:"items"`
}

// JenkinsNode type metadata.
var (
	JenkinsNodeKind             = reflect.TypeOf(JenkinsNode{}).Name()
	JenkinsNodeGroupKind        = schema.GroupKind{Group: Group, Kind: JenkinsNodeKind}.String()
	JenkinsNodeKindAPIVersion   = JenkinsNodeKind + "." + SchemeGroupVersion.String()
	JenkinsNodeGroupVersionKind = SchemeGroupVersion.WithKind(JenkinsNodeKind)
)

func init() {
	SchemeBuilder.Register(&JenkinsNode{}, &JenkinsNodeList{})
}
