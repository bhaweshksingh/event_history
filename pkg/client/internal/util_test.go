package internal_test

import (
	"github.com/stretchr/testify/assert"
	"event-history/pkg/client/internal"
	"testing"
)

func TestBuildURL(t *testing.T) {
	testCases := map[string]struct {
		input          func() (string, map[string]string)
		expectedResult string
		expectedError  error
	}{
		"test build url with query params": {
			input: func() (string, map[string]string) {
				return "check", map[string]string{"ko": "vo", "ke": "ve"}
			},
			expectedResult: "http://localhost/check?ke=ve&ko=vo",
		},
		"test build url with out query params": {
			input: func() (string, map[string]string) {
				return "check", map[string]string{}
			},
			expectedResult: "http://localhost/check",
		},
		"test build url with out path and params": {
			input: func() (string, map[string]string) {
				return "", map[string]string{}
			},
			expectedResult: "http://localhost",
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			path, params := testCase.input()

			res, err := internal.BuildURL("http://localhost", path, params)

			assert.Equal(t, testCase.expectedError, err)
			assert.Equal(t, testCase.expectedResult, res.String())
		})
	}
}
