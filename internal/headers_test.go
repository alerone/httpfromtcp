package headers

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidSingleHeader(t *testing.T) {
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["Host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)
}

func TestInvalidSpacingHeader(t *testing.T) {
	headers := NewHeaders()
	data := []byte("       Host   : localhost:42069     \r\n\r\n")
	n, done, err := headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}

func TestValidTwoHeadersWithExistingHeaders(t *testing.T) {
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\nAuthorization: 1dd90c72-daea-44c1-bc29-75f18fc4522b\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["Host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	copy(data, data[n:])
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	fmt.Println(len("Host: localhost:42069\r\n"))
	fmt.Println(len("Authorization: 1dd90c72-daea-44c1-bc29-75f18fc4522b\r\n"))
	assert.Equal(t, "1dd90c72-daea-44c1-bc29-75f18fc4522b", headers["Authorization"])
	assert.Equal(t, 53, n)
	assert.False(t, done)

	//Test valid done
	copy(data, data[n:])
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.True(t, done)
}


