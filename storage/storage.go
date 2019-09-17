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
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	corelisters "k8s.io/client-go/listers/core/v1"

	ddp "github.com/AmitKumarDas/storage-provisioner/pkg/apis/ddp/v1alpha1"
)

const (
	storageProviderKey string = "storageprovider.ddp.mayadata.io/storageclass-name"
)

func ptrTrue() *bool {
	t := true
	return &t
}

func findProvider(dict map[string]string) (string, bool) {
	if len(dict) == 0 {
		return "", false
	}
	provider, found := dict[storageProviderKey]
	return provider, found
}

func isStorageOwner(owners []metav1.OwnerReference, storage *ddp.Storage) bool {
	for _, o := range owners {
		if o.APIVersion != storage.APIVersion &&
			o.Kind != storage.Kind &&
			o.Name != storage.Name &&
			o.UID != storage.UID {
			return true
		}
	}
	return false
}

// Reconciler manages reconciling storage API
// in kubernetes cluster
type Reconciler struct {
	// instances to invoke various Kubernetes APIs
	Clientset kubernetes.Interface
	PVCLister corelisters.PersistentVolumeClaimLister

	// name of the storage provider
	provider string
}

func (r *Reconciler) String() string {
	return "StorageReconciler"
}

// Reconcile accepts storage as the desired state and starts executing
// the reconcile logic based on this desired state
//
// NOTE:
//	Reconcile logic needs to be idempotent
func (r *Reconciler) Reconcile(stor *ddp.Storage) error {
	var found bool
	if r.provider, found = findProvider(stor.GetAnnotations()); !found {
		return errors.Errorf(
			"%s %s/%s: Reconcile failed: Missing storage provider",
			r, stor.Namespace, stor.Name,
		)
	}

	// find if PVC is created in previous reconcile attempt
	pvc, err := r.findPVC(stor)
	if err != nil {
		return err
	}

	// create PVC if not found
	if pvc == nil {
		return r.createPVC(stor)
	}

	// update PVC if desired state was changed
	return r.updatePVC(pvc, stor)
}

// findPVC will list & find the correct PVC if available
func (r *Reconciler) findPVC(stor *ddp.Storage) (*v1.PersistentVolumeClaim, error) {
	list, err :=
		r.PVCLister.PersistentVolumeClaims(stor.Namespace).List(labels.Everything())
	if err != nil {
		return nil, err
	}

	for _, pvc := range list {
		isowner := isStorageOwner(pvc.OwnerReferences, stor)
		if isowner {
			return pvc, nil
		}
	}
	return nil, nil
}

func (r *Reconciler) updatePVC(pvc *v1.PersistentVolumeClaim, stor *ddp.Storage) error {
	if pvc.Spec.Resources.Requests[v1.ResourceStorage] == stor.Spec.Capacity {
		// no changes
		return nil
	}

	copy := pvc.DeepCopy()
	copy.Spec.Resources.Requests[v1.ResourceStorage] = stor.Spec.Capacity
	_, err :=
		r.Clientset.CoreV1().PersistentVolumeClaims(stor.Namespace).Update(copy)
	return err
}

func (r *Reconciler) createPVC(stor *ddp.Storage) error {
	pvc := r.buildPVCFromStorage(stor)
	_, err :=
		r.Clientset.CoreV1().PersistentVolumeClaims(stor.Namespace).Create(pvc)
	return err
}

func (r *Reconciler) buildPVCFromStorage(stor *ddp.Storage) *v1.PersistentVolumeClaim {
	return &v1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: stor.Name,
			Namespace:    stor.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				metav1.OwnerReference{
					APIVersion:         stor.APIVersion,
					Kind:               stor.Kind,
					Name:               stor.Name,
					UID:                stor.UID,
					Controller:         ptrTrue(),
					BlockOwnerDeletion: ptrTrue(),
				},
			},
		},
		Spec: v1.PersistentVolumeClaimSpec{
			Resources: v1.ResourceRequirements{
				Requests: map[v1.ResourceName]resource.Quantity{
					v1.ResourceStorage: stor.Spec.Capacity,
				},
			},
			StorageClassName: &r.provider,
		},
	}
}
