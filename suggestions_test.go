package cli

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSuggestFlag(t *testing.T) {
	// Given
	app := buildExtendedTestCommand()

	for _, testCase := range []struct {
		provided, expected string
	}{
		{"", ""},
		{"a", "--another-flag"},
		{"hlp", "--help"},
		{"k", ""},
		{"s", "-s"},
	} {
		// When
		res := suggestFlag(app.Flags, testCase.provided, false)

		// Then
		assert.Equal(t, testCase.expected, res)
	}
}

func TestSuggestFlagHideHelp(t *testing.T) {
	// Given
	app := buildExtendedTestCommand()

	// When
	res := suggestFlag(app.Flags, "hlp", true)

	// Then
	assert.Equal(t, "--fl", res)
}

func TestSuggestFlagFromError(t *testing.T) {
	// Given
	app := buildExtendedTestCommand()

	for _, testCase := range []struct {
		command, provided, expected string
	}{
		{"", "hel", "--help"},
		{"", "soccer", "--socket"},
		{"config", "anot", "--another-flag"},
	} {
		// When
		res, _ := app.suggestFlagFromError(
			errors.New(providedButNotDefinedErrMsg+testCase.provided),
			testCase.command,
		)

		// Then
		assert.Equal(t, fmt.Sprintf(SuggestDidYouMeanTemplate+"\n\n", testCase.expected), res)
	}
}

func TestSuggestFlagFromErrorWrongError(t *testing.T) {
	// Given
	app := buildExtendedTestCommand()

	// When
	_, err := app.suggestFlagFromError(errors.New("invalid"), "")

	// Then
	assert.Error(t, err)
}

func TestSuggestFlagFromErrorWrongCommand(t *testing.T) {
	// Given
	app := buildExtendedTestCommand()

	// When
	_, err := app.suggestFlagFromError(
		errors.New(providedButNotDefinedErrMsg+"flag"),
		"invalid",
	)

	// Then
	assert.Error(t, err)
}

func TestSuggestFlagFromErrorNoSuggestion(t *testing.T) {
	// Given
	app := buildExtendedTestCommand()

	// When
	_, err := app.suggestFlagFromError(
		errors.New(providedButNotDefinedErrMsg+""),
		"",
	)

	// Then
	assert.Error(t, err)
}

func TestSuggestCommand(t *testing.T) {
	// Given
	app := buildExtendedTestCommand()

	for _, testCase := range []struct {
		provided, expected string
	}{
		{"", ""},
		{"conf", "config"},
		{"i", "i"},
		{"information", "info"},
		{"inf", "info"},
		{"con", "config"},
		{"not-existing", "info"},
	} {
		// When
		res := suggestCommand(app.Commands, testCase.provided)

		// Then
		assert.Equal(t, testCase.expected, res)
	}
}
