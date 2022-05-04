package v1

import (
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VolumeGroup is a specification for a VolumeGroup resource
type VolumeGroup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VolumeGroupSpec   `json:"spec"`
	Status VolumeGroupStatus `json:"status"`
}

// VolumeGroupSpec is the spec for a VolumeGroup resource
type VolumeGroupSpec struct {
	NodeName string            `json:"nodeName"`
	VgName   string            `json:"vgName"`
	Size     resource.Quantity `json:"size"`
}

// VolumeGroupStatus is the status for a VolumeGroup resource
type VolumeGroupStatus struct {
	PvcStatus map[string]Item `json:"pvcStatus"`
	PvStatus  map[string]Item `json:"pvStatus"`
}

type Item struct {
	PvName      string            `json:"pvName"`
	CurrentSize resource.Quantity `json:"currentSize"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VolumeGroupList is a list of VolumeGroup resources
type VolumeGroupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []VolumeGroup `json:"items"`
}

func init() {
	SchemeBuilder.Register(&VolumeGroup{}, &VolumeGroupList{})
}
