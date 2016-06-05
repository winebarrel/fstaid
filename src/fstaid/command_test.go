package fstaid

import (
	. "."
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCommandRun(t *testing.T) {
	assert := assert.New(t)

	cmd, _ := NewCommand("primary", &CommandConfig{Command: "echo", Timeout: 1})

	out := logToBuffer(func() {
		exitCode, timeout := cmd.Run()
		assert.Equal(exitCode, 0)
		assert.Equal(timeout, false)
	})

	assert.Equal("primary: stdout: \n", out)
}

func TestCommandRunWithStderr(t *testing.T) {
	assert := assert.New(t)

	cmd, _ := NewCommand("primary", &CommandConfig{Command: "bash -c 'echo 1; sleep 0.1; echo 2 >&2'", Timeout: 1})

	out := logToBuffer(func() {
		exitCode, timeout := cmd.Run()
		assert.Equal(exitCode, 0)
		assert.Equal(timeout, false)
	})

	assert.Equal("primary: stdout: 1\nprimary: stderr: 2\n", out)
}

func TestCommandRunWithArgs(t *testing.T) {
	assert := assert.New(t)

	cmd, _ := NewCommand("primary", &CommandConfig{Command: "echo", Timeout: 1})

	out := logToBuffer(func() {
		exitCode, timeout := cmd.Run("1", "2", "3")
		assert.Equal(exitCode, 0)
		assert.Equal(timeout, false)
	})

	assert.Equal("primary: stdout: 1 2 3\n", out)
}

func TestCommandRunFail(t *testing.T) {
	assert := assert.New(t)

	cmd, _ := NewCommand("primary", &CommandConfig{Command: "false", Timeout: 1})

	out := logToBuffer(func() {
		exitCode, timeout := cmd.Run()
		assert.Equal(exitCode, 1)
		assert.Equal(timeout, false)
	})

	assert.Equal("primary: failed: exitCode=1\n", out)
}

func TestCommandTimeout(t *testing.T) {
	assert := assert.New(t)

	cmd, _ := NewCommand("primary", &CommandConfig{Command: "sleep 10", Timeout: 1})

	out := logToBuffer(func() {
		exitCode, timeout := cmd.Run()
		assert.Equal(exitCode, -1)
		assert.Equal(timeout, true)
	})

	assert.Equal("primary: failed: timed out\n", out)
}
