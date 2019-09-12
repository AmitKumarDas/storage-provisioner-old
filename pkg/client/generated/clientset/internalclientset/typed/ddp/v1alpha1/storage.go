/*
Copyright 2019 The MayaData Authors

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

// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"time"

	v1alpha1 "github.com/AmitKumarDas/storage-provisioner/pkg/apis/ddp/v1alpha1"
	scheme "github.com/AmitKumarDas/storage-provisioner/pkg/client/generated/clientset/internalclientset/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// StoragesGetter has a method to return a StorageInterface.
// A group's client should implement this interface.
type StoragesGetter interface {
	Storages(namespace string) StorageInterface
}

// StorageInterface has methods to work with Storage resources.
type StorageInterface interface {
	Create(*v1alpha1.Storage) (*v1alpha1.Storage, error)
	Update(*v1alpha1.Storage) (*v1alpha1.Storage, error)
	UpdateStatus(*v1alpha1.Storage) (*v1alpha1.Storage, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.Storage, error)
	List(opts v1.ListOptions) (*v1alpha1.StorageList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Storage, err error)
	StorageExpansion
}

// storages implements StorageInterface
type storages struct {
	client rest.Interface
	ns     string
}

// newStorages returns a Storages
func newStorages(c *DdpV1alpha1Client, namespace string) *storages {
	return &storages{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the storage, and returns the corresponding storage object, and an error if there is any.
func (c *storages) Get(name string, options v1.GetOptions) (result *v1alpha1.Storage, err error) {
	result = &v1alpha1.Storage{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("storages").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Storages that match those selectors.
func (c *storages) List(opts v1.ListOptions) (result *v1alpha1.StorageList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.StorageList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("storages").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested storages.
func (c *storages) Watch(opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("storages").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a storage and creates it.  Returns the server's representation of the storage, and an error, if there is any.
func (c *storages) Create(storage *v1alpha1.Storage) (result *v1alpha1.Storage, err error) {
	result = &v1alpha1.Storage{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("storages").
		Body(storage).
		Do().
		Into(result)
	return
}

// Update takes the representation of a storage and updates it. Returns the server's representation of the storage, and an error, if there is any.
func (c *storages) Update(storage *v1alpha1.Storage) (result *v1alpha1.Storage, err error) {
	result = &v1alpha1.Storage{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("storages").
		Name(storage.Name).
		Body(storage).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *storages) UpdateStatus(storage *v1alpha1.Storage) (result *v1alpha1.Storage, err error) {
	result = &v1alpha1.Storage{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("storages").
		Name(storage.Name).
		SubResource("status").
		Body(storage).
		Do().
		Into(result)
	return
}

// Delete takes name of the storage and deletes it. Returns an error if one occurs.
func (c *storages) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("storages").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *storages) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("storages").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched storage.
func (c *storages) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Storage, err error) {
	result = &v1alpha1.Storage{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("storages").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
