// SPDX-FileCopyrightText: 2023 Weston Schmidt <weston_schmidt@alumni.purdue.edu>
// SPDX-License-Identifier: Apache-2.0

package ghsecretmgr

import (
	"context"

	"github.com/google/go-github/v53/github"
)

type repo struct{}

var _ ghAPI = (*repo)(nil)

func NewRepo(ctx context.Context, m Manager) (ghAPI, *github.Response, error) {
	return new(repo), nil, nil
}

func (repo) GetPublicKey(ctx context.Context, m Manager) (*github.PublicKey, *github.Response, error) {
	return m.Client.Actions.GetRepoPublicKey(ctx, m.Owner, m.Repo)
}

func (repo) CreateOrUpdateSecret(ctx context.Context, m Manager, secret github.EncryptedSecret) (*github.Response, error) {
	return m.Client.Actions.CreateOrUpdateRepoSecret(ctx, m.Owner, m.Repo, &secret)
}
func (repo) GetSecret(ctx context.Context, m Manager, name string) (*github.Secret, *github.Response, error) {
	return m.Client.Actions.GetRepoSecret(ctx, m.Owner, m.Repo, name)
}

func (repo) DeleteSecret(ctx context.Context, m Manager, name string) (*github.Response, error) {
	return m.Client.Actions.DeleteRepoSecret(ctx, m.Owner, m.Repo, name)
}
