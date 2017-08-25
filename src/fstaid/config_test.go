package fstaid

import (
	. "."
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	assert := assert.New(t)

	tml := `
[global]
port = 8080
interval = 1
maxattempts = 2
attempt_interval = 3.0
lockdir = "/var/tmp"
log = "/var/log/fstaid.log"
mode = "release"
continue_if_self_check_failed = true

[handler]
command = "handler.rb"
timeout = 300

[primary]
command = "echo 1"
timeout = 3

[secondary]
command = "echo 2"
timeout = 4

[self]
command = "echo 3"
timeout = 5

[[user]]
userid = "foo"
password = "bar"
  `

	tempFile(tml, func(f *os.File) {
		flag := &Flags{Config: f.Name()}
		config, _ := LoadConfig(flag)

		assert.Equal(Config{
			Global: GlobalConfig{
				Port:            8080,
				Maxattempts:     2,
				AttemptInterval: 3.0,
				Interval:        1,
				Lockdir:         "/var/tmp",
				Log:             "/var/log/fstaid.log",
				Mode:            "release",
				ContinueIfSelfCheckFailed: true,
			},
			Primary: CommandConfig{
				Command: "echo 1",
				Timeout: 3,
			},
			Secondary: CommandConfig{
				Command: "echo 2",
				Timeout: 4,
			},
			Self: CommandConfig{
				Command: "echo 3",
				Timeout: 5,
			},
			Handler: CommandConfig{
				Command: "handler.rb",
				Timeout: 300,
			},
			User: []UserConfig{
				UserConfig{
					Userid:   "foo",
					Password: "bar",
				},
			},
		}, *config)

		assert.Equal("/var/tmp/fstaid.8080.lock", config.Global.LockFile())
	})
}

func TestLoadConfigWithoutAny(t *testing.T) {
	assert := assert.New(t)

	tml := `
[global]
port = 8080
interval = 1
maxattempts = 2
attempt_interval = 3.0

[handler]
command = "handler.rb"
timeout = 300

[primary]
command = "echo 1"
timeout = 3
  `

	tempFile(tml, func(f *os.File) {
		flag := &Flags{Config: f.Name()}
		config, _ := LoadConfig(flag)

		assert.Equal(Config{
			Global: GlobalConfig{
				Port:            8080,
				Maxattempts:     2,
				Interval:        1,
				AttemptInterval: 3.0,
				Lockdir:         "/tmp",
				Log:             "",
				Mode:            "debug",
			},
			Primary: CommandConfig{
				Command: "echo 1",
				Timeout: 3,
			},
			Secondary: CommandConfig{
				Command: "",
				Timeout: 0,
			},
			Self: CommandConfig{
				Command: "",
				Timeout: 0,
			},
			Handler: CommandConfig{
				Command: "handler.rb",
				Timeout: 300,
			},
		}, *config)

		assert.Equal("/tmp/fstaid.8080.lock", config.Global.LockFile())
	})
}
