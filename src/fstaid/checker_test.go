package fstaid

import (
	. "."
	"github.com/bouk/monkey"
	"github.com/stretchr/testify/assert"
	"log"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestCheckerLockFile(t *testing.T) {
	assert := assert.New(t)

	tempDir(func(path string) {
		checker := Checker{
			Config: &Config{
				Global: GlobalConfig{
					Port:    8080,
					Lockdir: path,
				},
			},
		}

		checker.TouchLockFile()
		assert.Equal(true, fileExists(checker.Config.Global.LockFile()))

		monkey.Patch(log.Fatalf, func(format string, v ...interface{}) {
			defer monkey.Unpatch(log.Fatalf)
			assert.Equal("fstaid is is already locked: %s", format)
			assert.Equal(1, len(v))
			assert.Equal(checker.Config.Global.LockFile(), v[0])
		})

		checker.CheckLockFile()
	})
}

func TestCheckerHandleFailWithoutShutdown(t *testing.T) {
	assert := assert.New(t)

	logToBuffer(func() {
		tempDir(func(path string) {
			checker := Checker{
				Config: &Config{
					Global: GlobalConfig{
						Port:    8080,
						Lockdir: path,
					},
				},
				Handler: &Command{CmdArgs: []string{"handler.rb"}},
				Running: true,
			}

			var guard *monkey.PatchGuard
			guard = monkey.PatchInstanceMethod(
				reflect.TypeOf(checker.Handler), "Run",
				func(_ *Command, args ...string) (int, bool) {
					defer guard.Unpatch()
					guard.Restore()

					assert.Equal(4, len(args))
					assert.Equal("1", args[0])
					assert.Equal("false", args[1])
					assert.Equal("0", args[2])
					assert.Equal("true", args[3])

					return 0, true
				})

			result := &CheckResult{
				Primary:   &CommandResult{ExitCode: 1, Timeout: false},
				Secondary: &CommandResult{ExitCode: 0, Timeout: true},
			}

			checker.HandleFailWithoutShutdown(result)

			assert.Equal(false, checker.Running)
			assert.Equal(true, fileExists(checker.Config.Global.LockFile()))
		})
	})
}

func TestCheckerHandleFail(t *testing.T) {
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
			return
		})

	monkey.Patch(ServerShutdown, func() {
		defer monkey.Unpatch(ServerShutdown)
		serverShutdownCalled = true
	})

	result := &CheckResult{}
	checker.HandleFail(result)

	assert.Equal(true, handleFailWithoutShutdownCalled)
	assert.Equal(true, serverShutdownCalled)
}

func TestCheckerRunStop(t *testing.T) {
	assert := assert.New(t)

	logToBuffer(func() {
		tempDir(func(path string) {
			checker := &Checker{
				Config: &Config{Global: GlobalConfig{
					Port:     8080,
					Lockdir:  path,
					Interval: 2,
				}},
				Running:   true,
				WaitGroup: &sync.WaitGroup{},
			}

			checkCalled := false

			var guard *monkey.PatchGuard
			guard = monkey.PatchInstanceMethod(
				reflect.TypeOf(checker), "Check",
				func(_ *Checker) {
					defer guard.Unpatch()
					guard.Restore()
					checkCalled = true
					return
				})

			checker.Run()
			time.Sleep(1 * time.Second)
			checker.Stop()

			assert.Equal(true, checkCalled)
		})
	})
}
