// Code generated by lister-gen. DO NOT EDIT.

package v1alpha2

import (
	v1alpha2 "github.com/nais/naiserator/pkg/apis/storage.cnrm.cloud.google.com/v1alpha2"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// StorageBucketAccessControlLister helps list StorageBucketAccessControls.
type StorageBucketAccessControlLister interface {
	// List lists all StorageBucketAccessControls in the indexer.
	List(selector labels.Selector) (ret []*v1alpha2.StorageBucketAccessControl, err error)
	// StorageBucketAccessControls returns an object that can list and get StorageBucketAccessControls.
	StorageBucketAccessControls(namespace string) StorageBucketAccessControlNamespaceLister
	StorageBucketAccessControlListerExpansion
}

// storageBucketAccessControlLister implements the StorageBucketAccessControlLister interface.
type storageBucketAccessControlLister struct {
	indexer cache.Indexer
}

// NewStorageBucketAccessControlLister returns a new StorageBucketAccessControlLister.
func NewStorageBucketAccessControlLister(indexer cache.Indexer) StorageBucketAccessControlLister {
	return &storageBucketAccessControlLister{indexer: indexer}
}

// List lists all StorageBucketAccessControls in the indexer.
func (s *storageBucketAccessControlLister) List(selector labels.Selector) (ret []*v1alpha2.StorageBucketAccessControl, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha2.StorageBucketAccessControl))
	})
	return ret, err
}

// StorageBucketAccessControls returns an object that can list and get StorageBucketAccessControls.
func (s *storageBucketAccessControlLister) StorageBucketAccessControls(namespace string) StorageBucketAccessControlNamespaceLister {
	return storageBucketAccessControlNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// StorageBucketAccessControlNamespaceLister helps list and get StorageBucketAccessControls.
type StorageBucketAccessControlNamespaceLister interface {
	// List lists all StorageBucketAccessControls in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1alpha2.StorageBucketAccessControl, err error)
	// Get retrieves the StorageBucketAccessControl from the indexer for a given namespace and name.
	Get(name string) (*v1alpha2.StorageBucketAccessControl, error)
	StorageBucketAccessControlNamespaceListerExpansion
}

// storageBucketAccessControlNamespaceLister implements the StorageBucketAccessControlNamespaceLister
// interface.
type storageBucketAccessControlNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all StorageBucketAccessControls in the indexer for a given namespace.
func (s storageBucketAccessControlNamespaceLister) List(selector labels.Selector) (ret []*v1alpha2.StorageBucketAccessControl, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha2.StorageBucketAccessControl))
	})
	return ret, err
}

// Get retrieves the StorageBucketAccessControl from the indexer for a given namespace and name.
func (s storageBucketAccessControlNamespaceLister) Get(name string) (*v1alpha2.StorageBucketAccessControl, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha2.Resource("storagebucketaccesscontrol"), name)
	}
	return obj.(*v1alpha2.StorageBucketAccessControl), nil
}
