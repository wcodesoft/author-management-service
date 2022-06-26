package utils

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuildResponseWithError(t *testing.T) {
	expectedError := errors.New("expected error")
	expectedErrorString := expectedError.Error()
	response := BuildResponse(nil, expectedError)
	responseError := response.GetError()
	assert.Equal(t, responseError, expectedErrorString)
	assert.False(t, response.Success)
}

func TestBuildResponseCorrectValues(t *testing.T) {
	expectedMessage := "MessageTest"
	response := BuildResponse([]string{expectedMessage}, nil)
	assert.Equal(t, response.Result[0], expectedMessage)
}
