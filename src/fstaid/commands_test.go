package fstaid

import (
	. "."
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCommandResultIsSuccess(t *testing.T) {
	assert := assert.New(t)

	cr := CommandResult{ExitCode: 0, Timeout: false}
	assert.Equal(true, cr.IsSuccess())
}

func TestCommandResultIsNotSuccess(t *testing.T) {
	assert := assert.New(t)

	cr := CommandResult{ExitCode: 1, Timeout: false}
	assert.Equal(false, cr.IsSuccess())

	cr = CommandResult{ExitCode: 0, Timeout: true}
	assert.Equal(false, cr.IsSuccess())
}

func TestCheckResultSelfCheckIsSuccess(t *testing.T) {
	assert := assert.New(t)

	cr := CheckResult{}
	assert.Equal(true, cr.SelfCheckIsSuccess())

	cr = CheckResult{Self: &CommandResult{}}
	assert.Equal(true, cr.SelfCheckIsSuccess())
}

func TestCheckResultSelfCheckIsNotSuccess(t *testing.T) {
	assert := assert.New(t)

	cr := CheckResult{Self: &CommandResult{ExitCode: 1}}
	assert.Equal(false, cr.SelfCheckIsSuccess())
}

func TestCommandsCheck(t *testing.T) {
	assert := assert.New(t)

	config := &Config{
		Global:    GlobalConfig{Maxattempts: 3},
		Primary:   CommandConfig{Command: "echo -n", Timeout: 1},
		Secondary: CommandConfig{Command: "echo -n", Timeout: 1},
		Self:      CommandConfig{Command: "echo -n", Timeout: 1},
	}

	cmds, _ := NewCommands(config)
	result := cmds.Check()

	assert.Equal(true, result.Primary.IsSuccess())
	assert.Equal(false, result.Secondary.IsSuccess())
	assert.Equal(true, result.SelfCheckIsSuccess())
}

func TestCommandsCheckWithoutSelf(t *testing.T) {
	assert := assert.New(t)

	config := &Config{
		Global:    GlobalConfig{Maxattempts: 3},
		Primary:   CommandConfig{Command: "echo -n", Timeout: 1},
		Secondary: CommandConfig{Command: "echo -n", Timeout: 1},
	}

	cmds, _ := NewCommands(config)
	result := cmds.Check()

	assert.Equal(true, result.Primary.IsSuccess())
	assert.Equal(false, result.Secondary.IsSuccess())
	assert.Equal(true, result.SelfCheckIsSuccess())
}

func TestCommandsCheckWithoutSecondary(t *testing.T) {
	assert := assert.New(t)

	config := &Config{
		Global:  GlobalConfig{Maxattempts: 3},
		Primary: CommandConfig{Command: "echo -n", Timeout: 1},
	}

	cmds, _ := NewCommands(config)
	result := cmds.Check()

	assert.Equal(true, result.Primary.IsSuccess())
	assert.Equal(false, result.Secondary.IsSuccess())
	assert.Equal(true, result.SelfCheckIsSuccess())
}

func TestCommandsCheckFail(t *testing.T) {
	assert := assert.New(t)

	config := &Config{
		Global:    GlobalConfig{Maxattempts: 3},
		Primary:   CommandConfig{Command: "false", Timeout: 1},
		Secondary: CommandConfig{Command: "false", Timeout: 1},
	}

	cmds, _ := NewCommands(config)

	out := logToBuffer(func() {
		result := cmds.Check()
		assert.Equal(false, result.Primary.IsSuccess())
		assert.Equal(false, result.Secondary.IsSuccess())
	})

	assert.Equal(`primary: failed: exitCode=1
primary: failed: exitCode=1
primary: failed: exitCode=1
secondary: failed: exitCode=1
`, out)
}

func TestCommandsCheckFailWithoutSecondary(t *testing.T) {
	assert := assert.New(t)

	config := &Config{
		Global:  GlobalConfig{Maxattempts: 3},
		Primary: CommandConfig{Command: "false", Timeout: 1},
	}

	cmds, _ := NewCommands(config)

	out := logToBuffer(func() {
		result := cmds.Check()
		assert.Equal(false, result.Primary.IsSuccess())
		assert.Equal(false, result.Secondary.IsSuccess())
	})

	assert.Equal(`primary: failed: exitCode=1
primary: failed: exitCode=1
primary: failed: exitCode=1
`, out)
}

func TestCommandsSlefCheckFail(t *testing.T) {
	assert := assert.New(t)

	config := &Config{
		Global: GlobalConfig{Maxattempts: 3},
		Self:   CommandConfig{Command: "false", Timeout: 1},
	}

	cmds, _ := NewCommands(config)

	out := logToBuffer(func() {
		result := cmds.Check()
		assert.Equal(false, result.SelfCheckIsSuccess())
	})

	assert.Equal("self: failed: exitCode=1\n", out)
}
