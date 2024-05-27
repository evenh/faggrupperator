//go:build !ignore_autogenerated

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Employee) DeepCopyInto(out *Employee) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Employee.
func (in *Employee) DeepCopy() *Employee {
	if in == nil {
		return nil
	}
	out := new(Employee)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PracticeGroup) DeepCopyInto(out *PracticeGroup) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PracticeGroup.
func (in *PracticeGroup) DeepCopy() *PracticeGroup {
	if in == nil {
		return nil
	}
	out := new(PracticeGroup)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PracticeGroup) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PracticeGroupList) DeepCopyInto(out *PracticeGroupList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]PracticeGroup, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PracticeGroupList.
func (in *PracticeGroupList) DeepCopy() *PracticeGroupList {
	if in == nil {
		return nil
	}
	out := new(PracticeGroupList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PracticeGroupList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PracticeGroupSpec) DeepCopyInto(out *PracticeGroupSpec) {
	*out = *in
	out.Leader = in.Leader
	if in.Members != nil {
		in, out := &in.Members, &out.Members
		*out = make([]Employee, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PracticeGroupSpec.
func (in *PracticeGroupSpec) DeepCopy() *PracticeGroupSpec {
	if in == nil {
		return nil
	}
	out := new(PracticeGroupSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PracticeGroupStatus) DeepCopyInto(out *PracticeGroupStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PracticeGroupStatus.
func (in *PracticeGroupStatus) DeepCopy() *PracticeGroupStatus {
	if in == nil {
		return nil
	}
	out := new(PracticeGroupStatus)
	in.DeepCopyInto(out)
	return out
}
