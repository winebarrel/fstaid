package fstaid

import (
	. "."
	"github.com/bouk/monkey"
	"github.com/stretchr/testify/assert"
	"log"
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

func TestCheckerHandleFailureWithoutShutdown(t *testing.T) {
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

			patchInstanceMethod(checker.Handler, "Run", func(guard **monkey.PatchGuard) interface{} {
				return func(_ *Command, args ...string) (int, bool) {
					defer (*guard).Unpatch()
					(*guard).Restore()

					assert.Equal(4, len(args))
					assert.Equal("1", args[0])
					assert.Equal("false", args[1])
					assert.Equal("0", args[2])
					assert.Equal("true", args[3])

					return 0, true
				}
			})

			result := &CheckResult{
				Primary:   &CommandResult{ExitCode: 1, Timeout: false},
				Secondary: &CommandResult{ExitCode: 0, Timeout: true},
			}

			checker.HandleFailureWithoutShutdown(result)

			assert.Equal(false, checker.Running)
			assert.Equal(true, fileExists(checker.Config.Global.LockFile()))
		})
	})
}

func TestCheckerHandleFailure(t *testing.T) {
	assert := assert.New(t)

	checker := &Checker{}
	HandleFailureWithoutShutdownCalled := false
	serverShutdownCalled := false

	patchInstanceMethod(checker, "HandleFailureWithoutShutdown", func(guard **monkey.PatchGuard) interface{} {
		return func(_ *Checker, result *CheckResult) {
			defer (*guard).Unpatch()
			(*guard).Restore()
			HandleFailureWithoutShutdownCalled = true
			return
		}
	})

	monkey.Patch(ServerShutdown, func() {
		defer monkey.Unpatch(ServerShutdown)
		serverShutdownCalled = true
	})

	result := &CheckResult{}
	checker.HandleFailure(result)

	assert.Equal(true, HandleFailureWithoutShutdownCalled)
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

			patchInstanceMethod(checker, "Check", func(guard **monkey.PatchGuard) interface{} {
				return func(_ *Checker) {
					defer (*guard).Unpatch()
					(*guard).Restore()
					checkCalled = true
					return
				}
			})

			checker.Run()
			time.Sleep(1 * time.Second)
			checker.Stop()

			assert.Equal(true, checkCalled)
		})
	})
}

func TestCheckerCheck(t *testing.T) {
	assert := assert.New(t)

	commands := &Commands{}
	checker := &Checker{Commands: commands}
	checkCalled := false

	patchInstanceMethod(commands, "Check", func(guard **monkey.PatchGuard) interface{} {
		return func(_ *Commands) (result *CheckResult) {
			defer (*guard).Unpatch()
			(*guard).Restore()

			checkCalled = true

			result = &CheckResult{
				Primary: &CommandResult{},
				Self:    &CommandResult{},
			}

			return
		}
	})

	checker.Check()

	assert.Equal(true, checkCalled)
}

func TestCheckerCheckFail(t *testing.T) {
	assert := assert.New(t)

	commands := &Commands{}
	checker := &Checker{Commands: commands}
	checkCalled := false
	handleFailureCalled := false

	result := &CheckResult{
		Primary:   &CommandResult{ExitCode: 1},
		Secondary: &CommandResult{ExitCode: 1},
	}

	out := logToBuffer(func() {
		patchInstanceMethod(commands, "Check", func(guard **monkey.PatchGuard) interface{} {
			return func(_ *Commands) *CheckResult {
				defer (*guard).Unpatch()
				(*guard).Restore()
				checkCalled = true
				return result
			}
		})

		patchInstanceMethod(checker, "HandleFailure", func(guard **monkey.PatchGuard) interface{} {
			return func(_ *Checker, cr *CheckResult) {
				defer (*guard).Unpatch()
				(*guard).Restore()
				assert.Equal(result, cr)
				handleFailureCalled = true
				return
			}
		})

		checker.Check()
	})

	assert.Equal(true, checkCalled)
	assert.Equal(true, handleFailureCalled)
	assert.Equal("** Health check failed **\n", out)
}

func TestCheckerPrimaryCheckFailSecondaryCheckSuccess(t *testing.T) {
	assert := assert.New(t)

	commands := &Commands{}
	checker := &Checker{Commands: commands}
	checkCalled := false
	handleFailureCalled := false

	result := &CheckResult{
		Primary:   &CommandResult{ExitCode: 1},
		Secondary: &CommandResult{},
	}

	out := logToBuffer(func() {
		patchInstanceMethod(commands, "Check", func(guard **monkey.PatchGuard) interface{} {
			return func(_ *Commands) *CheckResult {
				defer (*guard).Unpatch()
				(*guard).Restore()
				checkCalled = true
				return result
			}
		})

		patchInstanceMethod(checker, "HandleFailure", func(guard **monkey.PatchGuard) interface{} {
			return func(_ *Checker, cr *CheckResult) {
				defer (*guard).Unpatch()
				(*guard).Restore()
				assert.Equal(result, cr)
				handleFailureCalled = true
				return
			}
		})

		checker.Check()
	})

	assert.Equal(true, checkCalled)
	assert.Equal(false, handleFailureCalled)
	assert.Equal("Primary check failed, but Secondary check succeeded\n", out)
}

