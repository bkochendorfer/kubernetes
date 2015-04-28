/*
Copyright 2015 Google Inc. All rights reserved.

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

package secret

import (
	"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/api/rest"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/fields"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/labels"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/watch"
)

// Registry is an interface implemented by things that know how to store Secret objects.
type Registry interface {
	// ListSecrets obtains a list of Secrets having labels which match selector.
	ListSecrets(ctx api.Context, selector labels.Selector) (*api.SecretList, error)
	// Watch for new/changed/deleted secrets
	WatchSecrets(ctx api.Context, label labels.Selector, field fields.Selector, resourceVersion string) (watch.Interface, error)
	// Get a specific Secret
	GetSecret(ctx api.Context, name string) (*api.Secret, error)
	// Create a Secret based on a specification.
	CreateSecret(ctx api.Context, Secret *api.Secret) (*api.Secret, error)
	// Update an existing Secret
	UpdateSecret(ctx api.Context, Secret *api.Secret) (*api.Secret, error)
	// Delete an existing Secret
	DeleteSecret(ctx api.Context, name string) error
}

// storage puts strong typing around storage calls
type storage struct {
	rest.StandardStorage
}

// NewRegistry returns a new Registry interface for the given Storage. Any mismatched
// types will panic.
func NewRegistry(s rest.StandardStorage) Registry {
	return &storage{s}
}

func (s *storage) ListSecrets(ctx api.Context, label labels.Selector) (*api.SecretList, error) {
	obj, err := s.List(ctx, label, fields.Everything())
	if err != nil {
		return nil, err
	}
	return obj.(*api.SecretList), nil
}

func (s *storage) WatchSecrets(ctx api.Context, label labels.Selector, field fields.Selector, resourceVersion string) (watch.Interface, error) {
	return s.Watch(ctx, label, field, resourceVersion)
}

func (s *storage) GetSecret(ctx api.Context, name string) (*api.Secret, error) {
	obj, err := s.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	return obj.(*api.Secret), nil
}

func (s *storage) CreateSecret(ctx api.Context, secret *api.Secret) (*api.Secret, error) {
	obj, err := s.Create(ctx, secret)
	return obj.(*api.Secret), err
}

func (s *storage) UpdateSecret(ctx api.Context, secret *api.Secret) (*api.Secret, error) {
	obj, _, err := s.Update(ctx, secret)
	return obj.(*api.Secret), err
}

func (s *storage) DeleteSecret(ctx api.Context, name string) error {
	_, err := s.Delete(ctx, name, nil)
	return err
}
