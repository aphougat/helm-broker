// Code generated by mockery v1.0.0
package automock

import internal "github.com/kyma-project/helm-broker/internal"
import mock "github.com/stretchr/testify/mock"

// DocsProvider is an autogenerated mock type for the DocsProvider type
type DocsProvider struct {
	mock.Mock
}

// EnsureDocsTopic provides a mock function with given fields: addon
func (_m *DocsProvider) EnsureDocsTopic(addon *internal.Addon) error {
	ret := _m.Called(addon)

	var r0 error
	if rf, ok := ret.Get(0).(func(*internal.Addon) error); ok {
		r0 = rf(addon)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// EnsureDocsTopicRemoved provides a mock function with given fields: id
func (_m *DocsProvider) EnsureDocsTopicRemoved(id string) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetNamespace provides a mock function with given fields: namespace
func (_m *DocsProvider) SetNamespace(namespace string) {
	_m.Called(namespace)
}
