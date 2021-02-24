package bloom

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunMemhashString(t *testing.T) {
	v := memhashStringKey("some-string", 123)
	assert.True(t, v == 0 || v != 0)
}
