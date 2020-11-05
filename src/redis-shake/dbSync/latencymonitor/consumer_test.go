package latencymonitor

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKeyFilterRegexp(t *testing.T) {
	cases := map[string]bool{
		"1":    true,
		"1a":   false,
		"":     false,
		"1110": true,
	}
	for k,v := range cases{
		assert.Equal(t, v, keyPrefixRegex.MatchString(keyPrefix+k), k, keyPrefixRegex.String())
	}
}

