package tracer

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEbpfBytesCorrect(t *testing.T) {
	bs, err := ioutil.ReadFile("../../ebpf/tcptracer-ebpf.o")
	require.NoError(t, err)

	actual, err := tcptracerEbpfOBytes()
	require.NoError(t, err)

	assert.Equal(t, bs, actual)
}
