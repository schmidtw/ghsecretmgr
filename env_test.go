// SPDX-FileCopyrightText: 2023 Weston Schmidt <weston_schmidt@alumni.purdue.edu>
// SPDX-License-Identifier: Apache-2.0

package ghsecretmgr

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/google/go-github/v53/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/stretchr/testify/assert"
)

func TestNewRepoEnv(t *testing.T) {
	unknownErr := errors.New("unknown")
	tests := []struct {
		description string
		mockErr     bool
		mockID      int64
		expectErr   error
	}{
		{
			description: "working test",
			mockID:      10,
		}, {
			description: "missing id test",
			expectErr:   ErrUnable,
		}, {
			description: "error returned",
			mockErr:     true,
			expectErr:   unknownErr,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			assert := assert.New(t)

			var mc *http.Client
			if !tc.mockErr {
				var repo github.Repository

				if tc.mockID > 0 {
					repo.ID = &tc.mockID
				}

				mc = mock.NewMockedHTTPClient(
					mock.WithRequestMatch(
						mock.GetReposByOwnerByRepo,
						repo,
					),
				)
			} else {
				mc = mock.NewMockedHTTPClient(
					mock.WithRequestMatchHandler(
						mock.GetRepositories,
						http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							mock.WriteError(
								w,
								http.StatusInternalServerError,
								"github went belly up or something",
							)
						}),
					),
				)
			}
			m := Manager{
				Client: github.NewClient(mc),
				Owner:  "owner",
				Repo:   "repo",
				Env:    "env",
			}

			ctx := context.Background()
			got, _, err := NewRepoEnv(ctx, m)

			if tc.expectErr != nil {
				if errors.Is(tc.expectErr, unknownErr) {
					assert.Error(err)
				} else {
					assert.ErrorIs(err, tc.expectErr)
				}
				assert.Nil(got)
				return
			}

			assert.NoError(err)
			assert.NotNil(got)
		})
	}
}
