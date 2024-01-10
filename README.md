# navmux

Test project evaluate writing the equivalent of the boat repro using go.

The concept has been improved by using yaml config which will it easier to define inputs and routing to outputs

A Raspberry Pi 3 (Rpi) or routes NMEA 0183 messages from different inputs to various outputs.

Typically a usb serial port adapter plugged into the Rpi collects legacy 0183 messages and pass then to different out put serial and udp ports.

Messages can be mixed between ports and other processing packages such as a logger and udp server.
If a 2000 network bridge is used then messages on the 2000 network can be shared and devices on the 2000 network can use data from 0183 devices.

Commercial muxs can be bought but they don't allow sharing with added code packages which is the main advantage of this approach apart from the cost saving. 

For example:

The udp server is used to connect to OpenCPN to allow display of combined AIS and navigation data 

The log package will be useful for basic back up logging

The 2023 update is to remove the auto-helm feature which is moved to another repro.


To run automatically add to .bashrc eg sudo nano /home/pi/.bashrc:

cd /home/pi/run_time
tmux new -s navmux "./navmux run"
tmux new -s console
tmux attach -t navmux

