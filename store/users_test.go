package store_test

import (
	"async-api/fixtures"
	"async-api/store"
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUserStore(t *testing.T) {
	testEnv := fixtures.NewTestEnv(t)
	cleanup := testEnv.SetupDb(t)
	t.Cleanup(func() {
		cleanup(t)
	})

	userStore := store.NewUserStore(testEnv.Db)
	user, err := userStore.CreateUser(context.Background(), "test@test.com", "testingpassword")
	require.NoError(t, err)
	require.NotNil(t, user)
	require.Equal(t, "test@test.com", user.Email)
	require.NoError(t, user.ComparePassword("testingpassword"))

	userId := user.Id
	user, err = userStore.ById(context.Background(), userId)
	require.NoError(t, err)
	require.NotNil(t, user)

	user, err = userStore.ByEmail(context.Background(), "test@test.com")
	require.NoError(t, err)
	require.NotNil(t, user)
}
