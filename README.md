#BSW

A tool for monitoring local beanstalk queues

![bsw](https://cloud.githubusercontent.com/assets/5681893/13606694/e3fd8514-e544-11e5-9183-2b459f211e7d.jpg)

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
