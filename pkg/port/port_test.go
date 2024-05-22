package port

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_ParsePort_ParsesPort(t *testing.T) {
	port, err := ParsePort("80")
	require.NoError(t, err)
	require.Equal(t, Port(80), port)
}

func Test_ParsePort_WhenPortInvalid_ReturnsError(t *testing.T) {
	_, err := ParsePort("invalid")
	require.ErrorContains(t, err, "invalid syntax")
}
