package protocol_test

import (
	"goTibia/internal/protocol"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Encode_Decode(t *testing.T) {
	ip, err := protocol.StringToIP("127.0.0.1")
	require.NoError(t, err)

	ipString := protocol.IPToString(ip)

	require.Equal(t, "127.0.0.1", ipString)
}
