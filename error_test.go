package dropbox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	c := client()
	out, err := c.Files.GetMetadata(ctx, &GetMetadataInput{
		Path: "/this/does/not/exist.txt",
	})

	assert.Nil(t, out, "no object should be returned on error")
	assert.Error(t, err, "error should occur for nonexistant path")

	derr, ok := err.(*Error)

	assert.True(t, ok, "error should be a package error")
	assert.NotNil(t, derr.Header, "error headers should be present")
	assert.NotNil(t, derr.Err, "error detail should be present")
	tag, value := derr.Tag()
	assert.Equal(t, "path", tag, "error should indicate the path was invalid")
	assert.Equal(t, "not_found", value, "error should indicate not found")
}
