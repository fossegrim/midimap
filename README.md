# midimap
## Setup
### Install portmidi
```sh
# Debian derivatives
sudo apt install libportmidi-dev
# MacOS
brew install portmidi
```
### Verify that portmidi is working
```sh
# The following command outputs the midi events of the default midi device to stdout.
# Try sending some midi events and see if they are outputted to test that portmidi is working correctly
go run midilog.go
```
## Alternatives
| Alternative                                                                         | License                              | Platform     | Comment          |
|-------------------------------------------------------------------------------------|--------------------------------------|--------------|------------------|
| [Midikey2Key](https://midikey2key.de)                                               | Propietary freeware                  | Windows only |                  |
| [xobs/midi-to-keypress](https://github.com/xobs/midi-to-keypress)                   | Source available propietary freeware | Windows only |                  |
| [davidlukerice/midi-to-keypress](https://github.com/davidlukerice/midi-to-keypress) | MIT                                  | Windows only | Uses Electron ;( |
| [mwicat/midimap](https://github.com/mwicat/midimap)                                 | Source available propietary freeware | Windows only | Uses Python ;(   |