package password

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var plainPassword string
var hashPassword string

func init() {
	plainPassword = "Hello World"
}

func TestHasPassword(t *testing.T) {
	hashPassword = Hash(plainPassword)

	assert.NotEqual(t, plainPassword, hashPassword, "Password Should Be Hashed")
}

func TestComparePasswordSuccess(t *testing.T) {
	compareResult := Compare(plainPassword, hashPassword)

	assert.Equal(t, true, compareResult, "Password Comparing Should Be Success")
}

func TestCompareWrongPassword(t *testing.T) {
	compareResult := Compare("World, Hello", hashPassword)

	assert.Equal(t, false, compareResult, "Password Comparing Should Be Success")
}