func TestCheckerCheckFailSecondarySelfCheckSuccess(t *testing.T) {
	assert := assert.New(t)

	commands := &Commands{}
	checker := &Checker{Commands: commands}
	checkCalled := false
	handleFailureCalled := false

	result := &CheckResult{
		Primary:       &CommandResult{ExitCode: 1},
		Secondary:     &CommandResult{ExitCode: 1},
		Self:          &CommandResult{ExitCode: 0},
		SecondarySelf: &CommandResult{ExitCode: 0},
	}

	out := logToBuffer(func() {
		patchInstanceMethod(commands, "Check", func(guard **monkey.PatchGuard) interface{} {
			return func(_ *Commands) *CheckResult {
				defer (*guard).Unpatch()
				(*guard).Restore()
				checkCalled = true
				return result
			}
		})

		patchInstanceMethod(checker, "HandleFailure", func(guard **monkey.PatchGuard) interface{} {
			return func(_ *Checker, cr *CheckResult) {
				defer (*guard).Unpatch()
				(*guard).Restore()
				assert.Equal(result, cr)
				handleFailureCalled = true
				return
			}
		})

		checker.Check()
	})

	assert.Equal(true, checkCalled)
	assert.Equal(true, handleFailureCalled)
	assert.Equal("** Health check failed **\n", out)
}

func TestCheckerCheckFailSelfCheckFail(t *testing.T) {
	assert := assert.New(t)

	commands := &Commands{}
	checker := &Checker{Commands: commands}
	checkCalled := false
	handleFailureCalled := false

	result := &CheckResult{
		Primary:       &CommandResult{ExitCode: 1},
		Secondary:     &CommandResult{ExitCode: 1},
		Self:          &CommandResult{ExitCode: 0},
		SecondarySelf: &CommandResult{ExitCode: 1},
	}

	out := logToBuffer(func() {
		patchInstanceMethod(commands, "Check", func(guard **monkey.PatchGuard) interface{} {
			return func(_ *Commands) *CheckResult {
				defer (*guard).Unpatch()
				(*guard).Restore()
				checkCalled = true
				return result
			}
		})

		patchInstanceMethod(checker, "HandleFailure", func(guard **monkey.PatchGuard) interface{} {
			return func(_ *Checker, cr *CheckResult) {
				defer (*guard).Unpatch()
				(*guard).Restore()
				assert.Equal(result, cr)
				handleFailureCalled = true
				return
			}
		})

		checker.Check()
	})

	assert.Equal(true, checkCalled)
	assert.Equal(false, handleFailureCalled)
	assert.Equal("Primary/Secondary check failed, but Self check failed\n", out)
}

func TestCheckerSelfCheckFail(t *testing.T) {
	assert := assert.New(t)

	commands := &Commands{}

	checker := &Checker{
		Config:   &Config{Global: GlobalConfig{ContinueIfSelfCheckFailed: false}},
		Commands: commands,
	}

	checkCalled := false
	fatalfCalled := false

	result := &CheckResult{
		Primary: &CommandResult{},
		Self:    &CommandResult{ExitCode: 1},
	}

	patchInstanceMethod(commands, "Check", func(guard **monkey.PatchGuard) interface{} {
		return func(_ *Commands) *CheckResult {
			defer (*guard).Unpatch()
			(*guard).Restore()
			checkCalled = true
			return result
		}
	})

	monkey.Patch(log.Fatalf, func(format string, v ...interface{}) {
		defer monkey.Unpatch(log.Fatalf)
		fatalfCalled = true
		assert.Equal("** Self check failed **", format)
	})

	checker.Check()

	assert.Equal(true, checkCalled)
	assert.Equal(true, fatalfCalled)
}

func TestCheckerSelfCheckFailAndContinue(t *testing.T) {
	assert := assert.New(t)

	commands := &Commands{}

	checker := &Checker{
		Config:   &Config{Global: GlobalConfig{ContinueIfSelfCheckFailed: true}},
		Commands: commands,
	}

	checkCalled := false

	result := &CheckResult{
		Self: &CommandResult{ExitCode: 1},
	}

	out := logToBuffer(func() {
		patchInstanceMethod(commands, "Check", func(guard **monkey.PatchGuard) interface{} {
			return func(_ *Commands) *CheckResult {
				defer (*guard).Unpatch()
				(*guard).Restore()
				checkCalled = true
				return result
			}
		})

		checker.Check()
	})

	assert.Equal(true, checkCalled)
	assert.Equal("** Self check failed ** (Health check will continue)\n", out)
}
