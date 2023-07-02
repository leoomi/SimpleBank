package api

import (
	"testing"

	db "github.com/leoomi/simplebank/db/sqlc"
	"github.com/leoomi/simplebank/util"
	"github.com/stretchr/testify/require"
)

func randomUser(t *testing.T) (user db.User, password string) {
	password = util.RandomString(6)
	hashedPassword, err := util.HashedPassword(password)
	require.NoError(t, err)

	user = db.User{
		Username: util.RandomOwner(),
		Password: hashedPassword,
		FullName: util.RandomOwner(),
		Email:    util.RandomEmail(),
	}
	return
}
