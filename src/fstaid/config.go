package fstaid

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
)

type Config struct {
	Global    GlobalConfig
	Primary   CommandConfig
	Secondary CommandConfig
	Self      CommandConfig
	Handler   CommandConfig
	User      []UserConfig
}

type GlobalConfig struct {
	Port                      int
	Maxattempts               int
	AttemptInterval           int `toml:"attempt_interval"`
	Interval                  int
	Lockdir                   string
	Log                       string
	Mode                      string
	ContinueIfSelfCheckFailed bool `toml:"continue_if_self_check_failed"`
}

func (config *GlobalConfig) LockFile() string {
	return fmt.Sprintf("%s/fstaid.%d.lock", config.Lockdir, config.Port)
}

type CommandConfig struct {
	Command string
	Timeout int
}

type UserConfig struct {
	Userid   string
	Password string
}

func chkCmd(section string, cmd *CommandConfig) (err error) {
	if cmd.Command == "" {
		err = fmt.Errorf("[%s] command is required", section)
		return
	}

	if cmd.Timeout < 1 {
		err = fmt.Errorf("[%s] timeout must be '>= 1'", section)
		return
	}

	return
}

func LoadConfig(flags *Flags) (config *Config, err error) {
	config = &Config{}
	_, err = toml.DecodeFile(flags.Config, config)

	if err != nil {
		return
	}

	if config.Global.Port < 1 || config.Global.Port > 65535 {
		err = fmt.Errorf("[global] port must be '>= 1 && <= 65535'")
		return
	}

	if config.Global.Maxattempts < 1 {
		err = fmt.Errorf("[global] maxattempts must be '>= 1'")
		return
	}

	if config.Global.AttemptInterval < 1 {
		err = fmt.Errorf("[global] attempt_interval must be '>= 1'")
		return
	}

	if config.Global.Interval < 1 {
		err = fmt.Errorf("[global] interval must be '>= 1'")
		return
	}

	if config.Global.Lockdir == "" {
		config.Global.Lockdir = "/tmp"
	}

	if config.Global.Mode == "" {
		ginMode := os.Getenv("GIN_MODE")

		if ginMode != "" {
			config.Global.Mode = ginMode
		} else {
			config.Global.Mode = "debug"
		}
	}

	err = chkCmd("handler", &config.Handler)

	if err != nil {
		return
	}

	err = chkCmd("primary", &config.Primary)

	if err != nil {
		return
	}

	if config.Self.Command != "" {
		err = chkCmd("secondary", &config.Secondary)

		if err != nil {
			return
		}
	}

	if config.Self.Command != "" {
		err = chkCmd("self", &config.Self)

		if err != nil {
			return
		}
	}

	return
}
