# navmux

## version v0.2.0

Designed to run on Raspberian and tested on a Raspberry PI 3B.

The idea to be be able to pass to and/or collate messages from different sources such as serial ports and
udp ports. To log data collected from various sources and to monitor and select the best sources for compass and gps.
Including auto fall back on failure. It was developed into the NMEA-MUX and NMEA0183 packages which are now in the go library.
This repro. now contains the bespoke config.  relating to my personal use whereas the library packages are designed for
configuration for other uses.  However, it may be useful as a practical example used in real sea conditions.
