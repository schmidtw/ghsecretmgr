// SPDX-FileCopyrightText: 2023 Weston Schmidt <weston_schmidt@alumni.purdue.edu>
// SPDX-License-Identifier: Apache-2.0

package ghsecretmgr

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/go-github/v53/github"
	"github.com/jefflinse/githubsecret"
)

var (
	ErrInvalidInput = errors.New("input is invalid")
	ErrUnable       = errors.New("unable to perform requested action")
)

type newSecretType func(context.Context, Manager) (ghAPI, *github.Response, error)

type ghAPI interface {
	GetPublicKey(context.Context, Manager) (*github.PublicKey, *github.Response, error)
	CreateOrUpdateSecret(context.Context, Manager, github.EncryptedSecret) (*github.Response, error)
	GetSecret(context.Context, Manager, string) (*github.Secret, *github.Response, error)
	DeleteSecret(context.Context, Manager, string) (*github.Response, error)
}

type Manager struct {
	Client *github.Client
	Owner  string
	Repo   string // if "" then assumed to be an owner/org secret
	Env    string // if not "" then assumed to be an environment secret

	public *github.PublicKey

	gh ghAPI
}

// GetSecret returns the secret with the name provided if possible, otherwise
// it returns a nil secret and always returns the response from github.
func (m *Manager) GetSecret(ctx context.Context, name string) (*github.Secret, *github.Response, error) {
	resp, err := m.config(ctx)
	if err != nil {
		return nil, resp, err
	}

	s, resp, err := m.gh.GetSecret(ctx, *m, name)
	if err != nil && resp.StatusCode == http.StatusNotFound {
		err = nil
	}
	return s, resp, err
}

func (m *Manager) UpdateSecret(ctx context.Context, name, value string) (*github.Response, error) {
	resp, err := m.config(ctx)
	if err != nil {
		return resp, err
	}

	if m.public == nil {
		pub, resp, err := m.gh.GetPublicKey(ctx, *m)
		if err != nil {
			return resp, err
		}
		if pub == nil || pub.GetKey() == "" || pub.GetKeyID() == "" {
			return resp, fmt.Errorf("%w: an invalid public key was returned", ErrUnable)
		}

		m.public = pub
	}

	encrypted, err := githubsecret.Encrypt(m.public.GetKey(), value)
	if err != nil {
		return nil, err
	}

	secret := github.EncryptedSecret{
		Name:           name,
		KeyID:          m.public.GetKeyID(),
		EncryptedValue: encrypted,
	}

	return m.gh.CreateOrUpdateSecret(ctx, *m, secret)
}

func (m *Manager) DeleteSecret(ctx context.Context, name string) (*github.Response, error) {
	resp, err := m.config(ctx)
	if err != nil {
		return resp, err
	}

	return m.gh.DeleteSecret(ctx, *m, name)
}

func (m *Manager) config(ctx context.Context) (*github.Response, error) {
	// Exit if the manager has already been setup.
	if m.gh != nil {
		return nil, nil
	}

	if m.Client == nil {
		return nil, fmt.Errorf("%w: Client must not be nil", ErrInvalidInput)
	}

	if m.Owner == "" {
		return nil, fmt.Errorf("%w: Owner must not be empty", ErrInvalidInput)
	}

	if m.Repo == "" && m.Env != "" {
		return nil, fmt.Errorf("%w: Env must be empty if Repo is empty", ErrInvalidInput)
	}

	var fn newSecretType
	switch {
	case m.Repo == "":
		fn = NewOrg
	case m.Env == "":
		fn = NewRepo
	default:
		fn = NewRepoEnv
	}

	tmp, resp, err := fn(ctx, *m)
	if err != nil {
		return resp, err
	}
	m.gh = tmp

	return nil, nil
}
