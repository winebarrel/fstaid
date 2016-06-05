package fstaid

import (
	"fmt"
)

type Commands struct {
	Config    *Config
	Primary   *Command
	Secondary *Command
	Self      *Command
}

type CommandResult struct {
	ExitCode int
	Timeout  bool
}

func (result *CommandResult) IsSuccess() bool {
	if result.ExitCode == 0 && !result.Timeout {
		return true
	} else {
		return false
	}
}

type CheckResult struct {
	Primary   *CommandResult
	Secondary *CommandResult
	Self      *CommandResult
}

func (result *CheckResult) SelfCheckIsSuccess() bool {
	if result.Self == nil || result.Self.IsSuccess() {
		return true
	} else {
		return false
	}
}

func NewCommands(config *Config) (cmds *Commands, err error) {
	cmds = &Commands{Config: config}

	cmds.Primary, err = NewCommand("primary", &config.Primary)

	if err != nil {
		err = fmt.Errorf("primary: %s", err)
		return
	}

	if config.Secondary.Command != "" {
		cmds.Secondary, err = NewCommand("secondary", &config.Secondary)

		if err != nil {
			err = fmt.Errorf("secondary: %s", err)
			return
		}
	}

	if config.Self.Command != "" {
		cmds.Self, err = NewCommand("self", &config.Self)

		if err != nil {
			err = fmt.Errorf("self: %s", err)
			return
		}
	}

	return
}

func (cmds *Commands) Check() (result *CheckResult) {
	result = &CheckResult{
		Primary:   &CommandResult{},
		Secondary: &CommandResult{ExitCode: -1},
	}

	if cmds.Self != nil {
		result.Self = &CommandResult{}
		result.Self.ExitCode, result.Self.Timeout = cmds.Self.Run()

		if !result.Self.IsSuccess() {
			return
		}
	}

	for i := 0; i < cmds.Config.Global.Maxattempts; i++ {
		result.Primary.ExitCode, result.Primary.Timeout = cmds.Primary.Run()

		if result.Primary.IsSuccess() {
			break
		}
	}

	if result.Primary.IsSuccess() {
		return
	}

	if cmds.Secondary != nil {
		result.Secondary.ExitCode, result.Secondary.Timeout = cmds.Secondary.Run()
	}

	return
}
