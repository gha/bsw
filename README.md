#BSW

A tool for monitoring local beanstalk queues

## Install

    $ go get github.com/george-infinity/bsw

## Usage

    $ bsw -help
    Usage of bsw:
      -debug
            Enable debug logging
      -host string
            Beanstalk host address (default "127.0.0.1:11300")
      -refresh int
            Refresh rate of the tube list (seconds) (default 1)
