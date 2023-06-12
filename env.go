// SPDX-FileCopyrightText: 2023 Weston Schmidt <weston_schmidt@alumni.purdue.edu>
// SPDX-License-Identifier: Apache-2.0

package ghsecretmgr

import (
	"context"
	"fmt"

	"github.com/google/go-github/v53/github"
)

type repoEnv struct {
	repoID int
}

var _ ghAPI = (*repoEnv)(nil)

func NewRepoEnv(ctx context.Context, m Manager) (ghAPI, *github.Response, error) {
	repo, resp, err := m.Client.Repositories.Get(ctx, m.Owner, m.Repo)
	if err != nil {
		return nil, resp, err
	}

	if repo.ID == nil {
		return nil, resp, fmt.Errorf("%w: github did not return a valid repo id", ErrUnable)
	}

	return &repoEnv{
		repoID: int(*repo.ID),
	}, resp, nil
}

func (r repoEnv) GetPublicKey(ctx context.Context, m Manager) (*github.PublicKey, *github.Response, error) {
	return m.Client.Actions.GetEnvPublicKey(ctx, r.repoID, m.Env)
}

func (r repoEnv) CreateOrUpdateSecret(ctx context.Context, m Manager, secret github.EncryptedSecret) (*github.Response, error) {
	return m.Client.Actions.CreateOrUpdateEnvSecret(ctx, r.repoID, m.Env, &secret)
}
func (r repoEnv) GetSecret(ctx context.Context, m Manager, name string) (*github.Secret, *github.Response, error) {
	return m.Client.Actions.GetEnvSecret(ctx, r.repoID, m.Env, name)
}

func (r repoEnv) DeleteSecret(ctx context.Context, m Manager, name string) (*github.Response, error) {
	return m.Client.Actions.DeleteEnvSecret(ctx, r.repoID, m.Env, name)
}
