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

// Code generated by lister-gen. DO NOT EDIT.

package v1

import (
	v1 "github.com/blackducksoftware/synopsys-operator/pkg/api/sample-component/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// SampleComponentLister helps list SampleComponents.
type SampleComponentLister interface {
	// List lists all SampleComponents in the indexer.
	List(selector labels.Selector) (ret []*v1.SampleComponent, err error)
	// SampleComponents returns an object that can list and get SampleComponents.
	SampleComponents(namespace string) SampleComponentNamespaceLister
	SampleComponentListerExpansion
}

// sampleComponentLister implements the SampleComponentLister interface.
type sampleComponentLister struct {
	indexer cache.Indexer
}

// NewSampleComponentLister returns a new SampleComponentLister.
func NewSampleComponentLister(indexer cache.Indexer) SampleComponentLister {
	return &sampleComponentLister{indexer: indexer}
}

// List lists all SampleComponents in the indexer.
func (s *sampleComponentLister) List(selector labels.Selector) (ret []*v1.SampleComponent, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.SampleComponent))
	})
	return ret, err
}

// SampleComponents returns an object that can list and get SampleComponents.
func (s *sampleComponentLister) SampleComponents(namespace string) SampleComponentNamespaceLister {
	return sampleComponentNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// SampleComponentNamespaceLister helps list and get SampleComponents.
type SampleComponentNamespaceLister interface {
	// List lists all SampleComponents in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1.SampleComponent, err error)
	// Get retrieves the SampleComponent from the indexer for a given namespace and name.
	Get(name string) (*v1.SampleComponent, error)
	SampleComponentNamespaceListerExpansion
}

// sampleComponentNamespaceLister implements the SampleComponentNamespaceLister
// interface.
type sampleComponentNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all SampleComponents in the indexer for a given namespace.
func (s sampleComponentNamespaceLister) List(selector labels.Selector) (ret []*v1.SampleComponent, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.SampleComponent))
	})
	return ret, err
}

// Get retrieves the SampleComponent from the indexer for a given namespace and name.
func (s sampleComponentNamespaceLister) Get(name string) (*v1.SampleComponent, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("samplecomponent"), name)
	}
	return obj.(*v1.SampleComponent), nil
}
