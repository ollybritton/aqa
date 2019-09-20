package token

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestLookupIdent tests that valid keywords will be associated with their token.Type.
func TestLookupIdent(t *testing.T) {
	assert.Equal(t, LookupIdent("constant"), Type(CONSTANT), "'constant' should be translated to the CONSTANT keyword, not a %v", LookupIdent("constant"))
	assert.Equal(t, LookupIdent("CONSTANT"), Type(CONSTANT), "'CONSTANT' should be translated to the CONSTANT keyword, not a %v", LookupIdent("CONSTANT"))
	assert.Equal(t, LookupIdent("cOnStAnT"), Type(IDENT), "'cOnStAnT' should be an ident, not a %v", LookupIdent("cOnStAnT"))
}
