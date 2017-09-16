package main

import (
	"testing"
	"time"

	"github.com/henry40408/ssh-shell-resource/internal"
	"github.com/henry40408/ssh-shell-resource/pkg/mockio"
	"github.com/spacemonkeygo/errors"
	"github.com/stretchr/testify/assert"
)

func TestCheckCommandReturnDifferentResponse(t *testing.T) {
	request := CheckRequest{}

	response := CheckCommand(&request)
	assert.Equal(t, 1, len(response))

	time.Sleep(1 * time.Millisecond)

	anotherResponse := CheckCommand(&request)
	assert.Equal(t, 1, len(anotherResponse))

	responseTime := response[0].Timestamp.UnixNano()
	anotherResponseTime := anotherResponse[0].Timestamp.UnixNano()
	assert.NotEqual(t, responseTime, anotherResponseTime)
}

func TestCheckCommandReturnPreviousVersion(t *testing.T) {
	version := internal.Version{Timestamp: time.Now()}
	request := CheckRequest{
		Request: internal.Request{Version: version},
	}

	time.Sleep(1 * time.Millisecond)

	response := CheckCommand(&request)
	assert.Equal(t, 2, len(response))

	requestTime := request.Version.Timestamp.UnixNano()
	responseTime := response[0].Timestamp.UnixNano()
	assert.Equal(t, requestTime, responseTime)
}

func TestCheckCommandResponseTimeIsGreaterThanRequestTime(t *testing.T) {
	version := internal.Version{Timestamp: time.Now()}
	request := CheckRequest{
		Request: internal.Request{Version: version},
	}

	time.Sleep(1 * time.Millisecond)

	response := CheckCommand(&request)

	requestTime := request.Version.Timestamp.UnixNano()
	responseTime := response[1].Timestamp.UnixNano()
	assert.True(t, responseTime > requestTime)
}

func TestMain(t *testing.T) {
	mockio, err := mockio.NewMockIO([]byte("{}"))
	if err != nil {
		t.Error(err)
	}
	defer mockio.Cleanup()

	err = Main(mockio.In, mockio.Out)
	if err != nil {
		t.Error(err)
	}
}

func TestMainNotValidJSON(t *testing.T) {
	mockio, err := mockio.NewMockIO([]byte(""))
	if err != nil {
		t.Error(err)
	}
	defer mockio.Cleanup()

	err = Main(mockio.In, mockio.Out)
	assert.Equal(t, "InvalidJSONError: stdin is not a valid JSON", errors.GetMessage(err))
}