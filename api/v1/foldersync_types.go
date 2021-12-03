/*
Copyright 2021 Amit Kshirsagar.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// FolderSyncSpec defines the desired state of FolderSync
type FolderSyncSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// FolderName is an example field of FolderSync. Edit foldersync_types.go to remove/update
	//+kubebuilder:object:required=true
	//+kubebuilder:validation:MinLength=6
	//+kubebuilder:validation:MaxLength=8
	//+kubebuilder:validation:Pattern="pv-*"
	FolderName string `json:"folderName,omitempty"`
	//+kubebuilder:default=5
	//+kubebuilder:validation:Optional
	//+kubebuilder:validation:Minimum=3
	//+kubebuilder:validation:Maximum=8
	SubFolderCount int32 `json:"subFolderCount,omitenty"`
}

// FolderSyncStatus defines the observed state of FolderSync
type FolderSyncStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	FolderName            string `json:"folderName,omitempty"`
	FolderNameExists      bool   `json:"folderNameExists"`
	DesiredSubFolderCount int32  `json:"subFolderCount,omitenty"`
	CurrentSubFolderCount *int32 `json:"currentFolderCount,omitenty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:shortName={"fs"}
//+kubebuilder:printcolumn:name="FolderName",type=string,JSONPath=`.status.folderName`
//+kubebuilder:printcolumn:name="FolderNameExists",type=bool,JSONPath=`.status.folderNameExists`
//+kubebuilder:printcolumn:name="Desired",type=string,JSONPath=`.spec.subFolderCount`
//+kubebuilder:printcolumn:name="Current",type=string,JSONPath=`.status.currentFolderCount`

// FolderSync is the Schema for the foldersyncs API
type FolderSync struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FolderSyncSpec   `json:"spec,omitempty"`
	Status FolderSyncStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// FolderSyncList contains a list of FolderSync
type FolderSyncList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []FolderSync `json:"items"`
}

func init() {
	SchemeBuilder.Register(&FolderSync{}, &FolderSyncList{})
}
