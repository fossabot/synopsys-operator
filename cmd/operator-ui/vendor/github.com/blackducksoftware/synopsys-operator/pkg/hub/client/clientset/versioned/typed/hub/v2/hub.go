/*
Copyright The Kubernetes Authors.

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

package v2

import (
	v2 "github.com/blackducksoftware/synopsys-operator/pkg/api/hub/v2"
	scheme "github.com/blackducksoftware/synopsys-operator/pkg/hub/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// HubsGetter has a method to return a HubInterface.
// A group's client should implement this interface.
type HubsGetter interface {
	Hubs(namespace string) HubInterface
}

// HubInterface has methods to work with Hub resources.
type HubInterface interface {
	Create(*v2.Hub) (*v2.Hub, error)
	Update(*v2.Hub) (*v2.Hub, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v2.Hub, error)
	List(opts v1.ListOptions) (*v2.HubList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v2.Hub, err error)
	HubExpansion
}

// hubs implements HubInterface
type hubs struct {
	client rest.Interface
	ns     string
}

// newHubs returns a Hubs
func newHubs(c *SynopsysV2Client, namespace string) *hubs {
	return &hubs{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the hub, and returns the corresponding hub object, and an error if there is any.
func (c *hubs) Get(name string, options v1.GetOptions) (result *v2.Hub, err error) {
	result = &v2.Hub{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("hubs").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Hubs that match those selectors.
func (c *hubs) List(opts v1.ListOptions) (result *v2.HubList, err error) {
	result = &v2.HubList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("hubs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested hubs.
func (c *hubs) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("hubs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a hub and creates it.  Returns the server's representation of the hub, and an error, if there is any.
func (c *hubs) Create(hub *v2.Hub) (result *v2.Hub, err error) {
	result = &v2.Hub{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("hubs").
		Body(hub).
		Do().
		Into(result)
	return
}

// Update takes the representation of a hub and updates it. Returns the server's representation of the hub, and an error, if there is any.
func (c *hubs) Update(hub *v2.Hub) (result *v2.Hub, err error) {
	result = &v2.Hub{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("hubs").
		Name(hub.Name).
		Body(hub).
		Do().
		Into(result)
	return
}

// Delete takes name of the hub and deletes it. Returns an error if one occurs.
func (c *hubs) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("hubs").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *hubs) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("hubs").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched hub.
func (c *hubs) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v2.Hub, err error) {
	result = &v2.Hub{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("hubs").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
