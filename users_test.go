package dropbox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUsers_GetCurrentAccount(t *testing.T) {
	c := client()
	_, err := c.Users.GetCurrentAccount()
	assert.NoError(t, err)
}

func TestUsers_GetAccountBatch(t *testing.T) {
	c := client()
	setup, err := c.Users.GetCurrentAccount()

	assert.NoError(t, err)
	assert.NotEmpty(t, setup.AccountID, "account id required for test setup")

	out, err := c.Users.GetAccountBatch(&GetAccountBatchInput{
		AccountIDs: []string{setup.AccountID},
	})

	assert.NoError(t, err, "getting account batch")
	assert.NotEmpty(t, out, "account batch output")
	assert.EqualValues(t, setup.AccountID, out[0].AccountID, "requested match retrieved")
}
