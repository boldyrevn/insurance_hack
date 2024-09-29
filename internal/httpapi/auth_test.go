package httpapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateToken(t *testing.T) {
	login := "sasdf"
	token, err := createToken(login)
	assert.NoError(t, err)

	login, err = getLoginFromToken(token)
	assert.NoError(t, err)
	assert.Equal(t, "sasdf", login)
}
