package entities

import "github.com/stretchr/testify/require"
import "testing"

func Test_Config_GetWithVersion_WithValueSet_ReturnsValue(t *testing.T) {
	config := NewConfig()

	require.NoError(t, config.SetWithVersion("foo", "bar", 1))
	value, err := config.GetWithVersion("foo")
	require.NoError(t, err)
	require.Equal(t, value.Value, "bar")
}

func Test_Config_GetWithVersion_WithNoValueSet_ReturnsError(t *testing.T) {
	config := NewConfig()

	_, err := config.GetWithVersion("foo")
	require.EqualError(t, err, _errValueDoesNotExist)
}

func Test_Config_SetWithVersion_WhenNewValue_Succeeds(t *testing.T) {
	config := NewConfig()

	require.NoError(t, config.SetWithVersion("foo", "bar", 1))
	value, err := config.GetWithVersion("foo")
	require.NoError(t, err)
	require.Equal(t, value.Value, "bar")
}

func Test_Config_SetWithVersion_WhenExistingValue_WithHigherVersion_Succeeds(t *testing.T) {
	config := NewConfig()

	require.NoError(t, config.SetWithVersion("foo", "zip", 1))
	require.NoError(t, config.SetWithVersion("foo", "bar", 2))

	value, err := config.GetWithVersion("foo")
	require.NoError(t, err)
	require.Equal(t, value.Value, "bar")
}

func Test_Config_SetWithVersion_WhenValueIsOlder_ReturnsError(t *testing.T) {
	config := NewConfig()

	require.NoError(t, config.SetWithVersion("foo", "bar", 2))
	require.EqualError(t, config.SetWithVersion("foo", "bar", 1), ErrNewValueIsOlderThanCurrentValue)
}

func Test_Config_List_ReturnsKeysWithValues(t *testing.T) {
	config := NewConfig()

	require.NoError(t, config.SetWithVersion("foo", "bar", 1))
	require.NoError(t, config.SetWithVersion("lu", "cy", 1))

	keys := config.List()
	require.Len(t, keys, 2)
	require.Equal(t, keys["foo"].Value, "bar")
	require.Equal(t, keys["lu"].Value, "cy")
}

func Test_Delete_WhenValueExists_DeletesValues(t *testing.T) {
	config := NewConfig()

	require.NoError(t, config.SetWithVersion("foo", "bar", 1))
	_, err := config.GetWithVersion("foo")
	require.NoError(t, err)

	err = config.Delete("foo")
	require.NoError(t, err)

	_, err = config.GetWithVersion("foo")
	require.EqualError(t, err, _errValueDoesNotExist)
}

func Test_Delete_WhenValueDoesNotExist_ReturnsError(t *testing.T) {
	config := NewConfig()

	require.EqualError(t, config.Delete("foo"), _errValueDoesNotExist)
}
