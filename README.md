# navmux

Test project evaluate writing the equivalent of boat repro using go.
The concept has been improved by using yaml config which will it easier to define inputs and routing to outputs

The concept is to use a Raspberry Pi (Rpi) or similar to route NMEA 0183 messages from different inputs to various outputs
Typically a usb serial port adapter plugged into the Rpi would collect legacy 0183 messages and pass then to different out put serial ports. 
Messages can be mixed between ports and other processing packages such as a logger, compass chip reader, udp server and autohelm.  If a 2000 network bridge is used then messages on the 2000 network can be shared and devices on the 2000 network can use data from 0183 devices.
 
 Commercial muxs can be bought but they don't allow sahring with added code packages which is the main advantage of this approach apart from cost. 

 For example:
 The udp server will be useful to connect to OpenCPN to allow display of navigation data
 The log package will be useful for basic logging

 Some optional packages will be available if prepeared to construct some basic hardware for the Rpi eg
 The compass chip reading code, tone outputs, autohelm drive 

 