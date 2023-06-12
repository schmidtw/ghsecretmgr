// SPDX-FileCopyrightText: 2023 Weston Schmidt <weston_schmidt@alumni.purdue.edu>
// SPDX-License-Identifier: Apache-2.0

package ghsecretmgr_test

import (
	"context"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/go-github/v53/github"
	"github.com/schmidtw/ghsecretmgr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func isAcceptance(t *testing.T) bool {
	if strings.ToLower(os.Getenv("ACCEPTANCE_TESTS")) != "run" {
		t.Skip("Acceptance tests only run when environment variable `ACCEPTANCE_TESTS=run` is set.")
		return false
	}

	return true
}

func getTestManager(ctx context.Context) *ghsecretmgr.Manager {
	owner := os.Getenv("GH_OWNER")
	if owner == "" {
		owner = "schmidtw"
	}

	repo := os.Getenv("GH_REPO")
	if repo == "" {
		repo = "ghsecretmgr"
	}

	env := os.Getenv("GH_ENV")
	if env == "" {
		env = "testing"
	}

	return &ghsecretmgr.Manager{
		Client: github.NewTokenClient(ctx, os.Getenv("GH_TOKEN")),
		Owner:  owner,
		Repo:   repo,
		Env:    env,
	}
}

func TestAcceptance_GetSecret(t *testing.T) {
	if !isAcceptance(t) {
		return
	}

	tests := []struct {
		description string
		name        string
		clearEnv    bool
		clearRepo   bool
		expect      bool
	}{
		{
			description: "present secret test",
			name:        "present",
			clearEnv:    true,
			expect:      true,
		},
		{
			description: "non-present secret test",
			name:        "non-present",
			clearEnv:    true,
		},
		{
			description: "present env based secret test",
			name:        "env_present",
			expect:      true,
		},
		{
			description: "non-present env based secret test",
			name:        "non-present",
		},
		/*{
			description: "present org based secret test",
			name:        "org_present",
			clearEnv:    true,
			clearRepo:   true,
			expect:      true,
		},*/
		{
			description: "non-present org based secret test",
			name:        "non-present",
			clearEnv:    true,
			clearRepo:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			ctx := context.Background()
			m := getTestManager(ctx)
			if tc.clearEnv {
				m.Env = ""
			}
			if tc.clearRepo {
				m.Repo = ""
			}

			secret, resp, err := m.GetSecret(ctx, tc.name)

			assert.NoError(err)
			assert.NotNil(resp)

			if !tc.expect {
				assert.Nil(secret)
				return
			}

			require.NotNil(secret)
			assert.Equal(strings.ToLower(tc.name), strings.ToLower(secret.Name))
			assert.True(secret.CreatedAt.After(time.Time{}))
			assert.True(secret.UpdatedAt.After(time.Time{}))
		})
	}
}

func TestAcceptance_CreateAndDestroyASecret(t *testing.T) {
	if !isAcceptance(t) {
		return
	}

	var mutex sync.Mutex

	tests := []struct {
		description string
		name        string
		value       string
		clearEnv    bool
		clearRepo   bool
	}{
		{
			description: "create/destroy a repo value",
			name:        "unique",
			value:       "value",
			clearEnv:    true,
		},
		{
			description: "create/destroy an env value",
			name:        "unique",
			value:       "value",
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			mutex.Lock()
			defer mutex.Unlock()

			assert := assert.New(t)

			ctx := context.Background()
			m := getTestManager(ctx)
			if tc.clearEnv {
				m.Env = ""
			}
			if tc.clearRepo {
				m.Repo = ""
			}

			resp, err := m.UpdateSecret(ctx, tc.name, tc.value)
			assert.NoError(err)
			assert.NotNil(resp)

			resp, err = m.DeleteSecret(ctx, tc.name)
			assert.NoError(err)
			assert.NotNil(resp)
		})
	}
}
