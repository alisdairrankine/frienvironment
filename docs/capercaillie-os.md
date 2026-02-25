# Capercaillie OS

Capercaillie OS is a microkernel based operating system for the capercaillie compute environment.

## System Processes:

These processes run as daemons and are never shutdown

* Init - Boots system services, restarts them if the fail
* Compositor - Draws application windows to display
* Notifications - Draws system notifications to display
* Dock - Draws application launcher

## System Applications:

These programs can be opened or closed at will.

* File browser - Integrates with object storage device to show files.
* Terminal - Allows interaction with other applications via terminal device. - can spawn and interact with programs.
* Clock - A basic clock for showing the time.