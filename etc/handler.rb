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
