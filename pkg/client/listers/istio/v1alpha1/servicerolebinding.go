// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/nais/naiserator/pkg/apis/istio/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// ServiceRoleBindingLister helps list ServiceRoleBindings.
type ServiceRoleBindingLister interface {
	// List lists all ServiceRoleBindings in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.ServiceRoleBinding, err error)
	// ServiceRoleBindings returns an object that can list and get ServiceRoleBindings.
	ServiceRoleBindings(namespace string) ServiceRoleBindingNamespaceLister
	ServiceRoleBindingListerExpansion
}

// serviceRoleBindingLister implements the ServiceRoleBindingLister interface.
type serviceRoleBindingLister struct {
	indexer cache.Indexer
}

// NewServiceRoleBindingLister returns a new ServiceRoleBindingLister.
func NewServiceRoleBindingLister(indexer cache.Indexer) ServiceRoleBindingLister {
	return &serviceRoleBindingLister{indexer: indexer}
}

// List lists all ServiceRoleBindings in the indexer.
func (s *serviceRoleBindingLister) List(selector labels.Selector) (ret []*v1alpha1.ServiceRoleBinding, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.ServiceRoleBinding))
	})
	return ret, err
}

// ServiceRoleBindings returns an object that can list and get ServiceRoleBindings.
func (s *serviceRoleBindingLister) ServiceRoleBindings(namespace string) ServiceRoleBindingNamespaceLister {
	return serviceRoleBindingNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// ServiceRoleBindingNamespaceLister helps list and get ServiceRoleBindings.
type ServiceRoleBindingNamespaceLister interface {
	// List lists all ServiceRoleBindings in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1alpha1.ServiceRoleBinding, err error)
	// Get retrieves the ServiceRoleBinding from the indexer for a given namespace and name.
	Get(name string) (*v1alpha1.ServiceRoleBinding, error)
	ServiceRoleBindingNamespaceListerExpansion
}

// serviceRoleBindingNamespaceLister implements the ServiceRoleBindingNamespaceLister
// interface.
type serviceRoleBindingNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all ServiceRoleBindings in the indexer for a given namespace.
func (s serviceRoleBindingNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.ServiceRoleBinding, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.ServiceRoleBinding))
	})
	return ret, err
}

// Get retrieves the ServiceRoleBinding from the indexer for a given namespace and name.
func (s serviceRoleBindingNamespaceLister) Get(name string) (*v1alpha1.ServiceRoleBinding, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("servicerolebinding"), name)
	}
	return obj.(*v1alpha1.ServiceRoleBinding), nil
}
