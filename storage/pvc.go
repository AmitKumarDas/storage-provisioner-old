/*
Copyright 2019 The MayaData Authors.

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

package storage

import (
	"fmt"

	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	storage "k8s.io/api/storage/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	storagelisters "k8s.io/client-go/listers/storage/v1beta1"
	ref "k8s.io/client-go/tools/reference"
	"k8s.io/klog"
)

// PVCReconciler manages reconciling PVC API
// in kubernetes cluster
type PVCReconciler struct {
	// instances to invoke various Kubernetes APIs
	Clientset kubernetes.Interface
	VALister  storagelisters.VolumeAttachmentLister

	// pvc object that will be reconciled
	pvc *v1.PersistentVolumeClaim

	// node where the storage should get attached
	nodeName string

	// name of the attacher that should attach the volume
	attacherName string
}

// String implements stringer interface
func (r *PVCReconciler) String() string {
	if r.pvc == nil {
		return "PVCReconciler"
	}
	return fmt.Sprintf("PVCReconciler %s/%s", r.pvc.Namespace, r.pvc.Name)
}

// Reconcile accepts PVC as the desired state and starts executing
// the reconcile logic based on this desired state
//
// NOTE:
//	Reconcile logic needs to be idempotent
func (r *PVCReconciler) Reconcile(pvc *v1.PersistentVolumeClaim) error {
	r.pvc = pvc

	if pvc.Spec.VolumeName == "" {
		// nothing to do since PVC is not yet bound to any PV
		klog.V(3).Infof(
			"%s: Reconcile ignored: Volume not bound", r,
		)
		return nil
	}

	// find if VolumeAttachment is created in previous reconcile attempt
	va, err := r.findVA()
	if err != nil {
		return err
	}

	// create VolumeAttachment if not found
	if va == nil {
		return r.createVA()
	}

	// update VolumeAttachment if desired state was changed
	update, err := r.updateVA(va)
	if !update {
		klog.V(3).Infof("%s: No change to desired state", r)
	}
	return err
}

// findVA will list & find the correct VolumeAttachment if available
func (r *PVCReconciler) findVA() (*storage.VolumeAttachment, error) {

	list, err := r.VALister.List(labels.Everything())
	if err != nil {
		return nil, err
	}

	for _, va := range list {
		isowner := isPVCOwner(va.OwnerReferences, r.pvc)
		if isowner {
			return va, nil
		}
	}
	return nil, nil
}

// updateVA updates the given VolumeAttachment in case of any change
// in the desired state
func (r *PVCReconciler) updateVA(va *storage.VolumeAttachment) (bool, error) {

	nodeName, found := findNodeNameFromPVC(r.pvc)
	if !found || nodeName == "" {
		// nothing to verify & evaluate further
		return false, nil
	}

	if nodeName == va.Spec.NodeName {
		// no change
		return false, nil
	}

	// we shall delete the VolumeAttachment & expect a new one
	// to get created as part of next reconcile invocation
	return true, r.Clientset.StorageV1beta1().VolumeAttachments().
		Delete(va.Name, &metav1.DeleteOptions{})
}

func (r *PVCReconciler) createVA() error {
	var found bool
	r.nodeName, found = findNodeNameFromPVC(r.pvc)
	if !found {
		return errors.Errorf(
			"%s: Reconcile failed: Node name not found", r,
		)
	}

	r.attacherName, found = findAttacherFromPVC(r.pvc)
	if !found {
		return errors.Errorf(
			"%s: Reconcile failed: Attacher name not found", r,
		)
	}

	va, err := r.newVA()
	if err != nil {
		return err
	}

	_, err =
		r.Clientset.StorageV1beta1().VolumeAttachments().Create(va)
	return err
}

func (r *PVCReconciler) newVA() (*storage.VolumeAttachment, error) {
	pvcref, err := ref.GetReference(scheme.Scheme, r.pvc)
	if err != nil {
		return nil, err
	}

	return &storage.VolumeAttachment{
		ObjectMeta: metav1.ObjectMeta{
			// pvc name is supposed to be generated in this case
			Name: pvcref.Name,
			OwnerReferences: []metav1.OwnerReference{
				metav1.OwnerReference{
					APIVersion:         pvcref.APIVersion,
					Kind:               pvcref.Kind,
					Name:               pvcref.Name,
					UID:                pvcref.UID,
					Controller:         boolPtr(true),
					BlockOwnerDeletion: boolPtr(true),
				},
			},
		},
		Spec: storage.VolumeAttachmentSpec{
			Source: storage.VolumeAttachmentSource{
				PersistentVolumeName: strPtr(r.pvc.Spec.VolumeName),
			},
			NodeName: r.nodeName,
			Attacher: r.attacherName,
		},
	}, nil
}
