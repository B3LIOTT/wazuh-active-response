# Quarantine script

This scripts aims to isolate a potentially infected computer, and dump its memory. This script is not yet adapted for Linux.

## Network isolation
The isolation process for Windows is just the execution of these commands:
```bash
netsh advfirewall set allprofiles firewallpolicy blockinbound,blockoutbound
Get-NetAdapter | Disable-NetAdapter -Confirm:$false
```

The first one applies a full deny rule on the firewall, and the second one disables every network adapter.

To go back at the initial state:
```bash
netsh advfirewall reset
Get-NetAdapter | Enable-NetAdapter -Confirm:$false
```

(NOTE: admin privileges are needed)

## Memory dump

### Dumpit (not tested)
Dumpit is another tool for dumping memory, it is available [here](https://www.magnetforensics.com/resources/magnet-dumpit-for-windows/) for Windows and [here](https://github.com/MagnetForensics/dumpit-linux) for Linux. To execute it in quiet mode you have to use it like this:
```bash
DumpIt.exe /QUIET /OUTPUT dump.dmp
```

### Winpmem
Winpmem is another option available for Windows [here](https://github.com/Velocidex/WinPmem/releases/tag/v4.0.rc1) and for Linux [here](http://github.com/Velocidex/Linpmem). However this script is not yet designed to work on Linux.

