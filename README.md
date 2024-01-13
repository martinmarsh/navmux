# navmux

Test project evaluate writing the equivalent of boat repro using go.

This was the first go version used on my boat between 2020 and 2023.  The autohelm feature was less useful than the python version it replaced
and was not used. The autohlem code was removed from this package in 2024 version 0.0.1 and a new autohelm package separate from this one was developed
to use the rudder position as well as heading.  This package now sends compass data via UDP to the new autohelm package and is now just a MUX and data logger.
By separtaing the features the system should be more robust as issues with the autohelm will not affect basic navigation data. It will also make updating the code
safer.

The concept has been improved by using yaml config which will it easier to define inputs and routing to outputs

The concept is to use a Raspberry Pi (Rpi) or similar to route NMEA 0183 messages from different inputs to various outputs.

Typically a usb serial port adapter plugged into the Rpi would collect legacy 0183 messages and pass then to different out put serial ports.

Messages can be mixed between ports and other processing packages such as a logger and udp server.
If a 2000 network bridge is used then messages on the 2000 network can be shared and devices on the 2000 network can use data from 0183 devices.
 
Commercial muxs can be bought but they don't allow sharing with added code packages which is the main advantage of this approach apart from cost. 

For example:

The udp server will be useful to connect to OpenCPN to allow display of navigation data

The log package will be useful for basic logging

Some optional packages will be available if prepeared to construct some basic hardware for the Rpi eg The compass chip reading code, tone outputs, autohelm drive 

