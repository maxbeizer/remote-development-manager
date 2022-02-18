package config

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

var testJson = `
{
	"commands": {
		"reltest": {
			"executablePath": "./foo"
		},

		"abstest": {
			"executablePath": "/tmp/bar",
			"longRunning": true 
		}
	}
}
`

func TestLoad(t *testing.T) {
	file, err := ioutil.TempFile("", "config.json")
	require.NoError(t, err)

	file.WriteString(testJson)
	file.Close()

	config := New()
	err = config.Load(file.Name())

	require.NoError(t, err)

	require.NotNil(t, config.Commands["reltest"])
	expectedRelTestPath := filepath.Join(filepath.Dir(file.Name()), "foo")
	require.Equal(t, expectedRelTestPath, config.Commands["reltest"].ExecutablePath)
	require.False(t, config.Commands["reltest"].LongRunning)

	require.NotNil(t, config.Commands["abstest"])
	expectedAbsTestPath := "/tmp/bar"
	require.Equal(t, expectedAbsTestPath, config.Commands["abstest"].ExecutablePath)
	require.True(t, config.Commands["abstest"].LongRunning)

	require.NotNil(t, config.Commands["abstest"])
}
