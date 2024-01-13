# navmux

## version v0.1.0

Designed to run on Raspberian and tested on a Raspberry PI 3B.

The idea to be be able to pass to and/or collate messages from different sources such as serial ports and
udp ports.  These would typically be NMEA 0183 navigation messages inlcuding AIS.

The package can also decode NMEA 0183 messages and log data. 

A yaml config file read on start up is used to define how messages are routed to different devices/ports

Typically a usb serial port adapter plugged into the Rpi would collect legacy 0183 messages and pass them to different out put serial ports, UDP clients and to a log function.

If a 2000 network bridge is used then messages on the NMEA 2000 network can be directed to NMEA 0183 devices and to OpenCPN via UDP.  Also NMEA 0183 messages can then be passed to the NMEA 2000 network.

Commercial muxs can be bought but they don't allow sharing with added code packages which is the main advantage of this approach apart from cost. 

For example:

The udp server will be useful to connect to OpenCPN to allow display of navigation data

The log package will be useful for basic logging of NMEA messages

## Example config and explantion:
```
 # A digital compass on serial port /dev/ttyUSB0
 # sends NMEA 0183 messages
 # to the queues listed in ouputs
compass:
    name: /dev/ttyUSB0
    type: serial
    baud:  4800
    outputs:
      - to_log
      - to_udp_autohelm
      - to_2000
      - to_udp_client

# A NMEA 2000 bridge sends NMEA 0183 messages to
# the queues listed in ouputs and reads messages
# from the "to_2000" queue
bridge:
    name: /dev/ttyUSB1
    type: serial
    baud: 38400
    input: to_2000
    outputs:
      - to_log
      - to_udp_client

# A AIS reciever sends NMEA 0183 messages to
# the queues listed in ouputs
ais:
    name: /dev/ttyUSB3
    type: serial
    baud: 38400
    outputs:
      - to_2000
      - to_udp_client

# The ships_log system reads the to_log queue,
# processes the NMEA 0183 messages
# and logs the data
log:
    type: ships_log
    input: to_log

# A UPD client reads the to to_udp_client queue and
# send messages to another
# RPI running openCPN plotter which listens
# on port 8011
udp_opencpn:
    type:  udp_client
    input: to_udp_client 
    server_address: 192.168.1.14:8011

# A UPD client reads the to to_udp_client queue and
# send messages to another
# RPI running openCPN plotter which listens
# on port 8011
udp_autohelm:
    type:  udp_client
    input: to_udp_autohelm
    server_address: 127.0.0.0:8006
```



