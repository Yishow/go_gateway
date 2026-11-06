package workspace

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRecordingHomeSourcesLabelTrimMatchesECMAScript(t *testing.T) {
	whitespace := "\t\n\v\f\r \u00a0\u1680\u2000\u2001\u2002\u2003\u2004\u2005\u2006\u2007\u2008\u2009\u200a\u2028\u2029\u202f\u205f\u3000\ufeff"
	require.Equal(t, "Label", homeTrimLabel(whitespace+"Label"+whitespace))
	require.Empty(t, homeTrimLabel(whitespace))
	for _, preserved := range []string{"\u0085", "\u180e", "\u200b"} {
		require.Equal(t, preserved+"Label"+preserved, homeTrimLabel(preserved+"Label"+preserved))
	}
	require.Equal(t, "A\ufeffB", homeTrimLabel(" A\ufeffB "))
}
