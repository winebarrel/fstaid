# fstaid

fstaid is a daemon that monitors the health condition of the server and executes the script if there is any problem.

[![Build Status](https://travis-ci.org/winebarrel/fstaid.svg?branch=master)](https://travis-ci.org/winebarrel/fstaid)

![](https://cdn.pbrd.co/images/fjskjOr4.png)

## Flowchart

![](https://cdn.pbrd.co/images/GHdiTqq.png)

## Usage

```
Usage of fstaid:
  -config string
    	config file path (default "fstaid.toml")
  -version
    	show version
```

## Configuration

```toml
[global]
port = 8080
interval = 1
maxattempts = 3
attempt_interval = 1
#lockdir = "/tmp"
#log = "/var/log/fstaid.log"
#mode = "debug"(default) / "release"
#continue_if_self_check_failed = false

[handler]
command = "/usr/libexec/fstaid/handler.rb"
timeout = 300

[primary]
command = "curl -s -f server-01"
timeout = 3

[secondary]
command = "curl -s -f -x server-02:8080 server-01"
timeout = 3
# secondary check is not required

[self]
command = "curl -s -f 169.254.169.254/latest/meta-data/instance-id"
timeout = 3
# self check is not required

[[user]]
userid = "foo"
password = "bar"
```

## Handler Example

```ruby
#!/usr/bin/env ruby
# The handler is called when both the primary check and the secondary check fail

primary_exit_code   =  ARGV[0].to_i
primary_timeout     = (ARGV[1] == 'true')
secondary_exit_code =  ARGV[2].to_i
secondary_timeout   = (ARGV[3] == 'true')

def failover
  # to fail over
end

if primary_exit_code != 0
  if secondary_timeout
    # Nothing to do if the secondary-check times out
    exit 1
  end

  failover
elsif primary_timeout
  # Nothing to do if the primary-check times out
  exit 2
end
```

## Web Interface

* `GET /ping`: ping
* `GET /fail`: force failover
