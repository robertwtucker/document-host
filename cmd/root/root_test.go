package root_test

import (
	"testing"

	subject "github.com/robertwtucker/document-host/cmd/root"
	"github.com/stretchr/testify/assert"
)

func TestRootCmd(t *testing.T) {
	err := subject.Cmd().Execute()
	assert.NoError(t, err)
}
