package fstaid

import (
	. "."
	"github.com/bouk/monkey"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestServerPing(t *testing.T) {
	assert := assert.New(t)

	ginMode("release", func() {
		server := NewServer(&Config{}, nil, ioutil.Discard)

		ts := httptest.NewServer(server.Engine)
		res, _ := http.Get(ts.URL + "/ping")
		body, status := readResponse(res)

		assert.Equal(200, status)
		assert.Equal(body, `{"message":"pong"}`+"\n")
	})
}

func TestServerFail(t *testing.T) {
	assert := assert.New(t)

	checker := &Checker{}
	handleFailWithoutShutdownCalled := false
	serverShutdownCalled := false

	var guard *monkey.PatchGuard
	guard = monkey.PatchInstanceMethod(
		reflect.TypeOf(checker), "HandleFailWithoutShutdown",
		func(_ *Checker, result *CheckResult) {
			defer guard.Unpatch()
			guard.Restore()
			handleFailWithoutShutdownCalled = true

			assert.Equal(&CheckResult{
				Primary:   &CommandResult{ExitCode: 1},
				Secondary: &CommandResult{ExitCode: 1},
			}, result)

			return
		})

	monkey.Patch(ServerShutdown, func() {
		defer monkey.Unpatch(ServerShutdown)
		serverShutdownCalled = true
	})

	ginMode("release", func() {
		server := NewServer(&Config{}, checker, ioutil.Discard)
		ts := httptest.NewServer(server.Engine)
		res, _ := http.Get(ts.URL + "/fail")
		body, status := readResponse(res)

		assert.Equal(200, status)
		assert.Equal(body, `{"accepted":true}`+"\n")
	})

	assert.Equal(true, handleFailWithoutShutdownCalled)
	assert.Equal(true, serverShutdownCalled)
}

func TestServerFailWithExitCode(t *testing.T) {
	assert := assert.New(t)

	checker := &Checker{}
	handleFailWithoutShutdownCalled := false
	serverShutdownCalled := false

	var guard *monkey.PatchGuard
	guard = monkey.PatchInstanceMethod(
		reflect.TypeOf(checker), "HandleFailWithoutShutdown",
		func(_ *Checker, result *CheckResult) {
			defer guard.Unpatch()
			guard.Restore()
			handleFailWithoutShutdownCalled = true

			assert.Equal(&CheckResult{
				Primary:   &CommandResult{ExitCode: 2},
				Secondary: &CommandResult{ExitCode: 2},
			}, result)

			return
		})

	monkey.Patch(ServerShutdown, func() {
		defer monkey.Unpatch(ServerShutdown)
		serverShutdownCalled = true
	})

	ginMode("release", func() {
		server := NewServer(&Config{}, checker, ioutil.Discard)
		ts := httptest.NewServer(server.Engine)
		res, _ := http.Get(ts.URL + "/fail?exit=2")
		body, status := readResponse(res)

		assert.Equal(200, status)
		assert.Equal(body, `{"accepted":true}`+"\n")
	})

	assert.Equal(true, handleFailWithoutShutdownCalled)
	assert.Equal(true, serverShutdownCalled)
}
