package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PracticeGroupSpec defines the desired state of PracticeGroup
type PracticeGroupSpec struct {
	// Name in human-readable form.
	Name string `json:"name"`
	// The year that is practice group started.
	StartYear int `json:"startYear,omitempty"`
	// The leader of the practice group.
	Leader Employee `json:"leader"`
	// Members of the practice group.
	Members []Employee `json:"members"`
}

// PracticeGroupStatus defines the observed state of PracticeGroup
type PracticeGroupStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

type Employee struct {
	// The full name of the employee.
	Name string `json:"name"`
	// The unique ID assigned to each employee, as seen in internal systems.
	EmployeeId uint32 `json:"employeeId"`
	// The seniority of the employee
	//+kubebuilder:validation:Enum=Consultant;Senior Consultant;Manager;Principal
	Seniority string `json:"seniority"`
	// The department that this practice group belongs to.
	//+kubebuilder:validation:Enum=Technology;Design;Operations;Trondheim
	Department string `json:"department"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName="pg"
// PracticeGroup is the Schema for the practicegroups API
type PracticeGroup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PracticeGroupSpec   `json:"spec,omitempty"`
	Status PracticeGroupStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PracticeGroupList contains a list of PracticeGroup
type PracticeGroupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PracticeGroup `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PracticeGroup{}, &PracticeGroupList{})
}
