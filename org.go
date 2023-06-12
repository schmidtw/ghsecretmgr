// SPDX-FileCopyrightText: 2023 Weston Schmidt <weston_schmidt@alumni.purdue.edu>
// SPDX-License-Identifier: Apache-2.0

package ghsecretmgr

import (
	"context"

	"github.com/google/go-github/v53/github"
)

type org struct{}

var _ ghAPI = (*org)(nil)

func NewOrg(ctx context.Context, m Manager) (ghAPI, *github.Response, error) {
	return new(org), nil, nil
}

func (org) GetPublicKey(ctx context.Context, m Manager) (*github.PublicKey, *github.Response, error) {
	return m.Client.Actions.GetOrgPublicKey(ctx, m.Owner)
}

func (org) CreateOrUpdateSecret(ctx context.Context, m Manager, secret github.EncryptedSecret) (*github.Response, error) {
	return m.Client.Actions.CreateOrUpdateOrgSecret(ctx, m.Owner, &secret)
}
func (org) GetSecret(ctx context.Context, m Manager, name string) (*github.Secret, *github.Response, error) {
	return m.Client.Actions.GetOrgSecret(ctx, m.Owner, name)
}

func (org) DeleteSecret(ctx context.Context, m Manager, name string) (*github.Response, error) {
	return m.Client.Actions.DeleteOrgSecret(ctx, m.Owner, name)
}
