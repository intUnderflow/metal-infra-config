package bootstrap

import (
	"encoding/json"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_GetBootstrapRecords_WhenInvalidJSON_ReturnsError(t *testing.T) {
	fs := memfs.New()
	fileToWrite, err := fs.Create("test.json")
	require.NoError(t, err)
	_, err = fileToWrite.Write([]byte("invalid"))
	require.NoError(t, err)
	require.NoError(t, fileToWrite.Close())

	fileToRead, err := fs.Open("test.json")
	require.NoError(t, err)
	_, err = GetBootstrapRecords(fileToRead)
	require.ErrorContains(t, err, "invalid character")
}

func Test_GetBootstrapRecords_ReturnsCorrectValues(t *testing.T) {
	testData := map[string]string{
		"foo": "bar",
		"baz": "982348234789",
	}
	testDataBytes, err := json.Marshal(testData)
	require.NoError(t, err)
	fs := memfs.New()
	fileToWrite, err := fs.Create("test.json")
	require.NoError(t, err)
	_, err = fileToWrite.Write(testDataBytes)
	require.NoError(t, err)
	require.NoError(t, fileToWrite.Close())

	fileToRead, err := fs.Open("test.json")
	require.NoError(t, err)
	result, err := GetBootstrapRecords(fileToRead)
	require.NoError(t, err)
	require.EqualValues(t, testData, result)
}
