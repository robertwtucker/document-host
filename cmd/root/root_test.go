package root

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRootCmd(t *testing.T) {
	err := Cmd().Execute()
	assert.NoError(t, err)
}
